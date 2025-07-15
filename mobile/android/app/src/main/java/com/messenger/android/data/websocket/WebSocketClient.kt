package com.messenger.android.data.websocket

import android.util.Log
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import okhttp3.*
import okio.ByteString
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class WebSocketClient @Inject constructor(
    private val okHttpClient: OkHttpClient
) {
    
    companion object {
        private const val TAG = "WebSocketClient"
    }
    
    private var webSocket: WebSocket? = null
    private val _connectionState = MutableStateFlow(ConnectionState.DISCONNECTED)
    val connectionState: StateFlow<ConnectionState> = _connectionState.asStateFlow()
    
    private val _messages = MutableStateFlow<WebSocketMessage?>(null)
    val messages: StateFlow<WebSocketMessage?> = _messages.asStateFlow()
    
    fun connect(token: String, serverUrl: String) {
        val request = Request.Builder()
            .url("$serverUrl/ws?token=$token")
            .build()
        
        webSocket = okHttpClient.newWebSocket(request, object : WebSocketListener() {
            override fun onOpen(webSocket: WebSocket, response: Response) {
                Log.d(TAG, "WebSocket connected")
                _connectionState.value = ConnectionState.CONNECTED
            }
            
            override fun onMessage(webSocket: WebSocket, text: String) {
                Log.d(TAG, "Received message: $text")
                try {
                    val message = parseMessage(text)
                    _messages.value = message
                } catch (e: Exception) {
                    Log.e(TAG, "Error parsing message", e)
                }
            }
            
            override fun onMessage(webSocket: WebSocket, bytes: ByteString) {
                Log.d(TAG, "Received binary message: ${bytes.hex()}")
            }
            
            override fun onClosing(webSocket: WebSocket, code: Int, reason: String) {
                Log.d(TAG, "WebSocket closing: $code $reason")
                _connectionState.value = ConnectionState.DISCONNECTING
            }
            
            override fun onClosed(webSocket: WebSocket, code: Int, reason: String) {
                Log.d(TAG, "WebSocket closed: $code $reason")
                _connectionState.value = ConnectionState.DISCONNECTED
            }
            
            override fun onFailure(webSocket: WebSocket, t: Throwable, response: Response?) {
                Log.e(TAG, "WebSocket error", t)
                _connectionState.value = ConnectionState.ERROR
            }
        })
    }
    
    fun disconnect() {
        webSocket?.close(1000, "User disconnected")
        webSocket = null
        _connectionState.value = ConnectionState.DISCONNECTED
    }
    
    fun sendMessage(message: WebSocketMessage) {
        val json = createMessageJson(message)
        webSocket?.send(json)
    }
    
    fun sendChatMessage(chatId: Int, content: String, messageType: String = "text") {
        val message = WebSocketMessage(
            type = "chat",
            data = mapOf(
                "chat_id" to chatId,
                "content" to content,
                "message_type" to messageType
            ),
            timestamp = System.currentTimeMillis()
        )
        sendMessage(message)
    }
    
    fun sendTypingIndicator(chatId: Int, isTyping: Boolean) {
        val message = WebSocketMessage(
            type = "typing",
            data = mapOf(
                "chat_id" to chatId,
                "is_typing" to isTyping
            ),
            timestamp = System.currentTimeMillis()
        )
        sendMessage(message)
    }
    
    fun markMessageAsRead(messageId: Int) {
        val message = WebSocketMessage(
            type = "message_read",
            data = mapOf(
                "message_id" to messageId
            ),
            timestamp = System.currentTimeMillis()
        )
        sendMessage(message)
    }
    
    private fun parseMessage(json: String): WebSocketMessage {
        // Здесь должен быть парсер JSON (например, Gson)
        // Для примера возвращаем простую структуру
        return WebSocketMessage(
            type = "unknown",
            data = mapOf("raw" to json),
            timestamp = System.currentTimeMillis()
        )
    }
    
    private fun createMessageJson(message: WebSocketMessage): String {
        // Здесь должен быть JSON сериализатор (например, Gson)
        // Для примера возвращаем простую строку
        return """
            {
                "type": "${message.type}",
                "data": ${message.data},
                "timestamp": ${message.timestamp}
            }
        """.trimIndent()
    }
}

data class WebSocketMessage(
    val type: String,
    val data: Map<String, Any>,
    val timestamp: Long,
    val userId: Int? = null
)

enum class ConnectionState {
    DISCONNECTED,
    CONNECTING,
    CONNECTED,
    DISCONNECTING,
    ERROR
}