package main

import (
	"encoding/json"
	"messenger/internal/websocket"
	"messenger/pkg/models"
	"github.com/google/uuid"
)

// Функция для создания сообщения о новом контакте
func createNewContactMessage(userID uuid.UUID, contact models.Contact) ([]byte, error) {
	message := websocket.Message{
		Type:      websocket.MessageTypeNewContact,
		UserID:    userID,
		Timestamp: 0, // Используйте реальный timestamp
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

// Функция для создания простого сообщения о новом контакте
func createSimpleNewContactMessage(userID uuid.UUID, contactID uuid.UUID, contactUsername string, contactFirstName string, contactLastName string, nickname string) ([]byte, error) {
	message := websocket.Message{
		Type:      websocket.MessageTypeNewContact,
		UserID:    userID,
		Timestamp: 0, // Используйте реальный timestamp
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

// Пример использования функции SendToUser для отправки уведомления о новом контакте
func exampleUsage(hub *websocket.Hub) {
	// Пример 1: Использование с полным объектом Contact
	userID := uuid.New()
	contactID := uuid.New()
	
	// Создаем объект контакта (обычно получается из базы данных)
	contact := models.Contact{
		ID:        uuid.New(),
		UserID:    userID,
		ContactID: contactID,
		Nickname:  "Друг",
		Contact: models.User{
			ID:        contactID,
			Username:  "friend_user",
			FirstName: "Иван",
			LastName:  "Иванов",
			Avatar:    "avatar.jpg",
			Status:    models.StatusOnline,
		},
	}
	
	// Создаем сообщение
	message, err := createNewContactMessage(userID, contact)
	if err != nil {
		// Обработка ошибки
		return
	}
	
	// Отправляем уведомление используя существующую функцию SendToUser
	hub.SendToUser(userID, message)
	
	// Пример 2: Использование с минимальными данными
	anotherUserID := uuid.New()
	anotherContactID := uuid.New()
	
	simpleMessage, err := createSimpleNewContactMessage(
		anotherUserID,
		anotherContactID,
		"another_user",
		"Петр",
		"Петров",
		"Коллега",
	)
	if err != nil {
		// Обработка ошибки
		return
	}
	
	hub.SendToUser(anotherUserID, simpleMessage)
}

// Пример интеграции в обработчик добавления контакта
func addContactHandler(hub *websocket.Hub, userID uuid.UUID, contactID uuid.UUID, nickname string) error {
	// Здесь должна быть логика добавления контакта в базу данных
	// ...
	
	// После успешного добавления контакта в БД, отправляем уведомление
	// Получаем данные контакта из БД
	var contact models.Contact
	// db.Preload("Contact").Where("user_id = ? AND contact_id = ?", userID, contactID).First(&contact)
	
	// Создаем сообщение
	message, err := createNewContactMessage(userID, contact)
	if err != nil {
		return err
	}
	
	// Отправляем уведомление используя существующую функцию SendToUser
	hub.SendToUser(userID, message)
	
	return nil
}