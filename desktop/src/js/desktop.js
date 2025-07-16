/**
 * Desktop-specific functionality for Electron app
 */

class DesktopApp {
    constructor() {
        this.isElectron = typeof window.electronAPI !== 'undefined';
        this.platform = this.isElectron ? window.electronAPI.app.platform : 'web';
        this.titleBarHeight = 32;
        this.currentZoom = 1;
        this.isDarkTheme = false;
        
        if (this.isElectron) {
            this.initElectronFeatures();
        }
    }

    initElectronFeatures() {
        this.setupTitleBar();
        this.setupTheme();
        this.setupKeyboardShortcuts();
        this.setupElectronEventListeners();
        this.setupWindowControls();
        this.setupContextMenu();
        this.loadSettings();
    }

    setupTitleBar() {
        const titleBar = document.getElementById('titleBar');
        if (!titleBar) return;

        // Hide title bar on macOS (system title bar is used)
        if (this.platform === 'darwin') {
            titleBar.style.display = 'none';
            document.body.classList.add('platform-darwin');
        } else {
            document.body.classList.add(`platform-${this.platform}`);
        }

        // Window controls
        const minimizeBtn = document.getElementById('minimizeBtn');
        const maximizeBtn = document.getElementById('maximizeBtn');
        const closeBtn = document.getElementById('closeBtn');

        if (minimizeBtn) {
            minimizeBtn.addEventListener('click', () => {
                // Electron will handle minimize via main process
                this.minimize();
            });
        }

        if (maximizeBtn) {
            maximizeBtn.addEventListener('click', () => {
                this.toggleMaximize();
            });
        }

        if (closeBtn) {
            closeBtn.addEventListener('click', () => {
                this.close();
            });
        }
    }

    minimize() {
        if (this.isElectron) {
            // This will be handled by main process
            const { remote } = require('electron');
            remote.getCurrentWindow().minimize();
        }
    }

    toggleMaximize() {
        if (this.isElectron) {
            const { remote } = require('electron');
            const win = remote.getCurrentWindow();
            if (win.isMaximized()) {
                win.unmaximize();
            } else {
                win.maximize();
            }
        }
    }

    close() {
        if (this.isElectron) {
            const { remote } = require('electron');
            remote.getCurrentWindow().close();
        }
    }

    setupTheme() {
        if (!this.isElectron) return;

        // Listen for system theme changes
        window.theme.onThemeUpdate(() => {
            this.updateTheme();
        });

        // Initial theme setup
        this.updateTheme();
    }

    updateTheme() {
        if (!this.isElectron) return;

        const shouldUseDark = window.theme.shouldUseDarkColors();
        this.isDarkTheme = shouldUseDark;
        
        document.body.classList.toggle('dark-theme', shouldUseDark);
        
        // Update theme variables
        if (shouldUseDark) {
            document.documentElement.classList.add('dark');
        } else {
            document.documentElement.classList.remove('dark');
        }
    }

    setupKeyboardShortcuts() {
        if (!this.isElectron) return;

        // Listen for shortcuts from main process
        window.electronAPI.on('new-chat', () => {
            this.handleNewChat();
        });

        window.electronAPI.on('show-settings', () => {
            this.showSettings();
        });

        window.electronAPI.on('zoom-in', () => {
            this.zoomIn();
        });

        window.electronAPI.on('zoom-out', () => {
            this.zoomOut();
        });

        window.electronAPI.on('zoom-reset', () => {
            this.zoomReset();
        });

        // Handle keyboard shortcuts in renderer
        document.addEventListener('keydown', (e) => {
            if (e.ctrlKey || e.metaKey) {
                switch (e.key) {
                    case 'n':
                        e.preventDefault();
                        this.handleNewChat();
                        break;
                    case ',':
                        e.preventDefault();
                        this.showSettings();
                        break;
                    case 'r':
                        e.preventDefault();
                        this.reload();
                        break;
                    case 'f':
                        e.preventDefault();
                        this.toggleFullscreen();
                        break;
                    case '=':
                    case '+':
                        e.preventDefault();
                        this.zoomIn();
                        break;
                    case '-':
                        e.preventDefault();
                        this.zoomOut();
                        break;
                    case '0':
                        e.preventDefault();
                        this.zoomReset();
                        break;
                }
            }
        });
    }

