## Как запустить

### Docker
```bash
# Сборка образа
docker build -t monolith .

# Запуск контейнера
docker run -p 8080:8080 monolith
```

Сервис будет доступен по адресу:
- monolith: http://localhost:8080
