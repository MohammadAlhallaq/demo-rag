package rag

import (
	"fmt"
	"log"
	"strings"

	"rag-demo/internal/chunk"
	"rag-demo/internal/embed"
	"rag-demo/internal/ingest"
	"rag-demo/internal/llm"
	"rag-demo/internal/retrieve"
	"rag-demo/internal/storage"
	"rag-demo/internal/types"
)

const (
	defaultTokenSize = 250
	defaultOverlap   = 5
	defaultTopK      = 3
)

type RAG struct {
	EmbedClient *embed.Client
	LLMClient   *llm.Client
	IndexPath   string
	DocsDir     string
	TokenSize   int
	Overlap     int
	TopK        int
	index       types.Index
	chunked     map[string][]string
}

func New(embedModel, llmModel, apiKey, indexPath, docsDir string) *RAG {
	return &RAG{
		EmbedClient: embed.New(embedModel, apiKey),
		LLMClient:   llm.New(llmModel, apiKey),
		IndexPath:   indexPath,
		DocsDir:     docsDir,
		TokenSize:   defaultTokenSize,
		Overlap:     defaultOverlap,
		TopK:        defaultTopK,
	}
}

func (r *RAG) BuildIndex() error {
	docs, err := ingest.LoadDocs(r.DocsDir)
	if err != nil {
		return fmt.Errorf("load docs: %w", err)
	}

	r.chunked = chunk.Chunks(docs, r.TokenSize, r.Overlap)

	var all []types.Doc
	for source, chunks := range r.chunked {
		for i, text := range chunks {
			vec, err := r.EmbedClient.Embed(text)
			if err != nil {
				log.Printf("embed chunk %s[%d]: %v", source, i, err)
				continue
			}
			all = append(all, types.Doc{
				ID:     fmt.Sprintf("%s-%d", source, i),
				Source: source,
				Text:   text,
				Vec:    vec,
			})
		}
	}

	r.index = types.Index{Docs: all}
	return storage.SaveIndex(r.IndexPath, r.index)
}

func (r *RAG) LoadIndex() error {
	idx, err := storage.LoadIndex(r.IndexPath)
	if err != nil {
		return err
	}
	r.index = idx
	return nil
}

func (r *RAG) Query(query string) (string, []string, error) {

	queryVec, err := r.EmbedClient.Embed(query)
	var sources []string
	if err != nil {
		return "", sources, fmt.Errorf("embed query: %w", err)
	}

	docs := retrieve.TopK(r.index, queryVec, r.TopK)

	var ctx strings.Builder
	for _, d := range docs {
		sources = append(sources, d.Source)
		ctx.WriteString(d.Text)
		ctx.WriteString("\n---\n")
	}

	answer, err := r.LLMClient.Ask(ctx.String(), query)
	if err != nil {
		return "", sources, fmt.Errorf("ask llm: %w", err)
	}

	return strings.TrimSpace(answer), sources, nil
}
