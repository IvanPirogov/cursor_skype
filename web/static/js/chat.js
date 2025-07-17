// Основной контроллер чата
class ChatController {
    constructor() {
        this.currentUser = null;
        this.currentChat = null;
        this.chats = new Map();
        this.contacts = new Map();
        this.messages = new Map();
        this.websocket = null;
        this.typingTimeout = null;
        this.isTyping = false;

        this.myContacts = new Map();
        this.selectedContactIds = new Set();        
        this.initializeElements();
        this.initializeEventListeners();
        this.checkAuthentication();
    }

    initializeElements() {
        // Основные элементы
        this.chatContainer = document.querySelector('.chat-container');
        this.sidebar = document.querySelector('.sidebar');
        this.mainChat = document.querySelector('.main-chat');
        
        // Элементы пользователя
        this.currentUserName = document.getElementById('current-user-name');
        this.currentUserStatus = document.getElementById('current-user-status');
        
        // Элементы чата
        this.chatItems = document.getElementById('chat-items');
        // this.contactItems = document.getElementById('contact-items');
        this.messages = document.getElementById('messages');
        this.chatTitle = document.getElementById('chat-title');
        this.chatStatus = document.getElementById('chat-status');
        this.onlineCount = document.getElementById('online-count');
        
        // Элементы ввода
        this.messageText = document.getElementById('message-text');
        this.sendBtn = document.getElementById('send-btn');
        this.fileInput = document.getElementById('file-input');
        this.attachFileBtn = document.getElementById('attach-file-btn');
        
        // Элементы поиска
        this.searchInput = document.getElementById('search-input');
        
        // Кнопки
        this.newChatBtn = document.getElementById('new-chat-btn');
        this.addContactBtn = document.getElementById('add-contact-btn');
        this.logoutBtn = document.getElementById('logout-btn');
        this.voiceCallBtn = document.getElementById('voice-call-btn');
        this.videoCallBtn = document.getElementById('video-call-btn');
        
        // Модальные окна
        this.modalOverlay = document.getElementById('modal-overlay');
        this.newChatModal = document.getElementById('new-chat-modal');
        this.addContactModal = document.getElementById('add-contact-modal');
        
        // Индикатор печати
        this.typingIndicator = document.getElementById('typing-indicator');
        this.myContactItems = document.getElementById('my-contact-items');
        this.tabChatsBtn = document.getElementById('tab-chats');
        this.tabMyContactsBtn = document.getElementById('tab-my-contacts');
        this.chatListBlock = document.getElementById('chat-list-block');
        this.myContactsListBlock = document.getElementById('my-contacts-list-block');
    }

