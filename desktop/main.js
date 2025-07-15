const { app, BrowserWindow, Menu, Tray, ipcMain, shell, dialog, nativeTheme } = require('electron');
const { autoUpdater } = require('electron-updater');
const Store = require('electron-store');
const path = require('path');
const contextMenu = require('electron-context-menu');
const AutoLaunch = require('auto-launch');

// Инициализация хранилища настроек
const store = new Store({
  defaults: {
    windowBounds: { width: 1200, height: 800 },
    theme: 'system',
    notifications: true,
    autoLaunch: false,
    minimizeToTray: true,
    closeToTray: true
  }
});

// Глобальные переменные
let mainWindow = null;
let tray = null;
let isQuiting = false;

// Настройка автозапуска
const autoLauncher = new AutoLaunch({
  name: 'Messenger',
  path: app.getPath('exe')
});

// Включаем контекстное меню
contextMenu({
  showLookUpSelection: false,
  showSearchWithGoogle: false,
  showCopyImage: true,
  showCopyImageAddress: true,
  showSaveImage: true,
  showSaveImageAs: true,
  showInspectElement: process.env.NODE_ENV === 'development'
});

// Создание главного окна
function createMainWindow() {
  const bounds = store.get('windowBounds');
  
  mainWindow = new BrowserWindow({
    width: bounds.width,
    height: bounds.height,
    minWidth: 800,
    minHeight: 600,
    icon: path.join(__dirname, 'assets', 'icon.png'),
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      enableRemoteModule: false,
      preload: path.join(__dirname, 'preload.js')
    },
    titleBarStyle: process.platform === 'darwin' ? 'hiddenInset' : 'default',
    show: false
  });

  // Загружаем приложение
  if (process.env.NODE_ENV === 'development') {
    mainWindow.loadFile('src/index.html');
    mainWindow.webContents.openDevTools();
  } else {
    mainWindow.loadFile('src/index.html');
  }

  // Показываем окно когда оно готово
  mainWindow.once('ready-to-show', () => {
    mainWindow.show();
    
    // Автообновления
    if (process.env.NODE_ENV !== 'development') {
      autoUpdater.checkForUpdatesAndNotify();
    }
  });

  // Сохраняем размер окна
  mainWindow.on('resize', () => {
    store.set('windowBounds', mainWindow.getBounds());
  });

  // Обработка закрытия окна
  mainWindow.on('close', (event) => {
    if (!isQuiting && store.get('closeToTray')) {
      event.preventDefault();
      mainWindow.hide();
    }
  });

  // Обработка минимизации
  mainWindow.on('minimize', (event) => {
    if (store.get('minimizeToTray')) {
      event.preventDefault();
      mainWindow.hide();
    }
  });

  return mainWindow;
}

// Создание трея
function createTray() {
  const trayIcon = path.join(__dirname, 'assets', 'tray-icon.png');
  tray = new Tray(trayIcon);
  
  const contextMenu = Menu.buildFromTemplate([
    {
      label: 'Показать',
      click: () => {
        mainWindow.show();
      }
    },
    { type: 'separator' },
    {
      label: 'Настройки',
      click: () => {
        mainWindow.show();
        mainWindow.webContents.send('show-settings');
      }
    },
    { type: 'separator' },
    {
      label: 'Выход',
      click: () => {
        isQuiting = true;
        app.quit();
      }
    }
  ]);

  tray.setToolTip('Messenger');
  tray.setContextMenu(contextMenu);
  
  // Двойной клик по трею показывает окно
  tray.on('double-click', () => {
    mainWindow.show();
  });
}

