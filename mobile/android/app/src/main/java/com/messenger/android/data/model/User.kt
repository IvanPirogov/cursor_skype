package com.messenger.android.data.model

import androidx.room.Entity
import androidx.room.PrimaryKey
import kotlinx.parcelize.Parcelize
import android.os.Parcelable

@Entity(tableName = "users")
@Parcelize
data class User(
    @PrimaryKey
    val id: Int,
    val username: String,
    val email: String,
    val firstName: String? = null,
    val lastName: String? = null,
    val avatar: String? = null,
    val status: String = "offline",
    val lastSeen: Long? = null,
    val isOnline: Boolean = false,
    val createdAt: Long = System.currentTimeMillis(),
    val updatedAt: Long = System.currentTimeMillis()
) : Parcelable {
    
    val displayName: String
        get() = when {
            !firstName.isNullOrBlank() && !lastName.isNullOrBlank() -> "$firstName $lastName"
            !firstName.isNullOrBlank() -> firstName
            !lastName.isNullOrBlank() -> lastName
            else -> username
        }
    
    val initials: String
        get() = when {
            !firstName.isNullOrBlank() && !lastName.isNullOrBlank() -> 
                "${firstName.first().uppercase()}${lastName.first().uppercase()}"
            !firstName.isNullOrBlank() -> firstName.first().uppercase()
            !lastName.isNullOrBlank() -> lastName.first().uppercase()
            else -> username.take(2).uppercase()
        }
}