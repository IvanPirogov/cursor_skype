import SwiftUI

struct AuthView: View {
    @EnvironmentObject var authManager: AuthManager
    @State private var isLogin = true
    @State private var username = ""
    @State private var email = ""
    @State private var password = ""
    @State private var firstName = ""
    @State private var lastName = ""
    @State private var showingAlert = false
    
    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 20) {
                    // Заголовок
                    VStack(spacing: 8) {
                        Text("Messenger")
                            .font(.largeTitle)
                            .fontWeight(.bold)
                            .foregroundColor(.blue)
                        
                        Text("Современный мессенджер для общения")
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                            .multilineTextAlignment(.center)
                    }
                    .padding(.top, 40)
                    .padding(.bottom, 20)
                    
                    // Переключатель режима
                    Picker("Режим", selection: $isLogin) {
                        Text("Вход").tag(true)
                        Text("Регистрация").tag(false)
                    }
                    .pickerStyle(SegmentedPickerStyle())
                    .padding(.horizontal)
                    
                    // Форма
                    VStack(spacing: 16) {
                        // Имя пользователя
                        TextField("Имя пользователя", text: $username)
                            .textFieldStyle(RoundedBorderTextFieldStyle())
                            .autocapitalization(.none)
                            .disableAutocorrection(true)
                        
                        // Email (только для регистрации)
                        if !isLogin {
                            TextField("Email", text: $email)
                                .textFieldStyle(RoundedBorderTextFieldStyle())
                                .autocapitalization(.none)
                                .disableAutocorrection(true)
                                .keyboardType(.emailAddress)
                        }
                        
                        // Пароль
                        SecureField("Пароль", text: $password)
                            .textFieldStyle(RoundedBorderTextFieldStyle())
                        
                        // Дополнительные поля для регистрации
                        if !isLogin {
                            TextField("Имя", text: $firstName)
                                .textFieldStyle(RoundedBorderTextFieldStyle())
                            
                            TextField("Фамилия", text: $lastName)
                                .textFieldStyle(RoundedBorderTextFieldStyle())
                        }
                    }
                    .padding(.horizontal)
                    
                    // Кнопка отправки
                    Button(action: {
                        Task {
                            if isLogin {
                                await authManager.login(username: username, password: password)
                            } else {
                                await authManager.register(
                                    username: username,
                                    email: email,
                                    password: password,
                                    firstName: firstName.isEmpty ? nil : firstName,
                                    lastName: lastName.isEmpty ? nil : lastName
                                )
                            }
                            
                            if authManager.errorMessage != nil {
                                showingAlert = true
                            }
                        }
                    }) {
                        HStack {
                            if authManager.isLoading {
                                ProgressView()
                                    .progressViewStyle(CircularProgressViewStyle(tint: .white))
                                    .scaleEffect(0.8)
                            }
                            
                            Text(isLogin ? "Войти" : "Зарегистрироваться")
                                .fontWeight(.semibold)
                        }
                        .frame(maxWidth: .infinity)
                        .padding()
                        .background(Color.blue)
                        .foregroundColor(.white)
                        .cornerRadius(12)
                    }
                    .disabled(authManager.isLoading || username.isEmpty || password.isEmpty)
                    .padding(.horizontal)
                    
                    Spacer()
                }
            }
            .navigationTitle("")
            .navigationBarHidden(true)
            .alert("Ошибка", isPresented: $showingAlert) {
                Button("OK") {
                    authManager.errorMessage = nil
                }
            } message: {
                Text(authManager.errorMessage ?? "Неизвестная ошибка")
            }
        }
    }
}

#Preview {
    AuthView()
        .environmentObject(AuthManager())
}