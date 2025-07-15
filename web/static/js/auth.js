// Управление аутентификацией
class AuthManager {
    constructor() {
        this.loginForm = document.getElementById('loginForm');
        this.registerForm = document.getElementById('registerForm');
        this.loginFormDiv = document.getElementById('login-form');
        this.registerFormDiv = document.getElementById('register-form');
        
        this.initializeEventListeners();
        this.checkAuthStatus();
    }

    initializeEventListeners() {
        // Переключение между формами
        document.getElementById('show-register').addEventListener('click', (e) => {
            e.preventDefault();
            this.showRegisterForm();
        });

        document.getElementById('show-login').addEventListener('click', (e) => {
            e.preventDefault();
            this.showLoginForm();
        });

        // Обработка отправки форм
        this.loginForm.addEventListener('submit', (e) => {
            e.preventDefault();
            this.handleLogin();
        });

        this.registerForm.addEventListener('submit', (e) => {
            e.preventDefault();
            this.handleRegister();
        });

        // Валидация в реальном времени
        this.setupFormValidation();
    }

    setupFormValidation() {
        const inputs = document.querySelectorAll('.auth-form input');
        inputs.forEach(input => {
            input.addEventListener('blur', () => {
                this.validateField(input);
            });

            input.addEventListener('input', () => {
                this.clearFieldError(input);
            });
        });
    }

    validateField(field) {
        const value = field.value.trim();
        const type = field.type;
        const name = field.name;
        let isValid = true;
        let errorMessage = '';

        // Проверка обязательных полей
        if (field.required && !value) {
            isValid = false;
            errorMessage = 'Это поле обязательно для заполнения';
        }

        // Специфичные проверки
        if (value && type === 'email') {
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(value)) {
                isValid = false;
                errorMessage = 'Введите корректный email адрес';
            }
        }

        if (value && name === 'username') {
            if (value.length < 3) {
                isValid = false;
                errorMessage = 'Имя пользователя должно содержать минимум 3 символа';
            }
        }

        if (value && name === 'password') {
            if (value.length < 6) {
                isValid = false;
                errorMessage = 'Пароль должен содержать минимум 6 символов';
            }
        }

        this.setFieldValidation(field, isValid, errorMessage);
        return isValid;
    }

    setFieldValidation(field, isValid, errorMessage) {
        const formGroup = field.closest('.form-group');
        const errorElement = formGroup.querySelector('.form-error') || this.createErrorElement(formGroup);

        if (isValid) {
            formGroup.classList.remove('has-error');
            formGroup.classList.add('has-success');
            errorElement.textContent = '';
        } else {
            formGroup.classList.remove('has-success');
            formGroup.classList.add('has-error');
            errorElement.textContent = errorMessage;
        }
    }

    clearFieldError(field) {
        const formGroup = field.closest('.form-group');
        formGroup.classList.remove('has-error', 'has-success');
    }

    createErrorElement(formGroup) {
        const errorElement = document.createElement('div');
        errorElement.className = 'form-error';
        formGroup.appendChild(errorElement);
        return errorElement;
    }

    showLoginForm() {
        this.registerFormDiv.classList.remove('active');
        this.loginFormDiv.classList.add('active');
        
        // Очистка форм
        this.clearFormErrors();
        this.registerForm.reset();
    }

    showRegisterForm() {
        this.loginFormDiv.classList.remove('active');
        this.registerFormDiv.classList.add('active');
        
        // Очистка форм
        this.clearFormErrors();
        this.loginForm.reset();
    }

    clearFormErrors() {
        const formGroups = document.querySelectorAll('.form-group');
        formGroups.forEach(group => {
            group.classList.remove('has-error', 'has-success');
            const errorElement = group.querySelector('.form-error');
            if (errorElement) {
                errorElement.textContent = '';
            }
        });
    }

    async handleLogin() {
        const formData = new FormData(this.loginForm);
        const credentials = {
            username: formData.get('username'),
            password: formData.get('password')
        };

        // Валидация формы
        if (!this.validateForm(this.loginForm)) {
            return;
        }

        const submitBtn = this.loginForm.querySelector('button[type="submit"]');
        this.setButtonLoading(submitBtn, true);

        try {
            const response = await api.login(credentials);
            
            // Сохранение токена
            api.setToken(response.token);
            
            // Сохранение информации о пользователе
            localStorage.setItem('user_info', JSON.stringify(response.user));
            
            notifications.success('Успешно!', 'Вы успешно вошли в систему');
            
            // Перенаправление на страницу чата
            setTimeout(() => {
                window.location.href = '/chat.html';
            }, 1000);
            
        } catch (error) {
            notifications.error('Ошибка входа', error.message);
        } finally {
            this.setButtonLoading(submitBtn, false);
        }
    }

    async handleRegister() {
        const formData = new FormData(this.registerForm);
        const userData = {
            username: formData.get('username'),
            email: formData.get('email'),
            password: formData.get('password'),
            first_name: formData.get('first_name'),
            last_name: formData.get('last_name')
        };

        // Валидация формы
        if (!this.validateForm(this.registerForm)) {
            return;
        }

        const submitBtn = this.registerForm.querySelector('button[type="submit"]');
        this.setButtonLoading(submitBtn, true);

        try {
            const response = await api.register(userData);
            
            // Сохранение токена
            api.setToken(response.token);
            
            // Сохранение информации о пользователе
            localStorage.setItem('user_info', JSON.stringify(response.user));
            
            notifications.success('Регистрация успешна!', 'Добро пожаловать в мессенджер');
            
            // Перенаправление на страницу чата
            setTimeout(() => {
                window.location.href = '/chat.html';
            }, 1000);
            
        } catch (error) {
            notifications.error('Ошибка регистрации', error.message);
        } finally {
            this.setButtonLoading(submitBtn, false);
        }
    }

    validateForm(form) {
        const inputs = form.querySelectorAll('input[required]');
        let isValid = true;

        inputs.forEach(input => {
            if (!this.validateField(input)) {
                isValid = false;
            }
        });

        return isValid;
    }

    setButtonLoading(button, isLoading) {
        if (isLoading) {
            button.classList.add('loading');
            button.disabled = true;
        } else {
            button.classList.remove('loading');
            button.disabled = false;
        }
    }

    checkAuthStatus() {
        const token = localStorage.getItem('auth_token');
        if (token) {
            // Если токен есть, проверим его валидность
            this.verifyToken(token);
        }
    }

    async verifyToken(token) {
        try {
            api.setToken(token);
            const user = await api.getCurrentUser();
            
            // Токен валиден, перенаправляем на чат
            window.location.href = '/chat.html';
        } catch (error) {
            // Токен невалиден, очищаем localStorage
            localStorage.removeItem('auth_token');
            localStorage.removeItem('user_info');
            api.setToken(null);
        }
    }

    static logout() {
        localStorage.removeItem('auth_token');
        localStorage.removeItem('user_info');
        api.setToken(null);
        window.location.href = '/';
    }
}

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
    new AuthManager();
});

// Экспорт для использования в других модулях
window.AuthManager = AuthManager;