    setupElectronEventListeners() {
        if (!this.isElectron) return;

        // Network status
        window.network.onOnline(() => {
            this.handleNetworkChange(true);
        });

        window.network.onOffline(() => {
            this.handleNetworkChange(false);
        });
    }

    setupWindowControls() {
        // Auto-resize textarea
        const messageInput = document.getElementById('messageInput');
        if (messageInput) {
            messageInput.addEventListener('input', () => {
                this.autoResizeTextarea(messageInput);
            });
        }

        // Handle window resize
        window.addEventListener('resize', () => {
            this.handleWindowResize();
        });
    }

    setupContextMenu() {
        if (!this.isElectron) return;

        // Custom context menu for specific elements
        document.addEventListener('contextmenu', (e) => {
            const target = e.target;
            
            if (target.classList.contains('message-bubble')) {
                e.preventDefault();
                this.showMessageContextMenu(e, target);
            } else if (target.classList.contains('chat-item')) {
                e.preventDefault();
                this.showChatContextMenu(e, target);
            }
        });
    }

    showMessageContextMenu(e, messageElement) {
        const contextMenu = document.getElementById('contextMenu');
        if (!contextMenu) return;

        const items = [
            { text: 'Копировать', action: () => this.copyMessage(messageElement) },
            { text: 'Ответить', action: () => this.replyToMessage(messageElement) },
            { text: 'Переслать', action: () => this.forwardMessage(messageElement) },
            { separator: true },
            { text: 'Удалить', action: () => this.deleteMessage(messageElement) }
        ];

        this.showContextMenu(contextMenu, e.pageX, e.pageY, items);
    }

    showChatContextMenu(e, chatElement) {
        const contextMenu = document.getElementById('contextMenu');
        if (!contextMenu) return;

        const items = [
            { text: 'Открыть чат', action: () => this.openChat(chatElement) },
            { text: 'Закрепить', action: () => this.pinChat(chatElement) },
            { separator: true },
            { text: 'Архивировать', action: () => this.archiveChat(chatElement) },
            { text: 'Удалить', action: () => this.deleteChat(chatElement) }
        ];

        this.showContextMenu(contextMenu, e.pageX, e.pageY, items);
    }

    showContextMenu(contextMenu, x, y, items) {
        contextMenu.innerHTML = '';

        items.forEach(item => {
            if (item.separator) {
                const separator = document.createElement('div');
                separator.className = 'context-menu-separator';
                contextMenu.appendChild(separator);
            } else {
                const menuItem = document.createElement('div');
                menuItem.className = 'context-menu-item';
                menuItem.textContent = item.text;
                menuItem.addEventListener('click', () => {
                    item.action();
                    this.hideContextMenu();
                });
                contextMenu.appendChild(menuItem);
            }
        });

        contextMenu.style.left = `${x}px`;
        contextMenu.style.top = `${y}px`;
        contextMenu.classList.remove('hidden');

        // Hide context menu on outside click
        const hideHandler = (e) => {
            if (!contextMenu.contains(e.target)) {
                this.hideContextMenu();
                document.removeEventListener('click', hideHandler);
            }
        };

        setTimeout(() => {
            document.addEventListener('click', hideHandler);
        }, 0);
    }

    hideContextMenu() {
        const contextMenu = document.getElementById('contextMenu');
        if (contextMenu) {
            contextMenu.classList.add('hidden');
        }
    }

