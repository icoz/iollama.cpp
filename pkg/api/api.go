package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// LLMEngine интерфейс для генерации текста (реализуется пакетом llm)
type LLMEngine interface {
	Generate(ctx context.Context, prompt string, maxTokens int, temperature float64) (string, error)
	GenerateStream(ctx context.Context, prompt string, maxTokens int, temperature float64, callback func(token string) error) error
}

type Server struct {
	engine LLMEngine
	model  string // имя модели (путь или ID)
}

func NewServer(engine LLMEngine, modelName string) *Server {
	return &Server{engine: engine, model: modelName}
}

// Handler возвращает http.Handler с маршрутами
func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	r.Get("/v1/models", s.listModels)
	r.Post("/v1/chat/completions", s.chatCompletions)

	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type ModelResponse struct {
	Object string `json:"object"`
	Data   []struct {
		ID     string `json:"id"`
		Object string `json:"object"`
	} `json:"data"`
}

func (s *Server) listModels(w http.ResponseWriter, r *http.Request) {
	resp := ModelResponse{
		Object: "list",
		Data: []struct {
			ID     string `json:"id"`
			Object string `json:"object"`
		}{{ID: s.model, Object: "model"}},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

func (s *Server) chatCompletions(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Извлекаем последнее сообщение пользователя как prompt
	var prompt string
	for i := len(req.Messages) - 1; i >= 0; i-- {
		if req.Messages[i].Role == "user" {
			prompt = req.Messages[i].Content
			break
		}
	}
	if prompt == "" {
		http.Error(w, "no user message found", http.StatusBadRequest)
		return
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 2048 // default
	}
	temp := req.Temperature
	if temp == 0 {
		temp = 0.7
	}

	if req.Stream {
		s.handleStream(w, r, prompt, maxTokens, temp)
		return
	}

	// Синхронный режим
	text, err := s.engine.Generate(r.Context(), prompt, maxTokens, temp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := ChatResponse{
		ID:      "chatcmpl-123",
		Object:  "chat.completion",
		Created: 1234567890,
		Model:   req.Model,
		Choices: []Choice{{
			Index:        0,
			Message:      Message{Role: "assistant", Content: text},
			FinishReason: "stop",
		}},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleStream(w http.ResponseWriter, r *http.Request, prompt string, maxTokens int, temp float64) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ctx := r.Context()
	err := s.engine.GenerateStream(ctx, prompt, maxTokens, temp, func(token string) error {
		// Формируем SSE сообщение в формате OpenAI
		data := map[string]interface{}{
			"id":      "chatcmpl-123",
			"object":  "chat.completion.chunk",
			"created": 1234567890,
			"model":   s.model,
			"choices": []map[string]interface{}{{
				"index": 0,
				"delta": map[string]string{"content": token},
			}},
		}
		jsonData, _ := json.Marshal(data)
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()
		return nil
	})
	if err != nil && err != ctx.Err() {
		log.Printf("stream error: %v", err)
	}
	// Отправляем сигнал завершения
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}
