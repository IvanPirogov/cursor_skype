import SwiftUI

struct ChatView: View {
    @EnvironmentObject var authManager: AuthManager
    @EnvironmentObject var chatManager: ChatManager
    @EnvironmentObject var webSocketClient: WebSocketClient
    @State private var messageText = ""
    @State private var showingNewChatSheet = false
    
    var body: some View {
        NavigationView {
            HStack(spacing: 0) {
                // Боковая панель с чатами
                ChatListView()
                    .frame(width: 300)
                    .background(Color(.systemGray6))
                
                // Основная область чата
                ChatDetailView()
            }
            .navigationTitle("Messenger")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button(action: {
                        showingNewChatSheet = true
                    }) {
                        Image(systemName: "plus")
                    }
                }
                
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Выйти") {
                        authManager.logout()
                    }
                }
            }
            .sheet(isPresented: $showingNewChatSheet) {
                NewChatView()
            }
        }
        .navigationViewStyle(DoubleColumnNavigationViewStyle())
    }
}

struct ChatListView: View {
    @EnvironmentObject var chatManager: ChatManager
    @EnvironmentObject var authManager: AuthManager
    
    var body: some View {
        VStack {
            // Информация о пользователе
            HStack {
                Circle()
                    .fill(Color.blue)
                    .frame(width: 40, height: 40)
                    .overlay(
                        Text(authManager.currentUser?.initials ?? "")
                            .foregroundColor(.white)
                            .font(.system(size: 16, weight: .semibold))
                    )
                
                VStack(alignment: .leading) {
                    Text(authManager.currentUser?.displayName ?? "Пользователь")
                        .font(.headline)
                    Text("В сети")
                        .font(.caption)
                        .foregroundColor(.green)
                }
                
                Spacer()
            }
            .padding()
            
            // Поиск
            HStack {
                Image(systemName: "magnifyingglass")
                    .foregroundColor(.gray)
                TextField("Поиск чатов", text: .constant(""))
                    .textFieldStyle(PlainTextFieldStyle())
            }
            .padding(.horizontal)
            .padding(.vertical, 8)
            .background(Color(.systemGray5))
            .cornerRadius(10)
            .padding(.horizontal)
            
            // Список чатов
            ScrollView {
                LazyVStack(spacing: 0) {
                    ForEach(chatManager.chats) { chat in
                        ChatRowView(chat: chat)
                            .onTapGesture {
                                chatManager.selectChat(chat)
                            }
                    }
                }
            }
            .refreshable {
                await chatManager.loadChats()
            }
        }
    }
}

struct ChatRowView: View {
    let chat: Chat
    @EnvironmentObject var chatManager: ChatManager
    
    var body: some View {
        HStack {
            // Аватар чата
            Circle()
                .fill(LinearGradient(
                    gradient: Gradient(colors: [Color.blue, Color.purple]),
                    startPoint: .topLeading,
                    endPoint: .bottomTrailing
                ))
                .frame(width: 50, height: 50)
                .overlay(
                    Text(chat.initials)
                        .foregroundColor(.white)
                        .font(.system(size: 18, weight: .semibold))
                )
            
            VStack(alignment: .leading, spacing: 4) {
                HStack {
                    Text(chat.name)
                        .font(.headline)
                        .foregroundColor(.primary)
                    
                    Spacer()
                    
                    Text(chat.lastMessageTime)
                        .font(.caption)
                        .foregroundColor(.secondary)
                }
                
                HStack {
                    Text(chat.lastMessage ?? "Нет сообщений")
                        .font(.subheadline)
                        .foregroundColor(.secondary)
                        .lineLimit(1)
                    
                    Spacer()
                    
                    if chat.unreadCount > 0 {
                        Text("\(chat.unreadCount)")
                            .font(.caption)
                            .foregroundColor(.white)
                            .padding(.horizontal, 8)
                            .padding(.vertical, 2)
                            .background(Color.blue)
                            .cornerRadius(10)
                    }
                }
            }
            
            Spacer()
        }
        .padding()
        .background(chatManager.currentChat?.id == chat.id ? Color.blue.opacity(0.1) : Color.clear)
        .overlay(
            Rectangle()
                .fill(Color.gray.opacity(0.2))
                .frame(height: 0.5),
            alignment: .bottom
        )
    }
}

struct ChatDetailView: View {
    @EnvironmentObject var chatManager: ChatManager
    @EnvironmentObject var webSocketClient: WebSocketClient
    @State private var messageText = ""
    
