# Примеры использования веб-клиента

## 1. Аутентификация

### Регистрация нового пользователя

```javascript
// Пример регистрации через API
const userData = {
    username: "john_doe",
    email: "john@example.com",
    password: "secure_password123",
    first_name: "John",
    last_name: "Doe"
};

try {
    const response = await api.register(userData);
    console.log('Пользователь зарегистрирован:', response.user);
    // Токен автоматически сохраняется в localStorage
} catch (error) {
    console.error('Ошибка регистрации:', error.message);
}
```

### Вход в систему

```javascript
// Пример входа через API
const credentials = {
    username: "john_doe",
    password: "secure_password123"
};

try {
    const response = await api.login(credentials);
    console.log('Пользователь вошел:', response.user);
    // Перенаправление на страницу чата
    window.location.href = '/chat.html';
} catch (error) {
    console.error('Ошибка входа:', error.message);
}
```

## 2. Управление чатами

### Создание нового чата

```javascript
// Создание группового чата
const chatData = {
    name: "Команда разработчиков",
    type: "group",
    participants: ["user1", "user2", "user3"]
};

try {
    const result = await api.createChat(chatData);
    console.log('Чат создан:', result);
    notifications.success('Чат создан', 'Новый чат успешно создан');
} catch (error) {
    console.error('Ошибка создания чата:', error);
    notifications.error('Ошибка', 'Не удалось создать чат');
}
```

### Получение списка чатов

```javascript
// Загрузка чатов пользователя
try {
    const response = await api.getChats();
    const chats = response.chats || [];
    
    chats.forEach(chat => {
        console.log(`Чат: ${chat.name}, Последнее сообщение: ${chat.last_message}`);
    });
} catch (error) {
    console.error('Ошибка загрузки чатов:', error);
}
```

## 3. Отправка сообщений

### Отправка текстового сообщения

```javascript
// Отправка сообщения через API
const messageData = {
    chat_id: 123,
    content: "Привет, как дела?",
    type: "text"
};

try {
    const result = await api.sendMessage(messageData);
    console.log('Сообщение отправлено:', result);
} catch (error) {
    console.error('Ошибка отправки:', error);
}
```

### Отправка через WebSocket

```javascript
// Real-time отправка сообщения
if (websocketClient && websocketClient.isConnected) {
    websocketClient.sendChatMessage(chatId, "Привет всем!");
}
```

## 4. Загрузка файлов

### Загрузка изображения

```javascript
// Обработка загрузки файла
const fileInput = document.getElementById('file-input');
fileInput.addEventListener('change', async (event) => {
    const files = event.target.files;
    
    for (const file of files) {
        // Проверка типа файла
        if (file.type.startsWith('image/')) {
            try {
                const result = await api.uploadFile(file, currentChatId);
                console.log('Файл загружен:', result);
                notifications.success('Файл загружен', `${file.name} успешно загружен`);
            } catch (error) {
                console.error('Ошибка загрузки:', error);
                notifications.error('Ошибка', `Не удалось загрузить ${file.name}`);
            }
        }
    }
});
```

### Drag & Drop загрузка

```javascript
// Добавление поддержки drag & drop
const messagesContainer = document.getElementById('messages');

messagesContainer.addEventListener('dragover', (e) => {
    e.preventDefault();
    messagesContainer.classList.add('drag-over');
});

messagesContainer.addEventListener('dragleave', () => {
    messagesContainer.classList.remove('drag-over');
});

messagesContainer.addEventListener('drop', async (e) => {
    e.preventDefault();
    messagesContainer.classList.remove('drag-over');
    
    const files = Array.from(e.dataTransfer.files);
    
    for (const file of files) {
        await handleFileUpload(file);
    }
});
```

## 5. WebSocket события

### Обработка входящих сообщений

```javascript
// Инициализация WebSocket
const websocket = initWebSocket();

// Обработка новых сообщений
websocket.on('new_message', (data) => {
    console.log('Новое сообщение:', data);
    
    if (data.chat_id === currentChatId) {
        const messageElement = createMessageElement(data);
        messagesContainer.appendChild(messageElement);
        scrollToBottom();
    }
    
    // Обновление списка чатов
    updateChatLastMessage(data.chat_id, data.content);
});

// Обработка индикатора печати
websocket.on('typing', (data) => {
    if (data.chat_id === currentChatId && data.user_id !== currentUserId) {
        if (data.is_typing) {
            showTypingIndicator(data.user_name);
        } else {
            hideTypingIndicator();
        }
    }
});

// Обработка статуса пользователей
websocket.on('user_status', (data) => {
    updateUserStatus(data.user_id, data.status);
});
```

## 6. Уведомления