    // Notification methods
    async showNotification(title, body, options = {}) {
        if (!this.isElectron) return;

        const notificationOptions = {
            title,
            body,
            icon: options.icon,
            ...options
        };

        try {
            await window.electronAPI.notification.show(notificationOptions);
        } catch (error) {
            console.error('Failed to show notification:', error);
        }
    }

    async updateBadgeCount(count) {
        if (!this.isElectron) return;

        try {
            await window.electronAPI.notification.setBadgeCount(count);
        } catch (error) {
            console.error('Failed to update badge count:', error);
        }
    }

    // Settings methods
    async loadSettings() {
        if (!this.isElectron) return;

        try {
            const settings = {
                theme: await window.electronAPI.store.get('theme'),
                notifications: await window.electronAPI.store.get('notifications'),
                autoLaunch: await window.electronAPI.autoLaunch.get(),
                minimizeToTray: await window.electronAPI.store.get('minimizeToTray')
            };

            this.applySettings(settings);
        } catch (error) {
            console.error('Failed to load settings:', error);
        }
    }

    async saveSettings(settings) {
        if (!this.isElectron) return;

        try {
            for (const [key, value] of Object.entries(settings)) {
                if (key === 'autoLaunch') {
                    await window.electronAPI.autoLaunch.set(value);
                } else {
                    await window.electronAPI.store.set(key, value);
                }
            }
        } catch (error) {
            console.error('Failed to save settings:', error);
        }
    }

    applySettings(settings) {
        // Apply theme
        if (settings.theme) {
            document.documentElement.setAttribute('data-theme', settings.theme);
        }

        // Apply other settings
        if (settings.notifications !== undefined) {
            this.notificationsEnabled = settings.notifications;
        }
    }

    // Utility methods
    autoResizeTextarea(textarea) {
        textarea.style.height = 'auto';
        textarea.style.height = Math.min(textarea.scrollHeight, 120) + 'px';
    }

    handleWindowResize() {
        // Handle any window resize logic here
        this.updateChatLayout();
    }

    updateChatLayout() {
        const sidebar = document.querySelector('.sidebar');
        const chatMain = document.querySelector('.chat-main');
        
        if (window.innerWidth < 768) {
            // Mobile layout adjustments
            if (sidebar) sidebar.style.width = '100%';
        } else {
            // Desktop layout
            if (sidebar) sidebar.style.width = '300px';
        }
    }

    handleNetworkChange(isOnline) {
        const statusIndicator = document.getElementById('networkStatus');
        if (statusIndicator) {
            statusIndicator.textContent = isOnline ? 'В сети' : 'Нет соединения';
            statusIndicator.className = isOnline ? 'online' : 'offline';
        }

        // Notify other components about network change
        if (window.chatApp) {
            window.chatApp.handleNetworkChange(isOnline);
        }
    }

    // Zoom methods
    zoomIn() {
        this.currentZoom = Math.min(this.currentZoom + 0.1, 3);
        this.applyZoom();
    }

    zoomOut() {
        this.currentZoom = Math.max(this.currentZoom - 0.1, 0.5);
        this.applyZoom();
    }

    zoomReset() {
        this.currentZoom = 1;
        this.applyZoom();
    }

    applyZoom() {
        document.body.style.zoom = this.currentZoom;
        
        // Save zoom level
        if (this.isElectron) {
            window.electronAPI.store.set('zoom', this.currentZoom);
        }
    }

    // Action methods
    handleNewChat() {
        // Trigger new chat creation
        const newChatBtn = document.getElementById('newChatBtn');
        if (newChatBtn) {
            newChatBtn.click();
        }
    }

    showSettings() {
        // Show settings screen
        const settingsBtn = document.getElementById('settingsBtn');
        if (settingsBtn) {
            settingsBtn.click();
        }
    }

    reload() {
        if (this.isElectron) {
            location.reload();
        }
    }

    toggleFullscreen() {
        if (this.isElectron) {
            const { remote } = require('electron');
            const win = remote.getCurrentWindow();
            win.setFullScreen(!win.isFullScreen());
        }
    }

