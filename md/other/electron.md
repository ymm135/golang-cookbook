- # Electron    

目录:
- [小试牛刀](#小试牛刀)
  - [简介](#简介)
  - [环境搭建](#环境搭建)
- [疑问拓展](#疑问拓展)
- [vscode开发工具是`Electron`？](#vscode开发工具是electron)



## 小试牛刀 

`Electron` 基于 Chromium，而 Chromium 是谷歌开发的一个开源浏览器项目，它同样包括了一个网页渲染引擎。Chromium 最初使用的是 `WebKit`，但谷歌后来开发了自己的一个分支称为 `Blink`。从2013年开始，Chromium 和 Chrome 都转而使用 Blink 作为其网页渲染引擎。
> Electron `/ɪˈlektrɒn/`  电子  
> Chromium `/ˈkrəʊmiəm/`  [化学]铬（24号元素，符号 Cr）  

WebKit 是一个开源的网页渲染引擎，它负责解析 HTML、CSS 代码和执行 JavaScript，从而在屏幕上呈现网页内容。WebKit 是许多浏览器的基础，包括早期的 Apple Safari 和以前的 Google Chrome。

因此，`Electron` 与 `WebKit` 的关系在于它们都是网页渲染引擎，但 Electron 实际上是使用 Chromium 的 Blink 引擎，而不是直接使用 WebKit。不过，由于 Blink 引擎是从 WebKit 分支出来的，两者在许多核心渲染方面共享了相似的技术和代码基础。

> https://github.com/electron/electron  


如果您是在说使用 Web 技术（HTML, CSS, JavaScript）来创建一个具有图形用户界面的本地 Windows 应用程序，并且想要在点击按钮时显示一个“Hello”对话框，那么通常这涉及到使用 Electron 或 NW.js 这样的框架。这些框架允许开发人员使用 Web 技术来构建跨平台的桌面应用程序。

这里有一个基本的 Electron 应用程序示例，它在点击按钮时会显示一个弹出对话框：

1. **初始化一个新的 Node.js 项目**：

```bash
mkdir my-electron-app
cd my-electron-app
npm init -y
npm install electron --save-dev
```

2. **创建主进程文件** (`index.js`):

```javascript
const { app, BrowserWindow, dialog } = require('electron');

function createWindow() {
  const win = new BrowserWindow({
    width: 800,
    height: 600,
    webPreferences: {
      nodeIntegration: true
    }
  });

  win.loadFile('index.html');

  // 打开开发者工具（如果需要）
  win.webContents.openDevTools();
}

app.whenReady().then(createWindow);

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  }
});
```

3. **创建前端 HTML 文件** (`index.html`):

```html
<!DOCTYPE html>
<html>
  <head>
    <title>Hello Dialog</title>
  </head>
  <body>
    <h1>My Electron App</h1>
    <button id="helloButton">Say Hello</button>

    <script>
      // 在渲染进程中使用 Electron 的 API
      const { dialog } = require('electron').remote;

      document.getElementById('helloButton').addEventListener('click', () => {
        dialog.showMessageBox({
          message: 'Hello!',
          buttons: ['OK']
        });
      });
    </script>
  </body>
</html>
```

4. **在 `package.json` 中设置启动脚本**:

在 `package.json` 文件的 `scripts` 部分添加：

```json
"start": "electron ."
```

完成以上步

骤后，你可以通过在命令行中运行 `npm start` 来启动你的 Electron 应用程序。点击按钮时，会调用 Electron 的 `dialog` API 弹出一个包含“Hello”消息的对话框。

请注意，由于 Electron 5.0.0 版本开始，渲染器进程中默认不再启用 Node.js 集成。如果你需要在渲染器进程中使用 Node.js 特性（如上面示例中所做的那样），你需要在 `BrowserWindow` 的 `webPreferences` 中设置 `nodeIntegration: true`。此外，为了提高应用程序的安全性，应该避免在可能的情况下在渲染器进程中启用 Node.js 集成。

最后，如果你指的是使用 WebKit 直接开发 Windows 应用程序，那么这通常涉及到更底层的编程，例如使用 C++ 和一个类似于 Qt 的框架来嵌入 WebKit。然而，Electron 已经成为了使用 Web 技术创建桌面应用程序的主流方式，因为它简化了开发过程，并且提供了丰富的 API。

### 简介

### 环境搭建


## 疑问拓展  
## vscode开发工具是`Electron`？

Visual Studio Code（VSCode）是一个由微软开发的开源代码编辑器。它是使用 Electron 框架开发的，Electron 允许使用前端技术如 HTML、CSS 和 JavaScript 来开发跨平台的桌面应用程序。

VSCode 的开发涉及多种技术和工具：

1. **Electron**: 用于将 VSCode 作为一个桌面应用程序打包和运行。

2. **Node.js**: 提供后端运行时环境，允许使用 JavaScript 进行系统级的操作。

3. **TypeScript**: VSCode 的主要编程语言，是 JavaScript 的一个超集，添加了静态类型检查和更高级的编程特性。

4. **Monaco Editor**: 作为 VSCode 的编辑器核心，是一个用于网页应用的代码编辑器，也由微软开发。

5. **Git**: 用于版本控制，VSCode 自身也提供了内置的 Git 支持。

6. **各种前端技术**: 包括 HTML、CSS 和 JavaScript，用于构建用户界面。

7. **npm**: 作为包管理器，用于管理 VSCode 的依赖。

8. **Azure DevOps**: 微软的持续集成和持续部署服务，用于 VSCode 的开发流程。

由于 VSCode 是开源的，你可以在其[GitHub 仓库](https://github.com/microsoft/vscode)中找到所有源代码和构建脚本。这不仅让人们可以自由地探索和学习 VSCode 是如何被构建的，也允许社区贡献代码和功能。