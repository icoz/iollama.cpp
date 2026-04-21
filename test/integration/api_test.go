//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestRealModelIntegration(t *testing.T) {
	modelPath := os.Getenv("IOLLAMA_TEST_MODEL")
	if modelPath == "" {
		t.Skip("IOLLAMA_TEST_MODEL not set, skipping integration test")
	}
	// Запускаем сервер в отдельной горутине (упрощённо, в реальности лучше запустить процесс)
	// Здесь мы предполагаем, что сервер уже запущен на :8080
	// Для полноценного теста можно использовать testcontainers, но это выходит за рамки.

	client := &http.Client{Timeout: 30 * time.Second}
	reqBody := `{"model":"test","messages":[{"role":"user","content":"Say hello"}],"stream":false}`
	resp, err := client.Post("http://localhost:8080/v1/chat/completions", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Skipf("server not running: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		t.Error("no choices in response")
	}
}

