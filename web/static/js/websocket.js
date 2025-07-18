// WebSocket клиент для real-time сообщений
class WebSocketClient {
    constructor() {
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectDelay = 1000;
        this.isConnected = false;
        this.messageQueue = [];
        this.eventHandlers = new Map();
        
        this.connect();
    }

    connect() {
        const token = localStorage.getItem('auth_token');
        if (!token) {
            console.error('No auth token found');
            return;
        }

        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws?token=${token}`;

        try {
            this.ws = new WebSocket(wsUrl);
            this.setupEventHandlers();
        } catch (error) {
            console.error('WebSocket connection error:', error);
            this.scheduleReconnect();
        }
    }

    setupEventHandlers() {
        this.ws.onopen = () => {
            console.log('WebSocket connected');
            this.isConnected = true;
            this.reconnectAttempts = 0;
            
            // Отправляем накопленные сообщения
            this.flushMessageQueue();
            
            // Уведомляем о подключении
            this.emit('connected');
        };

        this.ws.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                this.handleMessage(message);
            } catch (error) {
                console.error('Error parsing WebSocket message:', error);
            }
        };

        this.ws.onclose = (event) => {
            console.log('WebSocket disconnected:', event.code, event.reason);
            this.isConnected = false;
            this.emit('disconnected');
            
            // Переподключение только если это не преднамеренное закрытие
            if (!event.wasClean) {
                this.scheduleReconnect();
            }
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            this.emit('error', error);
        };
    }

    handleMessage(message) {
        console.log('WebSocket message received:', message);
        const { type, data, user_id, timestamp } = message;
        
        switch (type) {
            case 'chat':
                console.log('Chat message:', data);
                this.emit('message', { ...data, user_id, timestamp });
                break;
                
            case 'user_status':
                console.log('User status:', data);
                this.emit('user_status', data);
                break;
                
            case 'typing':
                console.log('Typing indicator:', data);
                this.emit('typing', { ...data, user_id });
                break;
                
            case 'call_offer':
                this.emit('call_offer', { ...data, user_id });
                break;
                
            case 'call_answer':
                this.emit('call_answer', { ...data, user_id });
                break;
                
            case 'call_reject':
                this.emit('call_reject', { ...data, user_id });
                break;
                
            case 'call_end':
                this.emit('call_end', { ...data, user_id });
                break;
                
            case 'new_message':
                console.log('New message:', data);
                this.emit('new_message', { ...data, user_id, timestamp });
                break;
                
            case 'message_read':
                this.emit('message_read', { ...data, user_id });
                break;
                
            case 'user_joined':
                this.emit('user_joined', { ...data, user_id });
                break;
                
            case 'user_left':
                this.emit('user_left', { ...data, user_id });
                break;
                
            case 'new_contact':
                console.log('New contact notification:', data);
                this.emit('new_contact', { ...data, user_id, timestamp });
                break;
                
            default:
                console.warn('Unknown message type:', type);
        }
    }

    send(type, data) {
        const message = {
            type,
            data,
            timestamp: Date.now()
        };

        if (this.isConnected) {
            this.ws.send(JSON.stringify(message));
        } else {
            // Добавляем в очередь для отправки при подключении
            this.messageQueue.push(message);
        }
    }

    flushMessageQueue() {
        while (this.messageQueue.length > 0) {
            const message = this.messageQueue.shift();
            this.ws.send(JSON.stringify(message));
        }
    }

    scheduleReconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);
            
            console.log(`Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`);
            
            setTimeout(() => {
                this.connect();
            }, delay);
        } else {
            console.error('Max reconnect attempts reached');
            this.emit('max_reconnect_attempts_reached');
        }
    }

    // Методы для отправки различных типов сообщений
    sendChatMessage(chatId, content, messageType = 'text') {
        this.send('chat', {
            chat_id: chatId,
            content,
            message_type: messageType
        });
    }

    sendTypingIndicator(chatId, isTyping) {
        this.send('typing', {
            chat_id: chatId,
            is_typing: isTyping
        });
    }

    sendCallOffer(targetUserId, callType, offer) {
        this.send('call_offer', {
            target_user_id: targetUserId,
            call_type: callType,
            offer
        });
    }

    sendCallAnswer(targetUserId, answer) {
        this.send('call_answer', {
            target_user_id: targetUserId,
            answer
        });
    }

    sendCallReject(targetUserId) {
        this.send('call_reject', {
            target_user_id: targetUserId
        });
    }

    sendCallEnd(targetUserId) {
        this.send('call_end', {
            target_user_id: targetUserId
        });
    }

    markMessageAsRead(messageId) {
        this.send('message_read', {
            message_id: messageId
        });
    }

    // Система событий
    on(event, handler) {
        if (!this.eventHandlers.has(event)) {
            this.eventHandlers.set(event, []);
        }
        this.eventHandlers.get(event).push(handler);
    }

    off(event, handler) {
        if (this.eventHandlers.has(event)) {
            const handlers = this.eventHandlers.get(event);
            const index = handlers.indexOf(handler);
            if (index > -1) {
                handlers.splice(index, 1);
            }
        }
    }

    emit(event, data) {
        if (this.eventHandlers.has(event)) {
            this.eventHandlers.get(event).forEach(handler => {
                try {
                    handler(data);
                } catch (error) {
                    console.error('Error in event handler:', error);
                }
            });
        }
    }

    // Состояние подключения
    isConnected() {
        return this.isConnected;
    }

    // Закрытие соединения
    close() {
        if (this.ws) {
            this.ws.close(1000, 'Client disconnecting');
            this.isConnected = false;
        }
    }

    // Получение информации о соединении
    getConnectionInfo() {
        return {
            isConnected: this.isConnected,
            reconnectAttempts: this.reconnectAttempts,
            messageQueueLength: this.messageQueue.length,
            readyState: this.ws ? this.ws.readyState : null
        };
    }
}

// Глобальная инициализация
let websocketClient = null;

// Функция для инициализации WebSocket клиента
function initWebSocket() {
    if (websocketClient) {
        websocketClient.close();
    }
    websocketClient = new WebSocketClient();
    return websocketClient;
}

// Экспорт для использования в других модулях
window.initWebSocket = initWebSocket;
window.WebSocketClient = WebSocketClient;