// API клиент для взаимодействия с backend
class ApiClient {
    constructor(baseUrl = '/api/v1') {
        this.baseUrl = baseUrl;
        this.token = localStorage.getItem('auth_token');
    }

    // Установить токен аутентификации
    setToken(token) {
        this.token = token;
        if (token) {
            localStorage.setItem('auth_token', token);
        } else {
            localStorage.removeItem('auth_token');
        }
    }

    // Получить заголовки запроса
    getHeaders(includeAuth = true) {
        const headers = {
            'Content-Type': 'application/json',
        };

        if (includeAuth && this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        return headers;
    }

    // Базовый метод для HTTP запросов
    async request(endpoint, options = {}) {
        const url = `${this.baseUrl}${endpoint}?_t=${Date.now()}`;
        
        const config = {
            method: 'GET',
            headers: this.getHeaders(options.auth !== false),
            ...options,
        };

        try {
            const response = await fetch(url, config);
            
            if (!response.ok) {
                const error = await response.json().catch(() => ({ error: 'Network error' }));
                throw new Error(error.error || `HTTP ${response.status}`);
            }

            return await response.json();
        } catch (error) {
            console.error('API request failed:', error);
            throw error;
        }
    }

    // GET запрос
    async get(endpoint, options = {}) {
        return this.request(endpoint, { ...options, method: 'GET' });
    }

    // POST запрос
    async post(endpoint, data, options = {}) {
        return this.request(endpoint, {
            ...options,
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    // PUT запрос
    async put(endpoint, data, options = {}) {
        return this.request(endpoint, {
            ...options,
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    // DELETE запрос
    async delete(endpoint, options = {}) {
        return this.request(endpoint, { ...options, method: 'DELETE' });
    }

    // Методы аутентификации
    async register(userData) {
        return this.post('/auth/register', userData, { auth: false });
    }

    async login(credentials) {
        return this.post('/auth/login', credentials, { auth: false });
    }

    async logout() {
        return this.post('/auth/logout', {});
    }

    // Методы пользователей
    async getCurrentUser() {
        const response = await this.get('/users/me');
        console.log('API getCurrentUser response:', response);
        console.log('API getCurrentUser returning:', response.user);
        return response.user;
    }

    async getUsers() {
        const response = await this.get('/users');
        return response.users;
    }

    async getUser(userId) {
        const response = await this.get(`/users/${userId}`);
        return response.user;
    }

    async updateProfile(profileData) {
        return this.put('/users/me', profileData);
    }

    async updateStatus(status) {
        return this.put('/users/status', { status });
    }

    // Методы чатов
    async getChats() {
        const response = await this.get('/chats');
        return response.chats;
    }

    async createChat(chatData) {
        const response = await this.post('/chats', chatData);
        return response.chat;
    }

    async getChat(chatId) {
        const response = await this.get(`/chats/${chatId}`);
        return response.chat;
    }

    async updateChat(chatId, chatData) {
        const response = await this.put(`/chats/${chatId}`, chatData);
        return response.chat;
    }

    async deleteChat(chatId) {
        const response = await this.delete(`/chats/${chatId}`);
        return response;
    }

    async addChatMember(chatId, userId) {
        return this.post(`/chats/${chatId}/members`, { user_id: userId });
    }

    async removeChatMember(chatId, userId) {
        return this.delete(`/chats/${chatId}/members/${userId}`);
    }

    // Методы сообщений
    async getMessages(chatId, limit = 20, offset = 0) {
        const response = await this.get(`/messages?chat_id=${chatId}&limit=${limit}&offset=${offset}`);
        return response.messages;
    }

    async sendMessage(messageData) {
        const response = await this.post('/messages', messageData);
        return response.message;
    }

    async getMessage(messageId) {
        const response = await this.get(`/messages/${messageId}`);
        return response.message;
    }

    async updateMessage(messageId, messageData) {
        const response = await this.put(`/messages/${messageId}`, messageData);
        return response.message;
    }

    async deleteMessage(messageId) {
        const response = await this.delete(`/messages/${messageId}`);
        return response;
    }

    async markMessageAsRead(messageId) {
        const response = await this.post(`/messages/${messageId}/read`);
        return response;
    }

    // Методы контактов
    async getContacts() {
        const response = await this.get('/contacts');
        return response.contacts;
    }

    async addContact(contactData) {
        const response = await this.post('/contacts', contactData);
        return response.contact;
    }

    async removeContact(contactId) {
        return this.delete(`/contacts/${contactId}`);
    }

    async blockContact(contactId) {
        return this.put(`/contacts/${contactId}/block`);
    }

    async unblockContact(contactId) {
        return this.put(`/contacts/${contactId}/unblock`);
    }

    // Методы звонков
    async getCalls() {
        return this.get('/calls');
    }

    async initiateCall(callData) {
        return this.post('/calls', callData);
    }

    async answerCall(callId) {
        return this.put(`/calls/${callId}/answer`);
    }

    async rejectCall(callId) {
        return this.put(`/calls/${callId}/reject`);
    }

    async endCall(callId) {
        return this.put(`/calls/${callId}/end`);
    }

    // Загрузка файлов
    async uploadFile(file, chatId) {
        const formData = new FormData();
        formData.append('file', file);
        formData.append('chat_id', chatId);

        const response = await fetch('/api/v1/upload', {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${this.token}`,
            },
            body: formData,
        });

        if (!response.ok) {
            const error = await response.json().catch(() => ({ error: 'Upload failed' }));
            throw new Error(error.error || `HTTP ${response.status}`);
        }

        return await response.json();
    }

    // Проверка health
    async checkHealth() {
        try {
            const response = await fetch('/health');
            return response.ok;
        } catch (error) {
            return false;
        }
    }
}

// Глобальная инициализация
const api = new ApiClient();

// Экспорт для использования в других модулях
window.api = api;