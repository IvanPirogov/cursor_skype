// Система уведомлений
class NotificationSystem {
    constructor() {
        this.container = document.getElementById('notifications');
        this.notifications = new Map();
        this.counter = 0;
    }

    show(type, title, message, duration = 5000) {
        const id = ++this.counter;
        const notification = this.createNotification(id, type, title, message);
        
        this.container.appendChild(notification);
        this.notifications.set(id, notification);

        // Показать уведомление с анимацией
        setTimeout(() => {
            notification.classList.add('show');
        }, 100);

        // Автоматически скрыть уведомление
        if (duration > 0) {
            setTimeout(() => {
                this.hide(id);
            }, duration);
        }

        return id;
    }

    createNotification(id, type, title, message) {
        const notification = document.createElement('div');
        notification.className = `notification ${type}`;
        notification.dataset.id = id;

        const icon = this.getIcon(type);
        
        notification.innerHTML = `
            <div class="notification-icon">
                <i class="${icon}"></i>
            </div>
            <div class="notification-content">
                <div class="notification-title">${title}</div>
                <div class="notification-message">${message}</div>
            </div>
            <button class="notification-close" onclick="notifications.hide(${id})">
                <i class="fas fa-times"></i>
            </button>
        `;

        return notification;
    }

    getIcon(type) {
        const icons = {
            success: 'fas fa-check-circle',
            error: 'fas fa-exclamation-circle',
            warning: 'fas fa-exclamation-triangle',
            info: 'fas fa-info-circle'
        };
        return icons[type] || icons.info;
    }

    hide(id) {
        const notification = this.notifications.get(id);
        if (notification) {
            notification.classList.remove('show');
            setTimeout(() => {
                if (notification.parentNode) {
                    notification.parentNode.removeChild(notification);
                }
                this.notifications.delete(id);
            }, 300);
        }
    }

    clear() {
        this.notifications.forEach((notification, id) => {
            this.hide(id);
        });
    }

    // Вспомогательные методы
    success(title, message, duration) {
        return this.show('success', title, message, duration);
    }

    error(title, message, duration) {
        return this.show('error', title, message, duration);
    }

    warning(title, message, duration) {
        return this.show('warning', title, message, duration);
    }

    info(title, message, duration) {
        return this.show('info', title, message, duration);
    }
}

// Глобальная инициализация
const notifications = new NotificationSystem();

// Экспорт для использования в других модулях
window.notifications = notifications;