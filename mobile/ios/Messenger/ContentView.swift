import SwiftUI

struct ContentView: View {
    @EnvironmentObject var authManager: AuthManager
    @EnvironmentObject var chatManager: ChatManager
    @EnvironmentObject var webSocketClient: WebSocketClient
    
    var body: some View {
        Group {
            if authManager.isAuthenticated {
                ChatView()
                    .onAppear {
                        Task {
                            await chatManager.loadChats()
                            if let token = UserDefaults.standard.string(forKey: "auth_token") {
                                await webSocketClient.connect(token: token)
                            }
                        }
                    }
            } else {
                AuthView()
            }
        }
        .animation(.easeInOut, value: authManager.isAuthenticated)
    }
}

#Preview {
    ContentView()
        .environmentObject(AuthManager())
        .environmentObject(ChatManager())
        .environmentObject(WebSocketClient())
}