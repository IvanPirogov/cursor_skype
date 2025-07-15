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
        const url = `${this.baseUrl}${endpoint}`;
        
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
        return this.get('/users/me');
    }

    async getUsers() {
        return this.get('/users');
    }

    async getUser(userId) {
        return this.get(`/users/${userId}`);
    }

    async updateProfile(profileData) {
        return this.put('/users/me', profileData);
    }

    async updateStatus(status) {
        return this.put('/users/status', { status });
    }

    // Методы чатов
    async getChats() {
        return this.get('/chats');
    }

    async createChat(chatData) {
        return this.post('/chats', chatData);
    }

    async getChat(chatId) {
        return this.get(`/chats/${chatId}`);
    }

    async updateChat(chatId, chatData) {
        return this.put(`/chats/${chatId}`, chatData);
    }

    async deleteChat(chatId) {
        return this.delete(`/chats/${chatId}`);
    }

    async addChatMember(chatId, userId) {
        return this.post(`/chats/${chatId}/members`, { user_id: userId });
    }

    async removeChatMember(chatId, userId) {
        return this.delete(`/chats/${chatId}/members/${userId}`);
    }

    // Методы сообщений
    async getMessages(chatId, limit = 20, offset = 0) {
        return this.get(`/messages?chat_id=${chatId}&limit=${limit}&offset=${offset}`);
    }

    async sendMessage(messageData) {
        return this.post('/messages', messageData);
    }

    async getMessage(messageId) {
        return this.get(`/messages/${messageId}`);
    }

    async updateMessage(messageId, messageData) {
        return this.put(`/messages/${messageId}`, messageData);
    }

    async deleteMessage(messageId) {
        return this.delete(`/messages/${messageId}`);
    }

    async markMessageAsRead(messageId) {
        return this.post(`/messages/${messageId}/read`);
    }

    // Методы контактов
    async getContacts() {
        return this.get('/contacts');
    }

    async addContact(contactData) {
        return this.post('/contacts', contactData);
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