    var body: some View {
        if let chat = chatManager.currentChat {
            VStack {
                // Заголовок чата
                HStack {
                    Circle()
                        .fill(LinearGradient(
                            gradient: Gradient(colors: [Color.blue, Color.purple]),
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        ))
                        .frame(width: 40, height: 40)
                        .overlay(
                            Text(chat.initials)
                                .foregroundColor(.white)
                                .font(.system(size: 16, weight: .semibold))
                        )
                    
                    VStack(alignment: .leading) {
                        Text(chat.name)
                            .font(.headline)
                        Text("В сети")
                            .font(.caption)
                            .foregroundColor(.green)
                    }
                    
                    Spacer()
                    
                    Button(action: {}) {
                        Image(systemName: "phone")
                    }
                    
                    Button(action: {}) {
                        Image(systemName: "video")
                    }
                }
                .padding()
                .background(Color(.systemGray6))
                
                // Сообщения
                ScrollView {
                    LazyVStack(spacing: 8) {
                        ForEach(chatManager.messages) { message in
                            MessageRowView(message: message)
                        }
                    }
                    .padding()
                }
                .background(Color(.systemBackground))
                
                // Поле ввода сообщения
                HStack {
                    Button(action: {}) {
                        Image(systemName: "paperclip")
                            .foregroundColor(.gray)
                    }
                    
                    TextField("Введите сообщение...", text: $messageText)
                        .textFieldStyle(RoundedBorderTextFieldStyle())
                    
                    Button(action: {
                        sendMessage()
                    }) {
                        Image(systemName: "arrow.up.circle.fill")
                            .font(.title2)
                            .foregroundColor(messageText.isEmpty ? .gray : .blue)
                    }
                    .disabled(messageText.isEmpty)
                }
                .padding()
                .background(Color(.systemGray6))
            }
        } else {
            VStack {
                Image(systemName: "message")
                    .font(.system(size: 60))
                    .foregroundColor(.gray)
                
                Text("Выберите чат")
                    .font(.title2)
                    .foregroundColor(.gray)
                
                Text("Выберите чат из списка слева, чтобы начать общение")
                    .font(.caption)
                    .foregroundColor(.secondary)
                    .multilineTextAlignment(.center)
            }
            .padding()
        }
    }
    
    private func sendMessage() {
        guard !messageText.isEmpty, let chat = chatManager.currentChat else { return }
        
        let content = messageText
        messageText = ""
        
        // Отправляем через WebSocket
        webSocketClient.sendMessage(content: content, chatId: chat.id)
        
        // Отправляем через API
        Task {
            await chatManager.sendMessage(content: content, chatId: chat.id)
        }
    }
}

struct MessageRowView: View {
    let message: Message
    
    var body: some View {
        HStack {
            if message.isOwn {
                Spacer()
                
                VStack(alignment: .trailing) {
                    Text(message.content)
                        .padding()
                        .background(Color.blue)
                        .foregroundColor(.white)
                        .cornerRadius(16)
                    
                    Text(message.formattedTime)
                        .font(.caption)
                        .foregroundColor(.secondary)
                }
                .frame(maxWidth: .infinity * 0.7, alignment: .trailing)
            } else {
                VStack(alignment: .leading) {
                    Text(message.content)
                        .padding()
                        .background(Color(.systemGray5))
                        .foregroundColor(.primary)
                        .cornerRadius(16)
                    
                    Text(message.formattedTime)
                        .font(.caption)
                        .foregroundColor(.secondary)
                }
                .frame(maxWidth: .infinity * 0.7, alignment: .leading)
                
                Spacer()
            }
        }
    }
}

struct NewChatView: View {
    @Environment(\.presentationMode) var presentationMode
    @State private var chatName = ""
    @State private var chatType = ChatType.private
    
    enum ChatType: String, CaseIterable {
        case `private` = "private"
        case group = "group"
        
        var displayName: String {
            switch self {
            case .private: return "Приватный чат"
            case .group: return "Групповой чат"
            }
        }
    }
    
    var body: some View {
        NavigationView {
            Form {
                Section {
                    TextField("Название чата", text: $chatName)
                }
                
                Section {
                    Picker("Тип чата", selection: $chatType) {
                        ForEach(ChatType.allCases, id: \.self) { type in
                            Text(type.displayName).tag(type)
                        }
                    }
                    .pickerStyle(SegmentedPickerStyle())
                }
            }
            .navigationTitle("Новый чат")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Отмена") {
                        presentationMode.wrappedValue.dismiss()
                    }
                }
                
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Создать") {
                        createChat()
                    }
                    .disabled(chatName.isEmpty)
                }
            }
        }
    }
    
    private func createChat() {
        // TODO: Реализовать создание чата
        presentationMode.wrappedValue.dismiss()
    }
}

#Preview {
    ChatView()
        .environmentObject(AuthManager())
        .environmentObject(ChatManager())
        .environmentObject(WebSocketClient())
}