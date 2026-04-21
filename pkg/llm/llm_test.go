package llm

import (
	"context"
	"os"
	"testing"
)

func TestLoadModel_FileNotFound(t *testing.T) {
	_, err := LoadModel("/nonexistent/model.gguf")
	if err == nil {
		t.Fatal("expected error for missing model")
	}
}

func TestGenerate(t *testing.T) {
	// Создаём временный файл как заглушку модели
	tmpFile, err := os.CreateTemp("", "model-*.gguf")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	model, err := LoadModel(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer model.Close()

	ctx := context.Background()
	resp, err := model.Generate(ctx, "test prompt", 100, 0.7)
	if err != nil {
		t.Fatal(err)
	}
	if resp == "" {
		t.Error("empty response")
	}
}

func TestGenerateStream(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "model-*.gguf")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	model, err := LoadModel(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer model.Close()

	var tokens []string
	err = model.GenerateStream(context.Background(), "hello world", 10, 0.7, func(token string) error {
		tokens = append(tokens, token)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) == 0 {
		t.Error("no tokens received")
	}
}

