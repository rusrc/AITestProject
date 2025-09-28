# 🏆 Доска почёта команды

Веб-приложение для управления участниками команды и их ачивками, построенное на DaisyUI (фронтенд) и Go (бэкенд).

<img width="1960" height="1030" alt="image" src="https://github.com/user-attachments/assets/1e27140c-9aa9-4c70-8ca4-83fb8f713830" />


## 🚀 Быстрый старт

### Предварительные требования
- Docker и Docker Compose
- Современный браузер

### Установка и запуск с Docker

1. **Клонируйте репозиторий и перейдите в папку проекта:**
   ```bash
   cd team-honor-board
   ```

2. **Соберите и запустите контейнер:**
   ```bash
   # Linux/Mac
   ./docker-scripts.sh build && ./docker-scripts.sh up
   
   # Windows
   docker-scripts.bat build && docker-scripts.bat up
   ```

3. **Откройте браузер и перейдите по адресу:**
   ```
   http://localhost:5000
   ```

### Альтернативный запуск (без Docker)

Если у вас установлен Go 1.21+:

1. **Установите зависимости:**
   ```bash
   go mod tidy
   ```

2. **Запустите сервер:**
   ```bash
   go run main.go
   ```

## 📡 API Endpoints

### Участники

#### GET `/api/members`
Получить список всех участников команды.

**Ответ:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Иван Петров",
      "role": "Frontend Developer",
      "avatar": "uploads/avatars/1234567890.jpg",
      "achievements": [
        {
          "id": 1,
          "member_id": 1,
          "image": "uploads/achievements/1234567891.jpg",
          "category": "positive",
          "created_at": "2024-01-01T12:00:00Z"
        }
      ],
      "created_at": "2024-01-01T12:00:00Z"
    }
  ]
}
```

#### POST `/api/members`
Создать нового участника команды.

**Параметры (multipart/form-data):**
- `name` (string, обязательный) - Имя участника
- `role` (string, обязательный) - Роль участника
- `image` (file, опциональный) - Фото участника

**Ответ:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Иван Петров",
    "role": "Frontend Developer",
    "avatar": "uploads/avatars/1234567890.jpg",
    "achievements": [],
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

### Ачивки

#### POST `/api/achievements`
Добавить ачивку участнику.

**Параметры (multipart/form-data):**
- `member_id` (int, обязательный) - ID участника
- `category` (string, обязательный) - Категория ачивки ("positive" или "negative")
- `image` (file, обязательный) - Изображение ачивки

**Ответ:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "member_id": 1,
    "image": "uploads/achievements/1234567891.jpg",
    "category": "positive",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

## 🎨 Функции фронтенда

- **Тёмная тема по умолчанию** с возможностью переключения на любую из 29 тем DaisyUI
- **Адаптивный дизайн** - корректно работает на всех устройствах
- **Модальные окна** для добавления участников и ачивок
- **Загрузка файлов** с предварительным просмотром
- **Валидация форм** на клиентской стороне
- **Современный UI** с использованием компонентов DaisyUI

## 📁 Структура проекта

```
team-honor-board/
├── assets/
│   └── images/          # SVG иконки для ачивок и аватаров
├── uploads/             # Загруженные пользователями файлы
│   ├── avatars/         # Фото участников
│   └── achievements/    # Изображения ачивок
├── main.go              # Основной файл сервера
├── go.mod              # Зависимости Go
├── index.html          # Фронтенд приложения
└── README.md           # Документация
```

## 🐳 Docker управление

### Основные команды

```bash
# Linux/Mac
./docker-scripts.sh {команда}

# Windows
docker-scripts.bat {команда}
```

### Доступные команды:

- **`build`** - Собрать Docker образ
- **`up`** - Запустить контейнер
- **`down`** - Остановить контейнер
- **`restart`** - Перезапустить контейнер
- **`logs`** - Показать логи
- **`shell`** - Подключиться к контейнеру
- **`clean`** - Очистить Docker ресурсы
- **`status`** - Показать статус

### Примеры использования:

```bash
# Собрать и запустить
./docker-scripts.sh build && ./docker-scripts.sh up

# Просмотр логов
./docker-scripts.sh logs

# Полная очистка
./docker-scripts.sh clean
```

## 🔧 Технические детали

### Бэкенд (Go)
- **Фреймворк:** Gorilla Mux для роутинга
- **Хранилище:** CSV файлы (данные сохраняются между перезапусками)
- **CORS:** Настроен для работы с фронтендом
- **Загрузка файлов:** Поддержка multipart/form-data до 10MB
- **Docker:** Многоэтапная сборка с Alpine Linux

### Фронтенд
- **CSS Framework:** DaisyUI + Tailwind CSS
- **JavaScript:** Vanilla JS (без дополнительных зависимостей)
- **Иконки:** Локальные SVG файлы
- **Темы:** 29 встроенных тем DaisyUI

## 🚨 Ограничения

- Максимальный размер загружаемого файла: 10MB
- Поддерживаются только изображения (jpg, png, gif, svg)
- Данные сохраняются в CSV файлах (требуют монтирования volume в Docker)

## 🔮 Возможные улучшения

- Добавить базу данных (SQLite, PostgreSQL)
- Реализовать аутентификацию и авторизацию
- Добавить пагинацию для больших списков
- Реализовать поиск и фильтрацию участников
- Добавить экспорт данных
- Реализовать уведомления в реальном времени

## 📝 Лицензия

MIT License




