package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockEngine struct{}

func (m *mockEngine) Generate(ctx context.Context, prompt string, maxTokens int, temperature float64) (string, error) {
	return "mocked response", nil
}
func (m *mockEngine) GenerateStream(ctx context.Context, prompt string, maxTokens int, temperature float64, callback func(token string) error) error {
	return callback("mocked ")
}

func TestChatCompletionsSync(t *testing.T) {
	engine := &mockEngine{}
	s := NewServer(engine, "test-model")
	handler := s.Handler()

	body := `{"model":"test","messages":[{"role":"user","content":"hello"}],"stream":false}`
	req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d, want 200", rec.Code)
	}
	var resp ChatResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content != "mocked response" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

