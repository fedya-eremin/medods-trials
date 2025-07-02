# MEDDOS TRIALS
## Запуск
Запуск происходит на localhost:8000
```bash
docker compose -f docker-compose.yml up -d
```
## Генерация сваггера
Сваггер расположен на localhost:8000/docs/index.html. В нем сгенерированы сущности и эндпоинты
```bash
# необходимо установить swag
swag init -d ./cmd,./internal/api,./internal/dto
```
Для удобства проверки .env исключен из .gitignore
