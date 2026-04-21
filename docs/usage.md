# Использование iollama.cpp

## Установка

### Из исходников
```bash
git clone https://github.com/yourusername/iollama.cpp
cd iollama.cpp
make build
./bin/iollama serve --model /path/to/model.gguf
```

### Через Docker
```bash
docker build -t iollama:latest .
docker run -p 8080:8080 -e IOLLAMA_MODEL_PATH=/models/model.gguf -v ./models:/models iollama:latest
```

## Быстрый старт

1. Скачайте модель (например, tinyllama):
   ```bash
   ./bin/iollama download TheBloke/TinyLlama-1.1B-GGUF tinyllama-1.1b.Q4_K_M.gguf
   ```
2. Запустите сервер:
   ```bash
   export IOLLAMA_MODEL_PATH=$HOME/.cache/iollama/models/tinyllama-1.1b.Q4_K_M.gguf
   ./bin/iollama serve
   ```
3. Отправьте запрос:
   ```bash
   curl http://localhost:8080/v1/chat/completions \
     -H "Content-Type: application/json" \
     -d '{"model":"tiny","messages":[{"role":"user","content":"Привет!"}]}'
   ```

## Примеры с Python (OpenAI SDK)

```python
from openai import OpenAI

client = OpenAI(base_url="http://localhost:8080/v1", api_key="dummy")
response = client.chat.completions.create(
    model="tiny",
    messages=[{"role": "user", "content": "Расскажи шутку"}],
    stream=True
)
for chunk in response:
    print(chunk.choices[0].delta.content or "", end="")
```

## Переменные окружения

| Переменная | Значение по умолчанию | Описание |
|------------|----------------------|-----------|
| IOLLAMA_MODEL_PATH | - | Путь к GGUF файлу (обязательно) |
| IOLLAMA_HOST | 127.0.0.1 | Хост для API |
| IOLLAMA_PORT | 8080 | Порт |
| IOLLAMA_MAX_TOKENS | 2048 | Максимум токенов генерации |
| IOLLAMA_TEMPERATURE | 0.7 | Температура |

## Выбор модели для слабого железа

- CPU 4GB RAM: `TheBloke/TinyLlama-1.1B-GGUF` (Q4_K_M)
- GPU 4GB VRAM: `TheBloke/Phi-3-mini-4k-instruct-GGUF` (Q4_K_M)
