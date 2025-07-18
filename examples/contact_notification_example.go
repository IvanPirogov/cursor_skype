package main

import (
	"messenger/internal/websocket"
	"messenger/pkg/models"
	"github.com/google/uuid"
)

// Пример использования функции SendNewContactNotification
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
	
	// Отправляем уведомление
	hub.SendNewContactNotification(userID, contact)
	
	// Пример 2: Использование с минимальными данными
	anotherUserID := uuid.New()
	anotherContactID := uuid.New()
	
	hub.SendNewContactNotificationSimple(
		anotherUserID,
		anotherContactID,
		"another_user",
		"Петр",
		"Петров",
		"Коллега",
	)
}

// Пример интеграции в обработчик добавления контакта
func addContactHandler(hub *websocket.Hub, userID uuid.UUID, contactID uuid.UUID, nickname string) error {
	// Здесь должна быть логика добавления контакта в базу данных
	// ...
	
	// После успешного добавления контакта в БД, отправляем уведомление
	// Получаем данные контакта из БД
	var contact models.Contact
	// db.Preload("Contact").Where("user_id = ? AND contact_id = ?", userID, contactID).First(&contact)
	
	// Отправляем уведомление
	hub.SendNewContactNotification(userID, contact)
	
	return nil
}