# Примеры использования API

## Аутентификация

### Регистрация пользователя
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

Ответ:
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "status": "offline",
    "is_active": true,
    "created_at": "2024-01-01T12:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Вход в систему
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "securepassword123"
  }'
```

### Выход из системы
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Управление пользователями

### Получить информацию о текущем пользователе
```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Обновить профиль
```bash
curl -X PUT http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Smith",
    "avatar": "https://example.com/avatar.jpg"
  }'
```

### Обновить статус
```bash
curl -X PUT http://localhost:8080/api/v1/users/status \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "online"
  }'
```

## Управление чатами

### Создать приватный чат
```bash
curl -X POST http://localhost:8080/api/v1/chats \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "private",
    "participants": ["550e8400-e29b-41d4-a716-446655440001"]
  }'
```

### Создать групповой чат
```bash
curl -X POST http://localhost:8080/api/v1/chats \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Team Discussion",
    "description": "Our team chat room",
    "type": "group",
    "participants": [
      "550e8400-e29b-41d4-a716-446655440001",
      "550e8400-e29b-41d4-a716-446655440002"
    ]
  }'
```

### Получить список чатов
```bash
curl -X GET http://localhost:8080/api/v1/chats \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Получить чат по ID
```bash
curl -X GET http://localhost:8080/api/v1/chats/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Сообщения

### Отправить текстовое сообщение
```bash
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "chat_id": "550e8400-e29b-41d4-a716-446655440000",
    "content": "Привет! Как дела?",
    "type": "text"
  }'
```

### Отправить личное сообщение
```bash
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "receiver_id": "550e8400-e29b-41d4-a716-446655440001",
    "content": "Личное сообщение",
    "type": "text"
  }'
```

### Получить сообщения чата
```bash
curl -X GET "http://localhost:8080/api/v1/messages?chat_id=550e8400-e29b-41d4-a716-446655440000&limit=20&offset=0" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Отметить сообщение как прочитанное
```bash
curl -X POST http://localhost:8080/api/v1/messages/550e8400-e29b-41d4-a716-446655440000/read \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Контакты

### Добавить контакт
```bash
curl -X POST http://localhost:8080/api/v1/contacts \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "contact_id": "550e8400-e29b-41d4-a716-446655440001",
    "nickname": "John"
  }'
```

### Получить список контактов
```bash
curl -X GET http://localhost:8080/api/v1/contacts \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Заблокировать контакт
```bash
curl -X PUT http://localhost:8080/api/v1/contacts/550e8400-e29b-41d4-a716-446655440001/block \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Звонки

### Инициировать голосовой звонок
```bash
curl -X POST http://localhost:8080/api/v1/calls \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "callee_id": "550e8400-e29b-41d4-a716-446655440001",
    "type": "voice"
  }'
```

### Инициировать видеозвонок
```bash
curl -X POST http://localhost:8080/api/v1/calls \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "callee_id": "550e8400-e29b-41d4-a716-446655440001",
    "type": "video"
  }'
```

### Принять звонок
```bash
curl -X PUT http://localhost:8080/api/v1/calls/550e8400-e29b-41d4-a716-446655440000/answer \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Отклонить звонок
```bash
curl -X PUT http://localhost:8080/api/v1/calls/550e8400-e29b-41d4-a716-446655440000/reject \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## WebSocket сообщения

### Подключение к WebSocket
```javascript
const token = 'YOUR_JWT_TOKEN';
const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

ws.onopen = function(event) {
    console.log('Connected to WebSocket');
};

ws.onmessage = function(event) {
    const message = JSON.parse(event.data);
    console.log('Received:', message);
};

ws.onclose = function(event) {
    console.log('WebSocket connection closed');
};
```

### Отправить сообщение через WebSocket
```javascript
const message = {
    type: 'chat',
    data: {
        chat_id: '550e8400-e29b-41d4-a716-446655440000',
        content: 'Hello from WebSocket!',
        message_type: 'text'
    }
};

ws.send(JSON.stringify(message));
```

### Индикатор печати
```javascript
const typingMessage = {
    type: 'typing',
    data: {
        chat_id: '550e8400-e29b-41d4-a716-446655440000',
        is_typing: true
    }
};

ws.send(JSON.stringify(typingMessage));
```

### Предложение звонка
```javascript
const callOffer = {
    type: 'call_offer',
    data: {
        target_user_id: '550e8400-e29b-41d4-a716-446655440001',
        call_type: 'video',
        offer: 'SDP_OFFER_DATA'
    }
};

ws.send(JSON.stringify(callOffer));
```

## Загрузка файлов

### Загрузить файл
```bash
curl -X POST http://localhost:8080/api/v1/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@/path/to/your/file.jpg" \
  -F "chat_id=550e8400-e29b-41d4-a716-446655440000"
```

## Коды ошибок

- `400` - Bad Request (неверные параметры)
- `401` - Unauthorized (не авторизован)
- `403` - Forbidden (доступ запрещен)
- `404` - Not Found (ресурс не найден)
- `409` - Conflict (конфликт данных)
- `422` - Unprocessable Entity (ошибка валидации)
- `500` - Internal Server Error (внутренняя ошибка сервера)

## Пример обработки ошибок

```bash
# Попытка входа с неверными данными
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "wronguser",
    "password": "wrongpassword"
  }'

# Ответ:
# HTTP/1.1 401 Unauthorized
# {
#   "error": "invalid credentials"
# }
```

## Примеры интеграции

### Python
```python
import requests
import json

# Вход в систему
response = requests.post('http://localhost:8080/api/v1/auth/login', 
    json={
        'username': 'johndoe',
        'password': 'securepassword123'
    })

token = response.json()['token']

# Отправка сообщения
headers = {'Authorization': f'Bearer {token}'}
requests.post('http://localhost:8080/api/v1/messages',
    headers=headers,
    json={
        'chat_id': '550e8400-e29b-41d4-a716-446655440000',
        'content': 'Hello from Python!',
        'type': 'text'
    })
```

### JavaScript/Node.js
```javascript
const fetch = require('node-fetch');

// Вход в систему
const loginResponse = await fetch('http://localhost:8080/api/v1/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        username: 'johndoe',
        password: 'securepassword123'
    })
});

const { token } = await loginResponse.json();

// Отправка сообщения
const messageResponse = await fetch('http://localhost:8080/api/v1/messages', {
    method: 'POST',
    headers: { 
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
    },
    body: JSON.stringify({
        chat_id: '550e8400-e29b-41d4-a716-446655440000',
        content: 'Hello from JavaScript!',
        type: 'text'
    })
});
```