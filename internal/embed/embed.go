package embed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strings"
	"unicode"
)

var htmlTagRx = regexp.MustCompile(`<[^>]*>`)

type embedContentRequest struct {
	Content content `json:"content"`
}

type content struct {
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text"`
}

type embedContentResponse struct {
	Embedding embedding `json:"embedding"`
}

type embedding struct {
	Values []float64 `json:"values"`
}

type Client struct {
	APIKey string
	Model  string
}

func New(model, apiKey string) *Client {
	return &Client{Model: model, APIKey: apiKey}
}

func sanitize(text string) string {
	text = html.UnescapeString(text)

	text = htmlTagRx.ReplaceAllString(text, "")

	var b strings.Builder
	for _, r := range text {
		if r == '\n' || r == '\t' || (r >= 32 && !unicode.IsControl(r)) {
			b.WriteRune(r)
		}
	}
	text = b.String()

	text = strings.ToLower(text)

	text = strings.Join(strings.Fields(text), " ")

	return strings.TrimSpace(text)
}

func (c *Client) Embed(text string) ([]float64, error) {
	text = sanitize(text)
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:embedContent?key=%s", c.Model, c.APIKey)

	body := embedContentRequest{
		Content: content{
			Parts: []part{
				{Text: text},
			},
		},
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embedding API: %s", resp.Status)
	}

	var result embedContentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Embedding.Values) == 0 {
		return nil, fmt.Errorf("no embedding values returned")
	}

	return result.Embedding.Values, nil
}
