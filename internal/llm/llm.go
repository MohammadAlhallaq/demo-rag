package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type part struct {
	Text string `json:"text"`
}

type content struct {
	Parts []part `json:"parts"`
}

type geminiRequest struct {
	Contents          []content `json:"contents"`
	SystemInstruction *content  `json:"systemInstruction,omitempty"`
}

type geminiResponse struct {
	Candidates []struct {
		Content content `json:"content"`
	} `json:"candidates"`
}

type Client struct {
	APIKey string
	Model  string
}

func New(model, apiKey string) *Client {
	return &Client{Model: model, APIKey: apiKey}
}

func (c *Client) Ask(context, query string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", c.Model, c.APIKey)

	systemPrompt := "You are a concise assistant. Answer the question directly using ONLY the provided context. Do not repeat or echo back the context formatting, markdown headings, or separators. If the context does not contain relevant information, say \"I don't know\"."

	body := geminiRequest{
		Contents: []content{
			{
				Parts: []part{
					{Text: "Context:\n" + context + "\n\nQuestion:\n" + query},
				},
			},
		},
		SystemInstruction: &content{
			Parts: []part{
				{Text: systemPrompt},
			},
		},
	}

	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Gemini API: %s", resp.Status)
	}

	var result geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Candidates) == 0 {
		return "", fmt.Errorf("no candidates returned")
	}

	if len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content in candidate")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}
