const { contextBridge, ipcRenderer } = require('electron');

// Безопасный API для renderer процесса
contextBridge.exposeInMainWorld('electronAPI', {
  // Хранилище настроек
  store: {
    get: (key) => ipcRenderer.invoke('get-store-value', key),
    set: (key, value) => ipcRenderer.invoke('set-store-value', key, value)
  },

  // Информация о приложении
  app: {
    getVersion: () => ipcRenderer.invoke('get-app-version'),
    platform: process.platform,
    arch: process.arch
  },

  // Диалоги
  dialog: {
    showMessageBox: (options) => ipcRenderer.invoke('show-message-box', options),
    showOpenDialog: (options) => ipcRenderer.invoke('show-open-dialog', options),
    showSaveDialog: (options) => ipcRenderer.invoke('show-save-dialog', options)
  },

  // Уведомления
  notification: {
    show: (options) => ipcRenderer.invoke('show-notification', options),
    setBadgeCount: (count) => ipcRenderer.invoke('set-badge-count', count)
  },

  // Автозапуск
  autoLaunch: {
    set: (enabled) => ipcRenderer.invoke('set-auto-launch', enabled),
    get: () => ipcRenderer.invoke('get-auto-launch')
  },

  // Обработчики событий от главного процесса
  on: (channel, callback) => {
    const validChannels = [
      'new-chat',
      'show-settings',
      'zoom-in',
      'zoom-out',
      'zoom-reset',
      'theme-changed'
    ];
    
    if (validChannels.includes(channel)) {
      ipcRenderer.on(channel, callback);
    }
  },

  // Удаление обработчиков событий
  removeListener: (channel, callback) => {
    ipcRenderer.removeListener(channel, callback);
  },

  // Отправка событий главному процессу
  send: (channel, data) => {
    const validChannels = [
      'window-ready',
      'user-authenticated',
      'user-logout',
      'new-message',
      'update-unread-count'
    ];
    
    if (validChannels.includes(channel)) {
      ipcRenderer.send(channel, data);
    }
  }
});

// Дополнительные утилиты
contextBridge.exposeInMainWorld('utils', {
  // Утилиты для работы с файлами
  path: {
    join: (...args) => require('path').join(...args),
    dirname: (path) => require('path').dirname(path),
    basename: (path) => require('path').basename(path),
    extname: (path) => require('path').extname(path)
  },

  // Утилиты для работы с URL
  url: {
    isValid: (url) => {
      try {
        new URL(url);
        return true;
      } catch {
        return false;
      }
    }
  },

  // Утилиты для работы с буфером обмена
  clipboard: {
    writeText: (text) => require('electron').clipboard.writeText(text),
    readText: () => require('electron').clipboard.readText(),
    writeImage: (image) => require('electron').clipboard.writeImage(image),
    readImage: () => require('electron').clipboard.readImage()
  }
});

// Безопасный WebSocket API
contextBridge.exposeInMainWorld('WebSocketSecure', {
  create: (url, protocols) => {
    // Создаем WebSocket с дополнительными проверками безопасности
    if (!url.startsWith('ws://') && !url.startsWith('wss://')) {
      throw new Error('Invalid WebSocket URL');
    }
    
    const ws = new WebSocket(url, protocols);
    
    // Прокси для безопасного доступа к WebSocket
    return {
      send: (data) => {
        if (ws.readyState === WebSocket.OPEN) {
          ws.send(data);
        }
      },
      close: (code, reason) => ws.close(code, reason),
      addEventListener: (type, listener) => ws.addEventListener(type, listener),
      removeEventListener: (type, listener) => ws.removeEventListener(type, listener),
      get readyState() { return ws.readyState; },
      get url() { return ws.url; },
      get protocol() { return ws.protocol; }
    };
  }
});

// Константы WebSocket
contextBridge.exposeInMainWorld('WebSocketConstants', {
  CONNECTING: 0,
  OPEN: 1,
  CLOSING: 2,
  CLOSED: 3
});

// Console proxy для отладки
contextBridge.exposeInMainWorld('console', {
  log: (...args) => console.log(...args),
  error: (...args) => console.error(...args),
  warn: (...args) => console.warn(...args),
  info: (...args) => console.info(...args),
  debug: (...args) => console.debug(...args)
});

// Детектор темы системы
contextBridge.exposeInMainWorld('theme', {
  shouldUseDarkColors: () => {
    return require('electron').nativeTheme.shouldUseDarkColors;
  },
  
  onThemeUpdate: (callback) => {
    require('electron').nativeTheme.on('updated', callback);
  },
  
  removeThemeUpdateListener: (callback) => {
    require('electron').nativeTheme.removeListener('updated', callback);
  }
});

// Проверка доступности сети
contextBridge.exposeInMainWorld('network', {
  isOnline: () => navigator.onLine,
  
  onOnline: (callback) => {
    window.addEventListener('online', callback);
  },
  
  onOffline: (callback) => {
    window.addEventListener('offline', callback);
  },
  
  removeOnlineListener: (callback) => {
    window.removeEventListener('online', callback);
  },
  
  removeOfflineListener: (callback) => {
    window.removeEventListener('offline', callback);
  }
});

// Дополнительные возможности для мессенджера
contextBridge.exposeInMainWorld('messenger', {
  // Проверка поддержки уведомлений
  notificationSupported: () => {
    return 'Notification' in window;
  },
  
  // Запрос разрешений на уведомления
  requestNotificationPermission: async () => {
    if ('Notification' in window) {
      return await Notification.requestPermission();
    }
    return 'denied';
  },
  
  // Проверка текущего разрешения на уведомления
  getNotificationPermission: () => {
    if ('Notification' in window) {
      return Notification.permission;
    }
    return 'denied';
  },
  
  // Звуковые уведомления
  playNotificationSound: () => {
    const audio = new Audio('data:audio/wav;base64,UklGRnoGAABXQVZFZm10IBAAAAABAAEAQB8AAEAfAAABAAgAZGF0YQoGAACBhYqFbF1fdJivrJBhNjVgodDbq2EcBj+a2/LDciUFLIHO8tiJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvWE5CzODx/LRfSMB');
    audio.play().catch(() => {});
  },
  
  // Вибрация (для устройств с поддержкой)
  vibrate: (pattern) => {
    if ('vibrate' in navigator) {
      navigator.vibrate(pattern);
    }
  },
  
  // Информация о батарее
  getBatteryInfo: async () => {
    if ('getBattery' in navigator) {
      try {
        const battery = await navigator.getBattery();
        return {
          charging: battery.charging,
          level: battery.level,
          chargingTime: battery.chargingTime,
          dischargingTime: battery.dischargingTime
        };
      } catch (error) {
        return null;
      }
    }
    return null;
  }
});

// Логирование загрузки preload скрипта
console.log('Preload script loaded successfully');

// Уведомляем главный процесс о готовности
window.addEventListener('DOMContentLoaded', () => {
  console.log('DOM loaded, preload script ready');
});