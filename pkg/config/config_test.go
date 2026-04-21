package config

import (
	"os"
	"testing"
)

func TestLoad_MissingModelPath(t *testing.T) {
	os.Clearenv()
	_, err := Load()
	if err == nil {
		t.Error("expected error when MODEL_PATH missing")
	}
}

func TestLoad_WithValidPath(t *testing.T) {
	os.Clearenv()
	tmpFile, err := os.CreateTemp("", "model.gguf")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	os.Setenv("IOLLAMA_MODEL_PATH", tmpFile.Name())
	os.Setenv("IOLLAMA_PORT", "9090")
	os.Setenv("IOLLAMA_TEMPERATURE", "0.9")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != 9090 {
		t.Errorf("port = %d, want 9090", cfg.Port)
	}
	if cfg.Temperature != 0.9 {
		t.Errorf("temperature = %f, want 0.9", cfg.Temperature)
	}
}