// Создание меню приложения
function createMenu() {
  const template = [
    {
      label: 'Файл',
      submenu: [
        {
          label: 'Новый чат',
          accelerator: 'CmdOrCtrl+N',
          click: () => {
            mainWindow.webContents.send('new-chat');
          }
        },
        { type: 'separator' },
        {
          label: 'Настройки',
          accelerator: 'CmdOrCtrl+,',
          click: () => {
            mainWindow.webContents.send('show-settings');
          }
        },
        { type: 'separator' },
        {
          label: 'Выход',
          accelerator: process.platform === 'darwin' ? 'Cmd+Q' : 'Ctrl+Q',
          click: () => {
            isQuiting = true;
            app.quit();
          }
        }
      ]
    },
    {
      label: 'Правка',
      submenu: [
        { role: 'undo', label: 'Отменить' },
        { role: 'redo', label: 'Повторить' },
        { type: 'separator' },
        { role: 'cut', label: 'Вырезать' },
        { role: 'copy', label: 'Копировать' },
        { role: 'paste', label: 'Вставить' },
        { role: 'selectall', label: 'Выделить всё' }
      ]
    },
    {
      label: 'Вид',
      submenu: [
        {
          label: 'Перезагрузить',
          accelerator: 'CmdOrCtrl+R',
          click: () => {
            mainWindow.reload();
          }
        },
        {
          label: 'Полноэкранный режим',
          accelerator: process.platform === 'darwin' ? 'Ctrl+Cmd+F' : 'F11',
          click: () => {
            mainWindow.setFullScreen(!mainWindow.isFullScreen());
          }
        },
        { type: 'separator' },
        {
          label: 'Увеличить',
          accelerator: 'CmdOrCtrl+Plus',
          click: () => {
            mainWindow.webContents.send('zoom-in');
          }
        },
        {
          label: 'Уменьшить',
          accelerator: 'CmdOrCtrl+-',
          click: () => {
            mainWindow.webContents.send('zoom-out');
          }
        },
        {
          label: 'Сбросить масштаб',
          accelerator: 'CmdOrCtrl+0',
          click: () => {
            mainWindow.webContents.send('zoom-reset');
          }
        }
      ]
    },
    {
      label: 'Помощь',
      submenu: [
        {
          label: 'О программе',
          click: () => {
            dialog.showMessageBox(mainWindow, {
              type: 'info',
              title: 'О программе',
              message: 'Messenger Desktop',
              detail: `Версия: ${app.getVersion()}\nElectron: ${process.versions.electron}\nNode.js: ${process.versions.node}`
            });
          }
        },
        {
          label: 'Открыть DevTools',
          accelerator: process.platform === 'darwin' ? 'Alt+Cmd+I' : 'Ctrl+Shift+I',
          click: () => {
            mainWindow.webContents.toggleDevTools();
          }
        }
      ]
    }
  ];

  // Особенности меню для macOS
  if (process.platform === 'darwin') {
    template.unshift({
      label: app.getName(),
      submenu: [
        { role: 'about', label: 'О программе' },
        { type: 'separator' },
        { role: 'services', label: 'Службы' },
        { type: 'separator' },
        { role: 'hide', label: 'Скрыть' },
        { role: 'hideothers', label: 'Скрыть другие' },
        { role: 'unhide', label: 'Показать все' },
        { type: 'separator' },
        { role: 'quit', label: 'Выйти' }
      ]
    });
  }

  const menu = Menu.buildFromTemplate(template);
  Menu.setApplicationMenu(menu);
}

// IPC обработчики
ipcMain.handle('get-store-value', (event, key) => {
  return store.get(key);
});

ipcMain.handle('set-store-value', (event, key, value) => {
  store.set(key, value);
});

ipcMain.handle('get-app-version', () => {
  return app.getVersion();
});

ipcMain.handle('show-message-box', (event, options) => {
  return dialog.showMessageBox(mainWindow, options);
});

ipcMain.handle('show-open-dialog', (event, options) => {
  return dialog.showOpenDialog(mainWindow, options);
});

ipcMain.handle('show-save-dialog', (event, options) => {
  return dialog.showSaveDialog(mainWindow, options);
});

ipcMain.handle('set-badge-count', (event, count) => {
  if (process.platform === 'darwin') {
    app.setBadgeCount(count);
  }
});

ipcMain.handle('show-notification', (event, options) => {
  const { Notification } = require('electron');
  
  if (store.get('notifications') && Notification.isSupported()) {
    const notification = new Notification({
      title: options.title,
      body: options.body,
      icon: options.icon || path.join(__dirname, 'assets', 'icon.png')
    });
    
    notification.show();
    
    notification.on('click', () => {
      mainWindow.show();
      mainWindow.focus();
    });
  }
});

ipcMain.handle('set-auto-launch', async (event, enabled) => {
  try {
    if (enabled) {
      await autoLauncher.enable();
    } else {
      await autoLauncher.disable();
    }
    store.set('autoLaunch', enabled);
    return true;
  } catch (error) {
    console.error('Auto launch error:', error);
    return false;
  }
});

ipcMain.handle('get-auto-launch', async () => {
  try {
    return await autoLauncher.isEnabled();
  } catch (error) {
    console.error('Auto launch check error:', error);
    return false;
  }
});

// События приложения
app.whenReady().then(() => {
  createMainWindow();
  createTray();
  createMenu();
  
  // Обработка активации на macOS
  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createMainWindow();
    } else {
      mainWindow.show();
    }
  });
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('before-quit', () => {
  isQuiting = true;
});

// Обработка протокола для глубоких ссылок
app.setAsDefaultProtocolClient('messenger');

// Автообновления
autoUpdater.on('checking-for-update', () => {
  console.log('Checking for update...');
});

autoUpdater.on('update-available', (info) => {
  console.log('Update available.');
});

autoUpdater.on('update-not-available', (info) => {
  console.log('Update not available.');
});

autoUpdater.on('error', (err) => {
  console.log('Error in auto-updater. ' + err);
});

autoUpdater.on('download-progress', (progressObj) => {
  let log_message = "Download speed: " + progressObj.bytesPerSecond;
  log_message = log_message + ' - Downloaded ' + progressObj.percent + '%';
  log_message = log_message + ' (' + progressObj.transferred + "/" + progressObj.total + ')';
  console.log(log_message);
});

autoUpdater.on('update-downloaded', (info) => {
  console.log('Update downloaded');
  autoUpdater.quitAndInstall();
});

// Предотвращение множественных экземпляров
const gotTheLock = app.requestSingleInstanceLock();

if (!gotTheLock) {
  app.quit();
} else {
  app.on('second-instance', (event, commandLine, workingDirectory) => {
    // Пользователь попытался запустить второй экземпляр
    if (mainWindow) {
      if (mainWindow.isMinimized()) mainWindow.restore();
      mainWindow.focus();
    }
  });
}