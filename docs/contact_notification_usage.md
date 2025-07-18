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

## Использование существующей функции SendToUser

Для отправки уведомлений используется существующая функция `SendToUser`:

```go
func (h *Hub) SendToUser(userID uuid.UUID, message []byte)
```

## Создание сообщения о новом контакте

### 1. С полным объектом Contact

```go
func createNewContactMessage(userID uuid.UUID, contact models.Contact) ([]byte, error) {
	message := websocket.Message{
		Type:      websocket.MessageTypeNewContact,
		UserID:    userID,
		Timestamp: time.Now().Unix(),
		Data: map[string]interface{}{
			"contact_id": contact.ContactID,
			"contact": map[string]interface{}{
				"id":         contact.Contact.ID,
				"username":   contact.Contact.Username,
				"first_name": contact.Contact.FirstName,
				"last_name":  contact.Contact.LastName,
				"avatar":     contact.Contact.Avatar,
				"status":     contact.Contact.Status,
			},
			"nickname": contact.Nickname,
		},
	}
	
	return json.Marshal(message)
}
```

### 2. С минимальными данными

```go
func createSimpleNewContactMessage(userID uuid.UUID, contactID uuid.UUID, contactUsername string, contactFirstName string, contactLastName string, nickname string) ([]byte, error) {
	message := websocket.Message{
		Type:      websocket.MessageTypeNewContact,
		UserID:    userID,
		Timestamp: time.Now().Unix(),
		Data: map[string]interface{}{
			"contact_id": contactID,
			"contact": map[string]interface{}{
				"id":         contactID,
				"username":   contactUsername,
				"first_name": contactFirstName,
				"last_name":  contactLastName,
			},
			"nickname": nickname,
		},
	}
	
	return json.Marshal(message)
}
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
    
    // 3. Создаем сообщение
    message, err := createNewContactMessage(userID, contact)
    if err != nil {
        return err
    }
    
    // 4. Отправляем уведомление используя существующую функцию SendToUser
    hub.SendToUser(userID, message)
    
    return nil
}
```

### Простой пример использования

```go
// Создаем сообщение с минимальными данными
message, err := createSimpleNewContactMessage(
    userID,
    contactID,
    "john_doe",
    "John",
    "Doe",
    "Коллега",
)
if err != nil {
    // Обработка ошибки
    return err
}

// Отправляем уведомление
hub.SendToUser(userID, message)
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