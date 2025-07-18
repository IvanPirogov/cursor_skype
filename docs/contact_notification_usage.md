# Уведомления о новых контактах

## Обзор

Система WebSocket поддерживает отправку уведомлений о добавлении новых контактов в реальном времени. Это позволяет клиентам получать мгновенные уведомления о новых контактах без необходимости обновления страницы.

## Тип сообщения

Новый тип сообщения: `new_contact`

## Структура сообщения

```json
{
  "type": "new_contact",
  "user_id": "uuid-пользователя",
  "timestamp": 1234567890,
  "data": {
    "contact_id": "uuid-контакта",
    "contact": {
      "id": "uuid-контакта",
      "username": "username_контакта",
      "first_name": "Имя",
      "last_name": "Фамилия",
      "avatar": "путь-к-аватару",
      "status": "online"
    },
    "nickname": "Прозвище контакта"
  }
}
```

## Функции для отправки уведомлений

### 1. SendNewContactNotification

Отправляет уведомление с полным объектом контакта.

```go
func (h *Hub) SendNewContactNotification(userID uuid.UUID, contact models.Contact)
```

**Параметры:**
- `userID` - UUID пользователя, которому отправляется уведомление
- `contact` - Полный объект контакта с загруженными данными пользователя

**Пример использования:**
```go
// Получаем контакт из БД с предзагруженными данными пользователя
var contact models.Contact
db.Preload("Contact").Where("user_id = ? AND contact_id = ?", userID, contactID).First(&contact)

// Отправляем уведомление
hub.SendNewContactNotification(userID, contact)
```

### 2. SendNewContactNotificationSimple

Отправляет уведомление с минимальными данными.

```go
func (h *Hub) SendNewContactNotificationSimple(
    userID uuid.UUID, 
    contactID uuid.UUID, 
    contactUsername string, 
    contactFirstName string, 
    contactLastName string, 
    nickname string
)
```

**Параметры:**
- `userID` - UUID пользователя, которому отправляется уведомление
- `contactID` - UUID контакта
- `contactUsername` - Имя пользователя контакта
- `contactFirstName` - Имя контакта
- `contactLastName` - Фамилия контакта
- `nickname` - Прозвище контакта

**Пример использования:**
```go
hub.SendNewContactNotificationSimple(
    userID,
    contactID,
    "john_doe",
    "John",
    "Doe",
    "Коллега"
)
```

## Интеграция в обработчики

### Пример интеграции в обработчик добавления контакта

```go
func addContactHandler(hub *websocket.Hub, userID uuid.UUID, contactID uuid.UUID, nickname string) error {
    // 1. Добавляем контакт в базу данных
    contact := models.Contact{
        UserID:    userID,
        ContactID: contactID,
        Nickname:  nickname,
    }
    
    if err := db.Create(&contact).Error; err != nil {
        return err
    }
    
    // 2. Загружаем данные пользователя-контакта
    var contactUser models.User
    if err := db.Where("id = ?", contactID).First(&contactUser).Error; err != nil {
        return err
    }
    
    contact.Contact = contactUser
    
    // 3. Отправляем уведомление
    hub.SendNewContactNotification(userID, contact)
    
    return nil
}
```

## Обработка на клиентской стороне

### JavaScript

```javascript
// Подписываемся на события новых контактов
websocketClient.on('new_contact', (data) => {
    console.log('Новый контакт добавлен:', data);
    
    // Обновляем список контактов
    updateContactsList(data.contact);
    
    // Показываем уведомление
    showNotification(`Добавлен новый контакт: ${data.contact.first_name} ${data.contact.last_name}`);
});

function updateContactsList(contact) {
    // Логика обновления списка контактов
    const contactsList = document.getElementById('contacts-list');
    const contactElement = createContactElement(contact);
    contactsList.appendChild(contactElement);
}

function showNotification(message) {
    // Показываем уведомление пользователю
    const notification = document.createElement('div');
    notification.className = 'notification';
    notification.textContent = message;
    document.body.appendChild(notification);
    
    setTimeout(() => {
        notification.remove();
    }, 5000);
}
```

### Android (Kotlin)

```kotlin
// В WebSocket клиенте
websocketClient.on("new_contact") { data ->
    val contact = data.contact
    val nickname = data.nickname
    
    // Обновляем UI в главном потоке
    runOnUiThread {
        updateContactsList(contact)
        showNotification("Добавлен новый контакт: ${contact.firstName} ${contact.lastName}")
    }
}
```

## Обработка ошибок

Функции автоматически обрабатывают ошибки сериализации JSON и логируют их:

```go
data, err := json.Marshal(message)
if err != nil {
    log.Printf("Error marshaling new contact notification: %v", err)
    return
}
```

## Проверка онлайн статуса

Перед отправкой уведомления можно проверить, находится ли пользователь онлайн:

```go
if hub.IsUserOnline(userID) {
    hub.SendNewContactNotification(userID, contact)
} else {
    // Пользователь офлайн, сохраняем уведомление для отправки при подключении
    saveOfflineNotification(userID, notification)
}
```

## Лучшие практики

1. **Всегда проверяйте онлайн статус** перед отправкой уведомлений
2. **Используйте полную версию функции** когда у вас есть доступ к базе данных
3. **Обрабатывайте ошибки** на клиентской стороне
4. **Обновляйте UI** только после успешного получения уведомления
5. **Логируйте события** для отладки