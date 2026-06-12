package routes

import (
	"encoding/json"
	"net/http"

	"rag-demo/internal/rag"
)

type queryRequest struct {
	Query string `json:"query"`
}

type queryResponse struct {
	Answer  string   `json:"answer"`
	Sources []string `json:"sources"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func Register(mux *http.ServeMux, r *rag.RAG) {
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("POST /api/query", handleQuery(r))
	mux.HandleFunc("POST /api/ingest", handleIngest(r))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleQuery(r *rag.RAG) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var body queryRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if body.Query == "" {
			writeError(w, http.StatusBadRequest, "query is required")
			return
		}

		answer, sources, err := r.Query(body.Query)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(queryResponse{Answer: answer, Sources: sources})
	}
}

func handleIngest(r *rag.RAG) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := r.BuildIndex(); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "index rebuilt"})
	}
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorResponse{Error: msg})
}
