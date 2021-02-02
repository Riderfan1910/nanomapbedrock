const { join } = require('path')

const { BrowserWindow, app } = require('electron')
const isDev = !app.isPackaged

let window = null

function createWindow () {
  window = new BrowserWindow({
    width: 960,
    height: 600,
    webPreferences: {
      webSecurity: false,
      nodeIntegration: true,
      enableRemoteModule: true
    }
  })

  const port = process.env.PORT || 3000
  const url = isDev
    ? `http://localhost:${port}`
    : join(__dirname, './dist/index.html')

  isDev ? window.loadURL(url) : window.loadFile(url)

  // window.webContents.openDevTools();
}

app.on('ready', createWindow)

app.on('activate', () => {
  if (!window) {
    createWindow()
  }
})

app.on('window-all-closed', function () {
  if (process.platform !== 'darwin') app.quit()
})
