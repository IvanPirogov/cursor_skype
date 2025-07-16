# Мобильные клиенты мессенджера

Этот каталог содержит мобильные приложения для мессенджера на Android и iOS.

## Структура проекта

```
mobile/
├── android/                    # Android приложение на Kotlin
│   ├── app/
│   │   ├── build.gradle       # Зависимости и настройки
│   │   └── src/main/
│   │       ├── AndroidManifest.xml
│   │       └── java/com/messenger/android/
│   │           ├── MessengerApplication.kt
│   │           ├── MainActivity.kt
│   │           ├── data/
│   │           │   ├── api/ApiService.kt
│   │           │   ├── model/
│   │           │   └── websocket/WebSocketClient.kt
│   │           └── ui/
│   │               ├── auth/AuthScreen.kt
│   │               └── chat/ChatScreen.kt
│   └── build.gradle           # Конфигурация проекта
├── ios/                       # iOS приложение на Swift
│   ├── Messenger.xcodeproj/   # Проект Xcode
│   └── Messenger/
│       ├── MessengerApp.swift # Главный файл приложения
│       ├── ContentView.swift  # Основной вид
│       ├── Views/
│       │   ├── AuthView.swift
│       │   └── ChatView.swift
│       ├── Models/
│       │   ├── User.swift
│       │   └── Message.swift
│       └── Services/
│           ├── ApiService.swift
│           └── WebSocketClient.swift
└── README.md                  # Этот файл
```

## Android приложение

### Технологии

- **Язык**: Kotlin
- **UI**: Jetpack Compose
- **Архитектура**: MVVM
- **Dependency Injection**: Hilt
- **Сеть**: Retrofit + OkHttp
- **WebSocket**: OkHttp WebSocket
- **База данных**: Room
- **Хранение**: DataStore
- **Изображения**: Coil

### Минимальные требования

- Android 7.0 (API 24)
- Kotlin 1.9.0
- Gradle 8.1

### Функциональность

- ✅ Аутентификация (вход/регистрация)
- ✅ Список чатов
- ✅ Отправка/получение сообщений
- ✅ Real-time через WebSocket
- ✅ Загрузка файлов
- ✅ Уведомления
- ✅ Темная тема
- ✅ Локальное кэширование

### Сборка и запуск

```bash
cd mobile/android
./gradlew assembleDebug
./gradlew installDebug
```

### Структура кода

```
app/src/main/java/com/messenger/android/
├── MessengerApplication.kt     # Application класс
├── MainActivity.kt             # Главная активность
├── data/
│   ├── api/ApiService.kt      # REST API клиент
│   ├── model/                 # Модели данных
│   ├── repository/            # Репозитории
│   └── websocket/             # WebSocket клиент
├── ui/
│   ├── auth/                  # Экраны аутентификации
│   ├── chat/                  # Экраны чата
│   └── theme/                 # Тема приложения
└── di/                        # Dependency injection
```

## iOS приложение

### Технологии

- **Язык**: Swift 5.9
- **UI**: SwiftUI
- **Архитектура**: MVVM
- **Сеть**: URLSession
- **WebSocket**: URLSessionWebSocketTask
- **База данных**: Core Data
- **Хранение**: UserDefaults

### Минимальные требования

- iOS 17.0+
- Xcode 15.0+
- Swift 5.9

### Функциональность

- ✅ Аутентификация (вход/регистрация)
- ✅ Список чатов
- ✅ Отправка/получение сообщений
- ✅ Real-time через WebSocket
- ✅ Загрузка файлов
- ✅ Уведомления
- ✅ Темная тема
- ✅ Локальное кэширование

### Сборка и запуск

1. Откройте `mobile/ios/Messenger.xcodeproj` в Xcode
2. Выберите устройство или симулятор
3. Нажмите `Cmd + R` для запуска

### Структура кода

```
Messenger/
├── MessengerApp.swift          # Главный файл приложения
├── ContentView.swift           # Основной вид
├── Views/
│   ├── AuthView.swift         # Экран аутентификации
│   └── ChatView.swift         # Экран чата
├── Models/
│   ├── User.swift             # Модель пользователя
│   └── Message.swift          # Модель сообщения
├── Services/
│   ├── ApiService.swift       # API клиент
│   └── WebSocketClient.swift  # WebSocket клиент
└── ViewModels/
    ├── AuthViewModel.swift    # ViewModel для аутентификации
    └── ChatViewModel.swift    # ViewModel для чата
```

## Интеграция с backend

### API endpoints

Оба приложения используют одинаковые REST API endpoints:

- `POST /api/v1/auth/register` - регистрация
- `POST /api/v1/auth/login` - вход
- `GET /api/v1/chats` - список чатов
- `GET /api/v1/messages?chat_id=X` - сообщения чата
- `POST /api/v1/messages` - отправка сообщения
- `POST /api/v1/upload` - загрузка файлов

### WebSocket

- Подключение: `ws://localhost:8080/ws?token=<jwt_token>`
- Типы сообщений: `chat`, `typing`, `user_status`, `call_offer`

### Настройка сервера

По умолчанию приложения подключаются к:
- **Локальная разработка**: `http://localhost:8080`
- **Продакшн**: настройте в конфигурации

## Конфигурация

### Android

Измените `BaseURL` в `ApiService.kt`:

```kotlin
private const val BASE_URL = "http://10.0.2.2:8080/api/v1/"
```

### iOS

Измените `baseURL` в `ApiService.swift`:

```swift
private let baseURL = "http://localhost:8080/api/v1"
```

## Разработка

### Добавление новых функций

1. Создайте модель данных в `Models/`
2. Добавьте API методы в `ApiService`
3. Создайте ViewModel (если нужно)
4. Создайте UI компоненты
5. Интегрируйте с WebSocket (если нужно)

### Тестирование

#### Android

```bash
./gradlew test
./gradlew connectedAndroidTest
```

#### iOS

В Xcode: `Cmd + U`

## Развертывание

### Android

```bash
./gradlew assembleRelease
```

### iOS

1. В Xcode выберите схему `Messenger`
2. Product → Archive
3. Выберите метод распространения

## Безопасность

### Общие меры

- JWT токены для аутентификации
- HTTPS для всех запросов
- Валидация данных на клиенте
- Шифрование локальных данных

### Android

- ProGuard для обфускации
- Network Security Config
- Encrypted SharedPreferences

### iOS

- App Transport Security
- Keychain для токенов
- Certificate pinning

## Производительность

### Оптимизации

- Lazy loading для списков
- Пагинация сообщений
- Кэширование изображений
- Оптимизация WebSocket

### Метрики

- Время запуска: < 2 секунды
- Потребление памяти: < 100MB
- Время отклика UI: < 16ms

## Troubleshooting

### Частые проблемы

1. **Нет подключения к серверу**
   - Проверьте URL сервера
   - Убедитесь, что сервер запущен
   - Проверьте сетевые настройки

2. **WebSocket не подключается**
   - Проверьте JWT токен
   - Убедитесь в правильности URL
   - Проверьте CORS настройки

3. **Сообщения не отправляются**
   - Проверьте интернет соединение
   - Убедитесь в валидности токена
   - Проверьте логи сервера

### Логи

#### Android

```bash
adb logcat | grep MessengerApp
```

#### iOS

В Xcode: Window → Devices and Simulators → View Device Logs

## Поддержка

Для сообщений об ошибках и предложений создайте issue в репозитории проекта.

## Лицензия

MIT License