    // Message context menu actions
    copyMessage(messageElement) {
        const text = messageElement.textContent;
        if (this.isElectron) {
            window.utils.clipboard.writeText(text);
        } else {
            navigator.clipboard.writeText(text);
        }
    }

    replyToMessage(messageElement) {
        // Handle reply to message
        console.log('Reply to message:', messageElement);
    }

    forwardMessage(messageElement) {
        // Handle forward message
        console.log('Forward message:', messageElement);
    }

    deleteMessage(messageElement) {
        // Handle delete message
        console.log('Delete message:', messageElement);
    }

    // Chat context menu actions
    openChat(chatElement) {
        chatElement.click();
    }

    pinChat(chatElement) {
        // Handle pin chat
        console.log('Pin chat:', chatElement);
    }

    archiveChat(chatElement) {
        // Handle archive chat
        console.log('Archive chat:', chatElement);
    }

    deleteChat(chatElement) {
        // Handle delete chat
        console.log('Delete chat:', chatElement);
    }

    // File handling
    async openFile() {
        if (!this.isElectron) return null;

        try {
            const result = await window.electronAPI.dialog.showOpenDialog({
                properties: ['openFile'],
                filters: [
                    { name: 'Images', extensions: ['jpg', 'jpeg', 'png', 'gif', 'webp'] },
                    { name: 'Videos', extensions: ['mp4', 'avi', 'mov', 'wmv', 'flv'] },
                    { name: 'Documents', extensions: ['pdf', 'doc', 'docx', 'txt', 'rtf'] },
                    { name: 'All Files', extensions: ['*'] }
                ]
            });

            return result.canceled ? null : result.filePaths[0];
        } catch (error) {
            console.error('Failed to open file dialog:', error);
            return null;
        }
    }

    async saveFile(data, filename) {
        if (!this.isElectron) return false;

        try {
            const result = await window.electronAPI.dialog.showSaveDialog({
                defaultPath: filename,
                filters: [
                    { name: 'All Files', extensions: ['*'] }
                ]
            });

            if (result.canceled) return false;

            // Save file logic would go here
            return true;
        } catch (error) {
            console.error('Failed to save file:', error);
            return false;
        }
    }

    // Audio methods
    playNotificationSound() {
        if (this.isElectron) {
            window.messenger.playNotificationSound();
        }
    }

    // Battery info (for laptops)
    async getBatteryInfo() {
        if (this.isElectron) {
            return await window.messenger.getBatteryInfo();
        }
        return null;
    }

    // Vibration (for devices that support it)
    vibrate(pattern = 200) {
        if (this.isElectron) {
            window.messenger.vibrate(pattern);
        }
    }

    // App info
    async getAppInfo() {
        if (!this.isElectron) return null;

        try {
            const version = await window.electronAPI.app.getVersion();
            return {
                version,
                platform: this.platform,
                isElectron: true
            };
        } catch (error) {
            console.error('Failed to get app info:', error);
            return null;
        }
    }

    // Cleanup
    destroy() {
        // Remove event listeners
        if (this.isElectron) {
            window.electronAPI.removeListener('new-chat', this.handleNewChat);
            window.electronAPI.removeListener('show-settings', this.showSettings);
            window.electronAPI.removeListener('zoom-in', this.zoomIn);
            window.electronAPI.removeListener('zoom-out', this.zoomOut);
            window.electronAPI.removeListener('zoom-reset', this.zoomReset);
            
            window.theme.removeThemeUpdateListener(this.updateTheme);
            window.network.removeOnlineListener(this.handleNetworkChange);
            window.network.removeOfflineListener(this.handleNetworkChange);
        }
    }
}

// Initialize desktop app
const desktopApp = new DesktopApp();

// Make desktop app available globally
window.desktopApp = desktopApp;

// Export for modules
if (typeof module !== 'undefined' && module.exports) {
    module.exports = DesktopApp;
}