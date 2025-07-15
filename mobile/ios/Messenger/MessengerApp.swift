import SwiftUI

@main
struct MessengerApp: App {
    @StateObject private var authManager = AuthManager()
    @StateObject private var chatManager = ChatManager()
    @StateObject private var webSocketClient = WebSocketClient()
    
    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(authManager)
                .environmentObject(chatManager)
                .environmentObject(webSocketClient)
        }
    }
}

// MARK: - AuthManager
class AuthManager: ObservableObject {
    @Published var isAuthenticated = false
    @Published var currentUser: User?
    @Published var isLoading = false
    @Published var errorMessage: String?
    
    private let apiService = ApiService()
    
    init() {
        checkAuthStatus()
    }
    
    func checkAuthStatus() {
        if let token = UserDefaults.standard.string(forKey: "auth_token") {
            // Проверяем валидность токена
            Task {
                await getCurrentUser()
            }
        }
    }
    
    @MainActor
    func login(username: String, password: String) async {
        isLoading = true
        errorMessage = nil
        
        do {
            let response = try await apiService.login(username: username, password: password)
            
            // Сохраняем токен и пользователя
            UserDefaults.standard.set(response.token, forKey: "auth_token")
            currentUser = response.user
            isAuthenticated = true
            
        } catch {
            errorMessage = error.localizedDescription
        }
        
        isLoading = false
    }
    
    @MainActor
    func register(username: String, email: String, password: String, firstName: String?, lastName: String?) async {
        isLoading = true
        errorMessage = nil
        
        do {
            let response = try await apiService.register(
                username: username,
                email: email,
                password: password,
                firstName: firstName,
                lastName: lastName
            )
            
            // Сохраняем токен и пользователя
            UserDefaults.standard.set(response.token, forKey: "auth_token")
            currentUser = response.user
            isAuthenticated = true
            
        } catch {
            errorMessage = error.localizedDescription
        }
        
        isLoading = false
    }
    
    @MainActor
    func logout() {
        UserDefaults.standard.removeObject(forKey: "auth_token")
        currentUser = nil
        isAuthenticated = false
    }
    
    @MainActor
    private func getCurrentUser() async {
        guard let token = UserDefaults.standard.string(forKey: "auth_token") else { return }
        
        do {
            let user = try await apiService.getCurrentUser(token: token)
            currentUser = user
            isAuthenticated = true
        } catch {
            // Токен невалиден
            UserDefaults.standard.removeObject(forKey: "auth_token")
            isAuthenticated = false
        }
    }
}

// MARK: - ChatManager
class ChatManager: ObservableObject {
    @Published var chats: [Chat] = []
    @Published var messages: [Message] = []
    @Published var currentChat: Chat?
    @Published var isLoading = false
    @Published var errorMessage: String?
    
    private let apiService = ApiService()
    
    @MainActor
    func loadChats() async {
        guard let token = UserDefaults.standard.string(forKey: "auth_token") else { return }
        
        isLoading = true
        
        do {
            chats = try await apiService.getChats(token: token)
        } catch {
            errorMessage = error.localizedDescription
        }
        
        isLoading = false
    }
    
    @MainActor
    func loadMessages(for chatId: Int) async {
        guard let token = UserDefaults.standard.string(forKey: "auth_token") else { return }
        
        isLoading = true
        
        do {
            messages = try await apiService.getMessages(chatId: chatId, token: token)
        } catch {
            errorMessage = error.localizedDescription
        }
        
        isLoading = false
    }
    
    @MainActor
    func sendMessage(content: String, chatId: Int) async {
        guard let token = UserDefaults.standard.string(forKey: "auth_token") else { return }
        
        do {
            let message = try await apiService.sendMessage(
                content: content,
                chatId: chatId,
                token: token
            )
            messages.append(message)
        } catch {
            errorMessage = error.localizedDescription
        }
    }
    
    func selectChat(_ chat: Chat) {
        currentChat = chat
        Task {
            await loadMessages(for: chat.id)
        }
    }
}