    initializeEventListeners() {
        // Отправка сообщения
        this.sendBtn.addEventListener('click', () => this.sendMessage());
        this.messageText.addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                this.sendMessage();
            }
        });

        // Индикатор печати
        this.messageText.addEventListener('input', () => {
            this.handleTypingIndicator();
        });

        // Прикрепление файлов
        this.attachFileBtn.addEventListener('click', () => {
            this.fileInput.click();
        });
        this.fileInput.addEventListener('change', (e) => {
            this.handleFileUpload(e.target.files);
        });

        // Поиск
        this.searchInput.addEventListener('input', (e) => {
            this.handleSearch(e.target.value);
        });

        // Кнопки
        this.newChatBtn.addEventListener('click', () => this.showNewChatModal());
        this.addContactBtn.addEventListener('click', () => this.showAddContactModal());
        this.logoutBtn.addEventListener('click', () => this.logout());
        this.voiceCallBtn.addEventListener('click', () => this.initiateCall('voice'));
        this.videoCallBtn.addEventListener('click', () => this.initiateCall('video'));

         // Вкладки (табы)
        this.tabChatsBtn.addEventListener('click', () => this.showTab('chats'));
        this.tabMyContactsBtn.addEventListener('click', () => this.showTab('myContacts'));

        // Модальные окна
        this.setupModalHandlers();

        // Автоматическое изменение размера textarea
        this.messageText.addEventListener('input', () => {
            this.adjustTextareaHeight();
        });

        // Инициализация состояния поля участников
        this.initializeParticipantsField();
    }

    setupModalHandlers() {
        // Закрытие модальных окон
        this.modalOverlay.addEventListener('click', (e) => {
            if (e.target === this.modalOverlay) {
                this.hideModals();
            }
        });

        // Кнопки закрытия
        document.querySelectorAll('.btn-close').forEach(btn => {
            btn.addEventListener('click', () => this.hideModals());
        });

        // Новый чат
        document.getElementById('create-chat').addEventListener('click', () => {
            this.createNewChat();
        });
        document.getElementById('cancel-new-chat').addEventListener('click', () => {
            this.hideModals();
        });

        // Изменение типа чата
        document.getElementById('chat-type').addEventListener('change', (e) => {
            this.handleChatTypeChange(e.target.value);
        });

        // Добавление контакта
        document.getElementById('add-contact').addEventListener('click', () => {
            this.addNewContact();
        });
        document.getElementById('cancel-add-contact').addEventListener('click', () => {
            this.hideModals();
        });
    }

    async checkAuthentication() {
        const token = localStorage.getItem('auth_token');
        if (!token) {
            window.location.href = '/';
            return;
        }

        try {
           const user = await api.getCurrentUser();
            this.currentUser = user.user;
            this.updateUserInfo();
            this.initializeWebSocket();
            this.loadInitialData();
        } catch (error) {
            console.error('Authentication failed:', error);
            AuthManager.logout();
        }
    }

    updateUserInfo() {
        if (this.currentUser) {
            this.currentUserName.textContent = this.currentUser.username;
            this.currentUserStatus.textContent = 'Онлайн';
        }
    }

    initializeWebSocket() {
        this.websocket = initWebSocket();
        
        // Обработчики WebSocket событий
        this.websocket.on('connected', () => {
            console.log('Connected to WebSocket');
            this.currentUserStatus.textContent = 'Онлайн';
        });

        this.websocket.on('disconnected', () => {
            console.log('Disconnected from WebSocket');
            this.currentUserStatus.textContent = 'Офлайн';
        });

        this.websocket.on('new_message', (data) => {
            this.handleNewMessage(data);
        });

        this.websocket.on('user_status', (data) => {
            this.updateUserStatus(data);
        });

        this.websocket.on('typing', (data) => {
            this.handleTypingIndicator(data);
        });

        this.websocket.on('call_offer', (data) => {
            this.handleCallOffer(data);
        });
    }

    async loadInitialData() {
        try {
            // Загрузка чатов
            const chatsResponse = await api.getChats();
            this.renderChats(chatsResponse.chats || []);

            // Загрузка контактов
            const contactsResponse = await api.getContacts();
            // this.renderContacts(contactsResponse.contacts || []);

            // Обновление счетчика онлайн
            this.updateOnlineCount();
        } catch (error) {
            console.error('Failed to load initial data:', error);
            notifications.error('Ошибка', 'Не удалось загрузить данные');
        }
    }

    renderChats(chats) {
       this.chatItems.innerHTML = '';
        // Собираем контакты с которыми есть переписка
        const myContactsMap = new Map();
        chats.forEach(chat => {
            this.chats.set(chat.id, chat);
            // Добавляем всех участников кроме себя в myContacts
            if (Array.isArray(chat.members)) {
                chat.members.forEach(m => {
                    if (m.user && m.user.id !== this.currentUser?.id) {
                        myContactsMap.set(m.user.id, m.user);
                    }
                });
            }
            const chatElement = this.createChatElement(chat);
            this.chatItems.appendChild(chatElement);
        });
        // Обновляем myContacts
        this.myContacts = myContactsMap;
        // Если открыта вкладка "Твои контакты" — обновить её
        if (this.tabMyContactsBtn.classList.contains('active')) {
            this.renderMyContacts();
        }
    }

    createChatElement(chat) {
        const div = document.createElement('div');
        div.className = 'chat-item';
        div.dataset.chatId = chat.id;
        
        const avatar = this.getAvatarInitials(chat.name || 'Chat');
        const unreadCount = chat.unread_count || 0;
        // Подсчёт участников и онлайн
        let membersCount = 0;
        let onlineCount = 0;
        if (Array.isArray(chat.members)) {
            membersCount = chat.members.length;
            onlineCount = chat.members.filter(m => m.user && m.user.status === 'online').length;
        }
        const membersInfo = `Участников: ${membersCount}, онлайн: ${onlineCount}`;
        const lastTime = chat.last_message_time ? this.formatTime(chat.last_message_time) : '';
        
        div.innerHTML = `
            <div class="chat-avatar">
                ${avatar}
            </div>
            <div class="chat-details">
                <div class="chat-name">${chat.name || 'Unknown'}</div>
                <div class="chat-last-message">${membersInfo}</div>
            </div>
            <div class="chat-meta">
                <div class="chat-time">${lastTime}</div>
                ${unreadCount > 0 ? `<div class="chat-unread">${unreadCount}</div>` : ''}
            </div>
        `;
        
        div.addEventListener('click', () => this.selectChat(chat.id));
        return div;
    }

    // renderContacts(contacts) {
    //     this.contactItems.innerHTML = '';
    //     contacts.forEach(contact => {
    //         this.contacts.set(contact.id, contact);
    //         const contactElement = this.createContactElement(contact);
    //         this.contactItems.appendChild(contactElement);
    //     });
    // }

    createContactElement(contact) {
        const div = document.createElement('div');
        div.className = 'contact-item';
        div.dataset.contactId = contact.id;
        
        const avatar = this.getAvatarInitials(contact.nickname || contact.username);
        
        div.innerHTML = `
            <div class="contact-avatar">
                ${avatar}
            </div>
            <div class="contact-name">${contact.nickname || contact.username}</div>
        `;
        
        div.addEventListener('click', () => this.startPrivateChat(contact.id));
        return div;
    }

    async selectChat(chatId) {
        
        // Убираем выделение с предыдущего чата
        document.querySelectorAll('.chat-item').forEach(item => {
            item.classList.remove('active');
        });
        
        // Выделяем текущий чат
        const chatElement = document.querySelector(`[data-chat-id="${chatId}"]`);
        if (chatElement) {
            chatElement.classList.add('active');
        }
        
        this.currentChat = this.chats.get(chatId);
        if (!this.currentChat) {
            console.error('Chat not found:', chatId);
            return;
        }
        
        // Обновляем заголовок чата
        this.chatTitle.textContent = this.currentChat.name || 'Unknown';
        // this.chatStatus.textContent = this.currentChat.is_online ? 'Онлайн' : 'Офлайн';
        // Устанавливаем статус в зависимости от типа чата
        if (this.currentChat.type === 'public') {
            this.chatStatus.textContent = 'Публичный канал';
        } else if (this.currentChat.type === 'private') {
            this.chatStatus.textContent = 'Приватный чат';
        } else if (this.currentChat.type === 'group') {
            this.chatStatus.textContent = 'Групповой чат';
        } else {
            this.chatStatus.textContent = 'Чат';
        }
        
        // Загружаем сообщения
        try {
            const messagesResponse = await api.getMessages(chatId);
            this.renderMessages(messagesResponse.messages || []);
        } catch (error) {
            console.error('Failed to load messages:', error);
            notifications.error('Ошибка', 'Не удалось загрузить сообщения');
        }
    }

    renderMessages(messages) {
        this.messages.innerHTML = '';
        
        if (messages.length === 0) {
            console.log('No messages to render');
            return;
        }
        
        messages.forEach(message => {
            const messageElement = this.createMessageElement(message);
            this.messages.appendChild(messageElement);
        });
        
        // Прокрутка к последнему сообщению
        this.scrollToBottom();
    }

    createMessageElement(message) {
        const isOwnMessage = message.sender_id === this.currentUser.id;
        const div = document.createElement('div');
        div.className = `message ${isOwnMessage ? 'own' : ''}`;
        div.dataset.messageId = message.id;
        
        const time = this.formatTime(message.created_at);
    
        if (isOwnMessage && (window.innerWidth < 768)) {
            // Сообщение автора на мобильном устройстве - отображаем справа без имени отправителя
            div.innerHTML = `
                <div class="message-content">
                    <p class="message-time">${time}</p>
                    <div class="message-bubble">
                        <p class="message-text">${this.escapeHtml(message.content)}</p>
                        ${message.files ? this.renderMessageFiles(message.files) : ''}
                    </div>
                    <div class="message-info">
                        <div class="message-status">${this.getMessageStatus(message.status)}</div>
                    </div>
                </div>
            `;
        } else {
            // Сообщение другого пользователя - отображаем слева с именем отправителя
            let senderName = 'Unknown';
            let senderInitials = 'U';
            
            if (message.sender) {
                // Если есть объект sender с полной информацией
                senderName = message.sender.username || message.sender.first_name || 'Unknown';
                if (message.sender.first_name && message.sender.last_name) {
                    senderName = `${message.sender.first_name} ${message.sender.last_name}`;
                } else if (message.sender.first_name) {
                    senderName = message.sender.first_name;
                }
                senderInitials = this.getAvatarInitials(senderName);
            } else if (message.sender_name) {
                // Fallback для старых сообщений
                senderName = message.sender_name;
                senderInitials = this.getAvatarInitials(senderName);
            } else {
                // Если нет информации об отправителе
                senderInitials = this.getAvatarInitials('User');
            }
            
            div.innerHTML = `
                <div class="message-avatar">
                    ${senderInitials}
                </div>
                <div class="message-content">
                    <div class="message-info">
                        <span class="message-sender">${senderName}</span>
                        <span class="message-time">${time}</span>
                    </div>
                    <div class="message-bubble">
                        <p class="message-text">${this.escapeHtml(message.content)}</p>
                        ${message.files ? this.renderMessageFiles(message.files) : ''}
                    </div>
                </div>
            `;
        }
        
        return div;
    }

    renderMessageFiles(files) {
        return files.map(file => {
            if (file.mime_type.startsWith('image/')) {
                return `<img src="${file.file_path}" alt="${file.file_name}" class="message-image">`;
            } else {
                return `
                    <div class="message-file">
                        <div class="file-icon">
                            <i class="fas fa-file"></i>
                        </div>
                        <div class="file-info">
                            <div class="file-name">${file.file_name}</div>
                            <div class="file-size">${this.formatFileSize(file.file_size)}</div>
                        </div>
                    </div>
                `;
            }
        }).join('');
    }

    async sendMessage() {
        const content = this.messageText.value.trim();
        if (!content || !this.currentChat) return;
        
        // Очищаем поле ввода
        this.messageText.value = '';
        this.adjustTextareaHeight();
        
        // Сначала отправляем через API для сохранения в базе
        try {
            const response = await api.sendMessage({
                chat_id: this.currentChat.id,
                content: content,
                type: 'text'
            });
            
            // Если сообщение успешно сохранено, отправляем через WebSocket
            if (response && response.message) {
                this.websocket.sendChatMessage(this.currentChat.id, content);
                
                // Добавляем сообщение в UI
                const messageElement = this.createMessageElement(response.message);
                this.messages.appendChild(messageElement);
                this.scrollToBottom();
                
                // Обновляем список чатов
                this.updateChatLastMessage(this.currentChat.id, content);
            }
        } catch (error) {
            console.error('Failed to send message:', error);
            notifications.error('Ошибка', 'Не удалось отправить сообщение');
        }
    }

    handleNewMessage(data) {
        if (data.chat_id === this.currentChat?.id) {
            const messageElement = this.createMessageElement(data);
            this.messages.appendChild(messageElement);
            this.scrollToBottom();
        } 
        // Обновляем список чатов
        this.updateChatLastMessage(data.chat_id, data.content);
    }

    handleTypingIndicator() {
        if (!this.currentChat) return;
        
        if (!this.isTyping) {
            this.isTyping = true;
            this.websocket.sendTypingIndicator(this.currentChat.id, true);
        }
        
        // Сброс таймера
        clearTimeout(this.typingTimeout);
        this.typingTimeout = setTimeout(() => {
            this.isTyping = false;
            this.websocket.sendTypingIndicator(this.currentChat.id, false);
        }, 1000);
    }

    handleTypingIndicatorReceived(data) {
        if (data.chat_id === this.currentChat?.id && data.user_id !== this.currentUser.id) {
            if (data.is_typing) {
                this.typingIndicator.classList.remove('hidden');
            } else {
                this.typingIndicator.classList.add('hidden');
            }
        }
    }

    async handleFileUpload(files) {
        if (!this.currentChat) return;
        
        for (const file of files) {
            try {
                const result = await api.uploadFile(file, this.currentChat.id);
                notifications.success('Файл загружен', `${file.name} успешно загружен`);
                
                // Обновляем сообщения
                this.selectChat(this.currentChat.id);
            } catch (error) {
                console.error('File upload failed:', error);
                notifications.error('Ошибка загрузки', `Не удалось загрузить ${file.name}`);
            }
        }
    }

    showNewChatModal() {
        this.modalOverlay.classList.remove('hidden');
        this.newChatModal.classList.remove('hidden');
        this.newChatModal.style.display = 'block';
        // Устанавливаем публичный чат по умолчанию и скрываем участников
        document.getElementById('chat-type').value = 'public';
        document.getElementById('chat-name').value = '';
        document.getElementById('chat-participants').value = '';
        this.handleChatTypeChange('public');
    }

    initializeParticipantsField() {
        // Устанавливаем начальное состояние поля участников
        const chatTypeSelect = document.getElementById('chat-type');
        if (chatTypeSelect) {
            this.handleChatTypeChange(chatTypeSelect.value);
        }
    }

    handleChatTypeChange(chatType) {
        const participantsGroup = document.getElementById('participants-group');
        if (participantsGroup) {
            if (chatType === 'public') {
                participantsGroup.classList.remove('show');
            } else {
                participantsGroup.classList.add('show');
            }
        }
    }

    showAddContactModal() {
        this.modalOverlay.classList.remove('hidden');
        this.addContactModal.classList.remove('hidden');
        this.addContactModal.style.display = 'block';
        // Очищаем поля при открытии
        document.getElementById('contact-username').value = '';
        document.getElementById('contact-nickname').value = '';
        // Новый функционал: поиск и выбор пользователя
        this.loadAndShowUserSearchList();
    }

    async loadAndShowUserSearchList() {
        const userListDiv = document.getElementById('user-search-list');
        userListDiv.innerHTML = '<div style="padding:8px;color:#888;">Загрузка...</div>';
        try {
            if (!this.allUsers) {
                const users = await api.getUsers();
                this.allUsers = users;
            }
            this.renderUserSearchList(this.allUsers);
            // Поиск по вводу
            const usernameInput = document.getElementById('contact-username');
            usernameInput.addEventListener('input', () => {
                const val = usernameInput.value.toLowerCase();
                const filtered = this.allUsers.filter(u =>
                    u.username.toLowerCase().includes(val) ||
                    (u.nickname && u.nickname.toLowerCase().includes(val))
                );
                this.renderUserSearchList(filtered);
            });
        } catch (e) {
            userListDiv.innerHTML = '<div style="padding:8px;color:#c00;">Ошибка загрузки пользователей</div>';
        }
    }

    renderUserSearchList(users) {
        const userListDiv = document.getElementById('user-search-list');
        if (!users.length) {
            userListDiv.innerHTML = '<div style="padding:8px;color:#888;">Пользователи не найдены</div>';
            return;
        }
        userListDiv.innerHTML = users.map(u =>
            `<div class="user-search-item" data-username="${u.username}">
                <span class="user-search-name">${u.nickname ? u.nickname + ' (' + u.username + ')' : u.username}</span>
            </div>`
        ).join('');
        // Клик по пользователю
        userListDiv.querySelectorAll('.user-search-item').forEach(item => {
            item.addEventListener('click', () => {
                document.getElementById('contact-username').value = item.dataset.username;
            });
        });
    }

    hideModals() {
        this.modalOverlay.classList.add('hidden');
        this.addContactModal.classList.add('hidden');
        this.newChatModal.classList.add('hidden');
        this.newChatModal.style.display = 'none';
        this.addContactModal.style.display = 'none';
    }

    async createNewChat() {
        const chatType = document.getElementById('chat-type').value;
        const chatName = document.getElementById('chat-name').value.trim();
        const participants = document.getElementById('chat-participants').value.trim();
        
        if (!chatName) {
            notifications.error('Ошибка', 'Введите название чата');
            return;
        }
        
        try {
            const chatData = {
                name: chatName,
                type: chatType
            };
            
            // Добавляем участников только для приватных и групповых чатов
            if (chatType !== 'public' && participants) {
                chatData.member_ids = participants.split(',').map(p => p.trim()).filter(p => p);
            }
            
            const result = await api.createChat(chatData);
            
            notifications.success('Чат создан', 'Новый чат успешно создан');
            this.hideModals();
            this.loadInitialData();
        } catch (error) {
            console.error('Failed to create chat:', error);
            notifications.error('Ошибка', 'Не удалось создать чат');
        }
    }

    async addNewContact() {
        const username = document.getElementById('contact-username').value.trim();
        const nickname = document.getElementById('contact-nickname').value.trim();
        
        if (!username) {
            notifications.error('Ошибка', 'Введите имя пользователя');
            return;
        }
        
        try {
            const result = await api.addContact({
                username: username,
                nickname: nickname
            });
            
            notifications.success('Контакт добавлен', 'Новый контакт успешно добавлен');
            this.hideModals();
            this.loadInitialData();
        } catch (error) {
            console.error('Failed to add contact:', error);
            notifications.error('Ошибка', 'Не удалось добавить контакт');
        }
    }

    initiateCall(type) {
        if (!this.currentChat) return;
        
        notifications.info('Звонок', `Инициирую ${type === 'voice' ? 'голосовой' : 'видео'} звонок`);
        
        // TODO: Implement call functionality
        // this.websocket.sendCallOffer(targetUserId, type, offer);
    }

    handleCallOffer(data) {
        notifications.info('Входящий звонок', `${data.caller_name} звонит вам`);
        
        // TODO: Show call interface
    }

    adjustTextareaHeight() {
        this.messageText.style.height = 'auto';
        this.messageText.style.height = Math.min(this.messageText.scrollHeight, 100) + 'px';
    }

    scrollToBottom() {
        this.messages.scrollTop = this.messages.scrollHeight;
    }

    updateChatLastMessage(chatId, content) {
        const chatElement = document.querySelector(`[data-chat-id="${chatId}"]`);
        if (chatElement) {
            const lastMessageElement = chatElement.querySelector('.chat-last-message');
            if (lastMessageElement) {
                lastMessageElement.textContent = content;
            }
        }
    }

    updateUserStatus(data) {
        // TODO: Update user status in UI
        console.log('User status updated:', data);
    }

    updateOnlineCount() {
        // TODO: Update online count
        this.onlineCount.textContent = '5 онлайн';
    }

    handleSearch(query) {
        // TODO: Implement search functionality
        console.log('Search:', query);
    }

    logout() {
        if (this.websocket) {
            this.websocket.close();
        }
        // AuthManager.logout();
        localStorage.removeItem('auth_token');
        localStorage.removeItem('user_info');
        api.setToken(null);
        window.location.href = '/';
    }

    showTab(tab) {
        if (tab === 'chats') {
            this.tabChatsBtn.classList.add('active');
            this.tabMyContactsBtn.classList.remove('active');
            this.chatListBlock.style.display = '';
            this.myContactsListBlock.style.display = 'none';
        } else {
            this.tabChatsBtn.classList.remove('active');
            this.tabMyContactsBtn.classList.add('active');
            this.chatListBlock.style.display = 'none';
            this.myContactsListBlock.style.display = '';
            this.renderMyContacts();
        }
    }

    renderMyContacts() {
        this.myContactItems.innerHTML = '';
        // Собираем всех: с кем есть чат + кого выбрали вручную
        const allContacts = new Map([...this.myContacts]);
        this.selectedContactIds.forEach(id => {
            if (this.contacts.has(id)) {
                allContacts.set(id, this.contacts.get(id));
            }
        });
        allContacts.forEach(contact => {
            const div = document.createElement('div');
            div.className = 'contact-item';
            div.dataset.contactId = contact.id;
            const avatar = this.getAvatarInitials(contact.nickname || contact.username || contact.first_name || contact.username);
            div.innerHTML = `
                <div class="contact-avatar">${avatar}</div>
                <div class="contact-name">${contact.nickname || contact.username || contact.first_name || contact.username}</div>
            `;
            div.addEventListener('click', () => this.startPrivateChat(contact.id));
            this.myContactItems.appendChild(div);
        });
    }

    // Утилиты
    getAvatarInitials(name) {
        return name.split(' ').map(n => n[0]).join('').toUpperCase().substr(0, 2);
    }

    formatTime(timestamp) {
        const date = new Date(timestamp);
        const now = new Date();

        const pad = (n) => n.toString().padStart(2, '0');

        const isToday =
            date.getFullYear() === now.getFullYear() &&
            date.getMonth() === now.getMonth() &&
            date.getDate() === now.getDate();
        
        const hours = pad(date.getHours());
        const minutes = pad(date.getMinutes());

        if (isToday) {
            return `${hours}:${minutes}`;
        } else {
            const year = date.getFullYear();
            const month = pad(date.getMonth() + 1);
            const day = pad(date.getDate());
            return `${year}-${month}-${day} ${hours}:${minutes}`;
        }
    }

    formatFileSize(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

    getMessageStatus(status) {
        const statuses = {
            'sent': 'Отправлено',
            'delivered': 'Доставлено',
            'read': 'Прочитано',
            'failed': 'Ошибка'
        };
        return statuses[status] || 'Неизвестно';
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
    new ChatController();
});

// Экспорт для использования в других модулях
window.ChatController = ChatController;