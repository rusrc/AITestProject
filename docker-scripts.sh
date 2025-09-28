#!/bin/bash

# Скрипт для управления Docker контейнером команды

case "$1" in
    "build")
        echo "🔨 Сборка Docker образа..."
        docker-compose build
        ;;
    "up")
        echo "🚀 Запуск контейнера..."
        docker-compose up -d
        echo "✅ Сервер запущен на http://localhost:5000"
        ;;
    "down")
        echo "🛑 Остановка контейнера..."
        docker-compose down
        ;;
    "restart")
        echo "🔄 Перезапуск контейнера..."
        docker-compose restart
        ;;
    "logs")
        echo "📋 Просмотр логов..."
        docker-compose logs -f
        ;;
    "shell")
        echo "🐚 Подключение к контейнеру..."
        docker-compose exec team-honor-board sh
        ;;
    "clean")
        echo "🧹 Очистка Docker ресурсов..."
        docker-compose down
        docker system prune -f
        ;;
    "status")
        echo "📊 Статус контейнера..."
        docker-compose ps
        ;;
    *)
        echo "🏆 Доска почёта команды - Docker управление"
        echo ""
        echo "Использование: $0 {команда}"
        echo ""
        echo "Команды:"
        echo "  build    - Собрать Docker образ"
        echo "  up       - Запустить контейнер"
        echo "  down     - Остановить контейнер"
        echo "  restart  - Перезапустить контейнер"
        echo "  logs     - Показать логи"
        echo "  shell    - Подключиться к контейнеру"
        echo "  clean    - Очистить Docker ресурсы"
        echo "  status   - Показать статус"
        echo ""
        echo "Примеры:"
        echo "  $0 build && $0 up    # Собрать и запустить"
        echo "  $0 logs              # Просмотр логов"
        echo "  $0 clean             # Полная очистка"
        ;;
esac
