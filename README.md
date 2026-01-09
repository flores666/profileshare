Проект состоит из 4 микросервисов:

- **auth** — сервис авторизации (JWT, refresh tokens)
- **authOrchestrator** — сервис оркестрации сообщений auth
- **content** — сервис работы с контентом (CRUD)
- **mailer** — сервис отправки писем

---

## 1. Сервисы и эндпоинты

### 1.1 Auth Service

| Метод | Эндпоинт | Описание | Требует авторизации |
|-------|----------|----------|------------------|
| POST | /auth/login | Логин пользователя | ❌ |
| POST | /auth/refresh | Обновление access и refresh токенов | ❌ |
| POST | /auth/logout | Выход пользователя, инвалидирует refresh токен | ✅ |
| POST | /auth/confirm | Подтверждение аккаунта пользователя | ❌ |

Чтобы зарегистрироваться без сервиса для отправки почты:
1. /auth/register
2. в бд находим код и отправляем запрос на /auth/confirm

и эндпоинты /users/ для получения информации о пользователях

> Все защищённые эндпоинты требуют `Authorization: Bearer <access_token>`  

---

### 1.2 Content Service

| Метод | Эндпоинт | Описание | Требует авторизации |
|-------|----------|----------|------------------|
| GET | /content | Получить список контента | ❌ |
| GET | /content/{id} | Получить запись по ID | ❌ |
| POST | /content | Создать запись | ✅ |
| PUT | /content | Обновить запись | ✅ (только владелец) |
| DELETE | /content/{id} | Удалить запись | ✅ (только владелец) |

> Все write-операции требуют JWT access token и проверки владельца записи.

---

### 1.3 Mailer Service

Работает через сообщения Kafka

---

## 2. Переменные окружения

Ниже приведены переменные для каждого сервиса. Обязательно создайте `.env` в корне сервиса.

### 2.1 Auth Service

```env
CONFIG_PATH=config/local.yaml
DB__CONNECTION_STRING=""
SECURITY__ACCESS_SECRET="BkHMGL5ZiN4kotSDzjG8J14adEkAbwiaOB31QzXB21"
SECURITY__ACCESS_LIFETIME_MINUTES=10
SECURITY__REFRESH_LIFETIME_DAYS=7
```

### 2.2 AuthOrchestrator Service

```env
CONFIG_PATH=config/local.yaml
DB__CONNECTION_STRING=""
```

### 2.3 Content Service

```env
CONFIG_PATH=config/local.yaml
DB__CONNECTION_STRING=""
SECURITY__ACCESS_SECRET="BkHMGL5ZiN4kotSDzjG8J14adEkAbwiaOB31QzXB21"
```

### 2.4 Mailer Service

```env
CONFIG_PATH=config/local.yaml
DB__CONNECTION_STRING=""
SMTP__USERNAME=""
SMTP__PASSWORD=""
SMTP__HOST=""
```

---

## 3. Как запустить сервисы

В корне проекта выполнить docker-compose up --build -d
