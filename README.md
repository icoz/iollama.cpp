# iollama.cpp

[![Go Reference](https://img.shields.io/badge/go-1.24-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

Лёгкий HTTP-сервер на Go, предоставляющий OpenAI-совместимый API для LLM в формате GGUF через llama.cpp.

## Особенности

- **Zero CGO** – использует purego для вызова llama.cpp
- **Единый бинарник** – всё необходимое в одном файле
- **Автоматическая загрузка моделей** с Hugging Face Hub
- **OpenAI-совместимый API** с поддержкой потоковой передачи (SSE)
- **Гибкая конфигурация** через переменные окружения или флаги CLI
- **Docker-образ** для быстрого развёртывания

## Быстрый старт

```bash
# Скачайте модель
go run ./cmd/iollama download TheBloke/TinyLlama-1.1B-GGUF tinyllama-1.1b.Q4_K_M.gguf

# Запустите сервер
export IOLLAMA_MODEL_PATH=$HOME/.cache/iollama/models/tinyllama-1.1b.Q4_K_M.gguf
go run ./cmd/iollama serve

# Отправьте запрос
curl http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"tiny","messages":[{"role":"user","content":"Hello!"}]}'
```

Подробнее в [docs/usage.md](docs/usage.md).

## Лицензия

GPLv2