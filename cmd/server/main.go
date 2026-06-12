package main

import (
	"log"
	"net/http"
	"os"

	"rag-demo/cmd/server/routes"
	"rag-demo/internal/rag"
)

func main() {
	embedModel := envOrDefault("EMBED_MODEL", "gemini-embedding-2")
	llmModel := envOrDefault("LLM_MODEL", "gemini-2.5-flash")
	apiKey := os.Getenv("GOOGLE_API_KEY")
	indexPath := envOrDefault("INDEX_PATH", "index/index.json")
	docsDir := envOrDefault("DOCS_DIR", "data/docs")
	addr := envOrDefault("ADDR", ":8080")

	if apiKey == "" {
		log.Fatal("GOOGLE_API_KEY environment variable is required")
	}

	r := rag.New(embedModel, llmModel, apiKey, indexPath, docsDir)

	if err := r.LoadIndex(); err != nil {
		log.Printf("no existing index found, building from %s ...", docsDir)
		if err := r.BuildIndex(); err != nil {
			log.Fatalf("build index: %v", err)
		}
	}

	mux := http.NewServeMux()
	routes.Register(mux, r)

	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
