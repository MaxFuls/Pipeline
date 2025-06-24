#!/bin/bash

# Тестовый скрипт для Pipeline Service

echo "🚀 Тестирование Pipeline Service"
echo "=================================="

# Проверяем, что сервис запущен
echo "📡 Проверка доступности сервиса..."
if curl -s http://localhost:8080 > /dev/null; then
    echo "✅ Сервис доступен на http://localhost:8080"
else
    echo "❌ Сервис недоступен. Убедитесь, что он запущен:"
    echo "   go run cmd/main.go -config=config/config.yaml"
    exit 1
fi

# Тестируем API
echo ""
echo "🧪 Тестирование API..."

# Тест 1: Простое сообщение
echo "📝 Тест 1: Простое сообщение"
curl -X POST http://localhost:8080/submit \
  -H "Content-Type: application/json" \
  -d '{
    "mainMessage": "Привет, мир!",
    "additionalFields": ["validator", "processor"]
  }' | jq '.'

echo ""
echo ""

# Тест 2: Сообщение с дополнительными полями
echo "📝 Тест 2: Сообщение с дополнительными полями"
curl -X POST http://localhost:8080/submit \
  -H "Content-Type: application/json" \
  -d '{
    "mainMessage": "Тестовое сообщение для обработки",
    "additionalFields": ["validator", "processor", "formatter"]
  }' | jq '.'

echo ""
echo ""

# Тест 3: Сообщение без дополнительных полей (использует дефолтные)
echo "📝 Тест 3: Сообщение без дополнительных полей"
curl -X POST http://localhost:8080/submit \
  -H "Content-Type: application/json" \
  -d '{
    "mainMessage": "Сообщение с дефолтными обработчиками"
  }' | jq '.'

echo ""
echo "✅ Тестирование завершено!" 