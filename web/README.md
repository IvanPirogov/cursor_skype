# Веб-клиент мессенджера

Современный веб-интерфейс для мессенджера, построенный на чистом JavaScript.

## Быстрый старт

1. **Запустите backend сервер:**
   ```bash
   cd ..
   make build
   ./bin/server
   ```

2. **Откройте браузер:**
   ```
   http://localhost:8080
   ```

3. **Зарегистрируйтесь или войдите в систему**

## Файлы

### HTML страницы
- `index.html` - Главная страница (аутентификация)
- `chat.html` - Страница чата

### Стили (CSS)
- `static/css/common.css` - Общие стили
- `static/css/auth.css` - Стили аутентификации
- `static/css/chat.css` - Стили чата

### JavaScript
- `static/js/notifications.js` - Система уведомлений
- `static/js/api.js` - API клиент
- `static/js/auth.js` - Аутентификация
- `static/js/websocket.js` - WebSocket клиент
- `static/js/chat.js` - Основной контроллер чата

### Изображения
- `static/images/favicon.svg` - Иконка сайта

## Функции

- ✅ Регистрация и вход
- ✅ Real-time чат
- ✅ Загрузка файлов
- ✅ Контакты
- ✅ Мобильная адаптация
- ✅ Уведомления

## Технологии

- **HTML5** - разметка страниц
- **CSS3** - стили и анимации
- **JavaScript ES6+** - логика приложения
- **WebSocket** - real-time связь
- **Font Awesome** - иконки
- **Google Fonts** - шрифт Inter

## Браузеры

- Chrome 60+
- Firefox 55+
- Safari 12+
- Edge 79+

## Разработка

Для разработки просто редактируйте файлы в папке `static/` и обновляйте страницу в браузере.

## Подробная документация

См. [docs/WEB_CLIENT.md](../docs/WEB_CLIENT.md) для полной документации.