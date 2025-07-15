package com.messenger.android.data.api

import com.messenger.android.data.model.User
import com.messenger.android.data.model.Message
import com.messenger.android.data.model.Chat
import retrofit2.Response
import retrofit2.http.*

interface ApiService {
    
    // Аутентификация
    @POST("auth/register")
    suspend fun register(@Body request: RegisterRequest): Response<AuthResponse>
    
    @POST("auth/login")
    suspend fun login(@Body request: LoginRequest): Response<AuthResponse>
    
    @POST("auth/logout")
    suspend fun logout(): Response<Unit>
    
    // Пользователи
    @GET("users/me")
    suspend fun getCurrentUser(): Response<User>
    
    @GET("users")
    suspend fun getUsers(): Response<List<User>>
    
    @GET("users/{id}")
    suspend fun getUser(@Path("id") id: Int): Response<User>
    
    @PUT("users/me")
    suspend fun updateProfile(@Body request: UpdateProfileRequest): Response<User>
    
    @PUT("users/status")
    suspend fun updateStatus(@Body request: UpdateStatusRequest): Response<Unit>
    
    // Чаты
    @GET("chats")
    suspend fun getChats(): Response<List<Chat>>
    
    @POST("chats")
    suspend fun createChat(@Body request: CreateChatRequest): Response<Chat>
    
    @GET("chats/{id}")
    suspend fun getChat(@Path("id") id: Int): Response<Chat>
    
    @PUT("chats/{id}")
    suspend fun updateChat(@Path("id") id: Int, @Body request: UpdateChatRequest): Response<Chat>
    
    @DELETE("chats/{id}")
    suspend fun deleteChat(@Path("id") id: Int): Response<Unit>
    
    // Сообщения
    @GET("messages")
    suspend fun getMessages(
        @Query("chat_id") chatId: Int,
        @Query("limit") limit: Int = 20,
        @Query("offset") offset: Int = 0
    ): Response<List<Message>>
    
    @POST("messages")
    suspend fun sendMessage(@Body request: SendMessageRequest): Response<Message>
    
    @GET("messages/{id}")
    suspend fun getMessage(@Path("id") id: Int): Response<Message>
    
    @PUT("messages/{id}")
    suspend fun updateMessage(@Path("id") id: Int, @Body request: UpdateMessageRequest): Response<Message>
    
    @DELETE("messages/{id}")
    suspend fun deleteMessage(@Path("id") id: Int): Response<Unit>
    
    @POST("messages/{id}/read")
    suspend fun markMessageAsRead(@Path("id") id: Int): Response<Unit>
    
    // Контакты
    @GET("contacts")
    suspend fun getContacts(): Response<List<User>>
    
    @POST("contacts")
    suspend fun addContact(@Body request: AddContactRequest): Response<User>
    
    @DELETE("contacts/{id}")
    suspend fun removeContact(@Path("id") id: Int): Response<Unit>
    
    // Загрузка файлов
    @Multipart
    @POST("upload")
    suspend fun uploadFile(
        @Part("chat_id") chatId: Int,
        @Part file: okhttp3.MultipartBody.Part
    ): Response<UploadResponse>
}

// Data classes для запросов
data class RegisterRequest(
    val username: String,
    val email: String,
    val password: String,
    val firstName: String? = null,
    val lastName: String? = null
)

data class LoginRequest(
    val username: String,
    val password: String
)

data class AuthResponse(
    val user: User,
    val token: String
)

data class UpdateProfileRequest(
    val firstName: String? = null,
    val lastName: String? = null,
    val avatar: String? = null
)

data class UpdateStatusRequest(
    val status: String
)

data class CreateChatRequest(
    val name: String,
    val type: String,
    val participants: List<String>
)

data class UpdateChatRequest(
    val name: String? = null,
    val description: String? = null
)

data class SendMessageRequest(
    val chatId: Int,
    val content: String,
    val type: String = "text",
    val replyToId: Int? = null
)

data class UpdateMessageRequest(
    val content: String
)

data class AddContactRequest(
    val username: String,
    val nickname: String? = null
)

data class UploadResponse(
    val fileUrl: String,
    val fileName: String,
    val fileSize: Long,
    val mimeType: String
)