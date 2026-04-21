package download

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	bufSize  = 32 * 1024 // 32KB буфер для потоковой записи
	retryMax = 3         // количество попыток при ошибке сети
)

type progressFn func(downloaded, total int64)

// DownloadModel скачивает GGUF файл с Hugging Face Hub в кэш-директорию.
// repoID: "TheBloke/CodeLlama-7B-GGUF"
// filename: "codellama-7b.Q4_K_M.gguf"
// cacheDir: если пусто, используется $HOME/.cache/iollama/models
// Возвращает путь к скачанному файлу.
func DownloadModel(repoID, filename, cacheDir string) (string, error) {
	return DownloadModelWithProgress(repoID, filename, cacheDir, nil)
}

// DownloadModelWithProgress скачивает GGUF файл с поддержкой коллбэка прогресса.
// onProgress: вызывается после каждого чтения буфера с текущим прогрессом.
// Параметры onProgress: (скачано байт, всего байт, -1 если неизвестно).
func DownloadModelWithProgress(repoID, filename, cacheDir string, onProgress progressFn) (string, error) {
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

	stat, err := os.Stat(localPath)
	if err == nil && stat.Size() > 0 {
		fmt.Fprintf(os.Stderr, "Model already exists at %s, skipping download.\n", localPath)
		return localPath, nil
	}

	url := fmt.Sprintf("https://huggingface.co/%s/resolve/main/%s", repoID, filename)
	fmt.Fprintf(os.Stderr, "Downloading %s ...\n", url)

	return downloadWithRetry(url, localPath, onProgress, retryMax)
}

// downloadWithRetry выполняет загрузку с повторными попытками при ошибках сети.
func downloadWithRetry(url, localPath string, onProgress progressFn, retries int) (string, error) {
	var lastErr error
	for i := 0; i < retries; i++ {
		lastErr = downloadFile(url, localPath, onProgress)
		if lastErr == nil {
			return localPath, nil
		}
		fmt.Fprintf(os.Stderr, "Download failed (attempt %d/%d): %v\n", i+1, retries, lastErr)
	}
	return "", fmt.Errorf("download failed after %d attempts: %w", retries, lastErr)
}

// downloadFile выполняет потоковую загрузку файла с буферизацией по частям.
func downloadFile(url, localPath string, onProgress progressFn) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download: %s", resp.Status)
	}

	total := resp.ContentLength
	downloaded := int64(0)

	out, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer out.Close()

	buf := make([]byte, bufSize)
	for {
		n, rerr := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := out.Write(buf[:n]); werr != nil {
				return werr
			}
			downloaded += int64(n)
			if onProgress != nil {
				onProgress(downloaded, total)
			}
		}
		if rerr != nil {
			if rerr.Error() == "EOF" {
				break
			}
			return rerr
		}
	}

	fmt.Fprintf(os.Stderr, "Saved to %s\n", localPath)
	return nil
}

// ParseHFURL парсит URL Hugging Face и извлекает repoID и имя файла.
// Пример: "https://huggingface.co/TheBloke/CodeLlama-7B-GGUF/resolve/main/codellama-7b.Q4_K_M.gguf"
// Возвращает (repoID, filename, ok).
func ParseHFURL(url string) (repoID, filename string, isOk bool) {
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "", "", false
	}
	repoID = strings.Join(parts[:len(parts)-2], "/")
	filename = parts[len(parts)-1]
	return repoID, filename, true
}
