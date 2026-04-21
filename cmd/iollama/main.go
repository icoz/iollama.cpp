package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/icoz/iollama.cpp/pkg/api"
	"github.com/icoz/iollama.cpp/pkg/config"
	"github.com/icoz/iollama.cpp/pkg/llm"
	"github.com/icoz/iollama.cpp/pkg/download"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "iollama",
		Short: "Lightweight OpenAI-compatible server for GGUF models",
	}

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the API server",
		RunE:  runServe,
	}
	serveCmd.Flags().String("model", "", "Path to GGUF model file (overrides IOLLAMA_MODEL_PATH)")
	serveCmd.Flags().String("port", "", "Port to listen on (overrides IOLLAMA_PORT)")
	serveCmd.Flags().Int("max-tokens", 0, "Max tokens for generation")
	serveCmd.Flags().Float64("temperature", 0, "Temperature for sampling")
	rootCmd.AddCommand(serveCmd)

	// Команда download (будет реализована в промпте 7)
	downloadCmd := &cobra.Command{
		Use:   "download [repo_id] [filename]",
		Short: "Download a GGUF model from Hugging Face Hub",
		Args:  cobra.ExactArgs(2),
		RunE:  runDownload,
	}
	downloadCmd.Flags().String("cache-dir", "", "Cache directory for models")
	rootCmd.AddCommand(downloadCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runServe(cmd *cobra.Command, args []string) error {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}
	// Переопределение флагами
	if modelFlag, _ := cmd.Flags().GetString("model"); modelFlag != "" {
		cfg.ModelPath = modelFlag
	}
	if portFlag, _ := cmd.Flags().GetString("port"); portFlag != "" {
		var p int
		fmt.Sscanf(portFlag, "%d", &p)
		cfg.Port = p
	}
	if mt, _ := cmd.Flags().GetInt("max-tokens"); mt > 0 {
		cfg.MaxTokens = mt
	}
	if temp, _ := cmd.Flags().GetFloat64("temperature"); temp > 0 {
		cfg.Temperature = temp
	}

	// Инициализируем LLM
	log.Printf("Loading model from %s", cfg.ModelPath)
	model, err := llm.LoadModel(cfg.ModelPath)
	if err != nil {
		return fmt.Errorf("failed to load model: %w", err)
	}
	defer model.Close()

	// Создаём API сервер
	apiServer := api.NewServer(model, cfg.ModelPath)
	handler := apiServer.Handler()

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}

func runDownload(cmd *cobra.Command, args []string) error {
	repoID := args[0]
	filename := args[1]
	cacheDir, _ := cmd.Flags().GetString("cache-dir")
	path, err := download.DownloadModel(repoID, filename, cacheDir)
	if err != nil {
		return err
	}
	fmt.Printf("Model downloaded to: %s\n", path)
	return nil
}

