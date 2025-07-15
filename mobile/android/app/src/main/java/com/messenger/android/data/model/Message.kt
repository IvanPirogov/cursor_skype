package com.messenger.android.data.model

import androidx.room.Entity
import androidx.room.PrimaryKey
import kotlinx.parcelize.Parcelize
import android.os.Parcelable

@Entity(tableName = "messages")
@Parcelize
data class Message(
    @PrimaryKey
    val id: Int,
    val chatId: Int,
    val senderId: Int,
    val content: String,
    val type: MessageType = MessageType.TEXT,
    val status: MessageStatus = MessageStatus.SENT,
    val replyToId: Int? = null,
    val isEdited: Boolean = false,
    val createdAt: Long = System.currentTimeMillis(),
    val updatedAt: Long = System.currentTimeMillis(),
    val localId: String? = null,
    val fileUrl: String? = null,
    val fileName: String? = null,
    val fileSize: Long? = null,
    val mimeType: String? = null
) : Parcelable {
    
    val isOwn: Boolean
        get() = senderId == getCurrentUserId() // Нужно реализовать получение текущего пользователя
    
    val formattedTime: String
        get() = formatTime(createdAt)
    
    private fun formatTime(timestamp: Long): String {
        val now = System.currentTimeMillis()
        val diff = now - timestamp
        
        return when {
            diff < 60_000 -> "сейчас"
            diff < 3600_000 -> "${diff / 60_000} мин назад"
            diff < 86400_000 -> "${diff / 3600_000} ч назад"
            else -> {
                val date = java.util.Date(timestamp)
                java.text.SimpleDateFormat("HH:mm", java.util.Locale.getDefault()).format(date)
            }
        }
    }
    
    private fun getCurrentUserId(): Int {
        // TODO: Реализовать получение ID текущего пользователя из DataStore
        return 0
    }
}

enum class MessageType {
    TEXT,
    IMAGE,
    VIDEO,
    AUDIO,
    FILE,
    LOCATION,
    STICKER,
    VOICE
}

enum class MessageStatus {
    SENDING,
    SENT,
    DELIVERED,
    READ,
    FAILED
}