### Различные типы уведомлений

```javascript
// Успешное уведомление
notifications.success('Успешно!', 'Операция выполнена успешно');

// Ошибка
notifications.error('Ошибка!', 'Что-то пошло не так');

// Предупреждение
notifications.warning('Внимание!', 'Проверьте введенные данные');

// Информационное
notifications.info('Информация', 'Новая версия доступна');

// Уведомление без автоматического скрытия
const notificationId = notifications.show('info', 'Загрузка...', 'Пожалуйста, подождите', 0);
// Скрыть вручную
notifications.hide(notificationId);
```

## 7. Управление контактами

### Добавление контакта

```javascript
// Добавление нового контакта
const contactData = {
    username: "jane_smith",
    nickname: "Jane"
};

try {
    const result = await api.addContact(contactData);
    console.log('Контакт добавлен:', result);
    notifications.success('Контакт добавлен', 'Новый контакт успешно добавлен');
    
    // Обновление списка контактов
    loadContacts();
} catch (error) {
    console.error('Ошибка добавления контакта:', error);
    notifications.error('Ошибка', 'Не удалось добавить контакт');
}
```

### Поиск контактов

```javascript
// Поиск контактов
const searchInput = document.getElementById('search-input');
searchInput.addEventListener('input', (e) => {
    const query = e.target.value.toLowerCase();
    const contactItems = document.querySelectorAll('.contact-item');
    
    contactItems.forEach(item => {
        const name = item.querySelector('.contact-name').textContent.toLowerCase();
        if (name.includes(query)) {
            item.style.display = 'flex';
        } else {
            item.style.display = 'none';
        }
    });
});
```

## 8. Локальное хранение

### Работа с localStorage

```javascript
// Сохранение данных пользователя
const userInfo = {
    id: 123,
    username: "john_doe",
    email: "john@example.com"
};
localStorage.setItem('user_info', JSON.stringify(userInfo));

// Получение данных пользователя
const savedUser = localStorage.getItem('user_info');
if (savedUser) {
    const user = JSON.parse(savedUser);
    console.log('Пользователь:', user);
}

// Очистка при выходе
function logout() {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user_info');
    window.location.href = '/';
}
```

## 9. Обработка ошибок

### Централизованная обработка ошибок

```javascript
// Обработка ошибок API
api.on('error', (error) => {
    console.error('API Error:', error);
    
    if (error.status === 401) {
        // Токен истек
        notifications.error('Сессия истекла', 'Пожалуйста, войдите снова');
        AuthManager.logout();
    } else if (error.status === 500) {
        notifications.error('Ошибка сервера', 'Попробуйте позже');
    } else {
        notifications.error('Ошибка', error.message);
    }
});

// Обработка ошибок WebSocket
websocket.on('error', (error) => {
    console.error('WebSocket Error:', error);
    notifications.warning('Проблемы с соединением', 'Переподключение...');
});

websocket.on('max_reconnect_attempts_reached', () => {
    notifications.error('Нет соединения', 'Не удается подключиться к серверу');
});
```

## 10. Адаптивность

### Мобильная адаптация

```javascript
// Определение мобильного устройства
function isMobile() {
    return window.innerWidth <= 768;
}

// Адаптация интерфейса
function adaptInterface() {
    const sidebar = document.querySelector('.sidebar');
    const mainChat = document.querySelector('.main-chat');
    
    if (isMobile()) {
        // Мобильная версия
        sidebar.classList.add('mobile');
        mainChat.classList.add('mobile');
    } else {
        // Десктопная версия
        sidebar.classList.remove('mobile');
        mainChat.classList.remove('mobile');
    }
}

// Обработка изменения размера окна
window.addEventListener('resize', adaptInterface);
```

## 11. Кастомизация

### Изменение темы

```javascript
// Переключение темы
function toggleTheme() {
    const body = document.body;
    const currentTheme = body.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    
    body.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
}

// Загрузка сохраненной темы
function loadTheme() {
    const savedTheme = localStorage.getItem('theme') || 'light';
    document.body.setAttribute('data-theme', savedTheme);
}

// Инициализация темы
loadTheme();
```

### Кастомные звуки уведомлений

```javascript
// Воспроизведение звука для уведомлений
function playNotificationSound() {
    const audio = new Audio('/static/sounds/notification.mp3');
    audio.play().catch(e => console.log('Не удалось воспроизвести звук:', e));
}

// Использование при получении сообщения
websocket.on('new_message', (data) => {
    if (data.chat_id !== currentChatId) {
        playNotificationSound();
    }
});
```

Эти примеры показывают основные способы использования веб-клиента мессенджера. Для получения более подробной информации обратитесь к документации API и исходному коду.