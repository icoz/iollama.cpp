package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadModel скачивает GGUF файл с Hugging Face Hub в кэш-директорию.
// repoID: "TheBloke/CodeLlama-7B-GGUF"
// filename: "codellama-7b.Q4_K_M.gguf"
// cacheDir: если пусто, используется $HOME/.cache/iollama/models
func DownloadModel(repoID, filename, cacheDir string) (string, error) {
	if cacheDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot get home dir: %w", err)
		}
		cacheDir = filepath.Join(home, ".cache", "iollama", "models")
	}
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}
	localPath := filepath.Join(cacheDir, filename)

	// Проверяем, есть ли уже файл
	if _, err := os.Stat(localPath); err == nil {
		fmt.Fprintf(os.Stderr, "Model already exists at %s, skipping download.\n", localPath)
		return localPath, nil
	}

	url := fmt.Sprintf("https://huggingface.co/%s/resolve/main/%s", repoID, filename)
	fmt.Fprintf(os.Stderr, "Downloading %s ...\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download: %s", resp.Status)
	}

	out, err := os.Create(localPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Прогресс в stderr (можно улучшить)
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Fprintf(os.Stderr, "Saved to %s\n", localPath)
	return localPath, nil
}

