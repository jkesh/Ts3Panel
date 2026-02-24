# Ts3Panel - TeamSpeak 3 Server Management Panel

Ts3Panel 是一个基于 Go (后端) 和 Vue 3 (前端) 构建的现代化 TeamSpeak 3 服务器管理面板。它提供了一个直观的 Web 界面，用于管理 TS3 服务器实例、频道、权限组、用户封禁等核心功能。

![Logo](https://github.com/user-attachments/assets/b05e3595-6787-4751-9397-9a9e979fb7d7)

## ✨ 功能特性

* **仪表盘 (Dashboard)**: 实时监控服务器状态（在线人数、版本、平台、运行时间）以及实时日志流 (SSE)。
* **用户管理**:
    * 查看在线用户列表。
    * 踢出用户 (Kick)。
    * 发送全服广播消息。
* **频道管理**:
    * 创建新频道（支持设置密码、话题）。
    * 删除频道（支持强制删除）。
    * 快捷修改频道权限。
* **权限组管理**:
    * 查看服务器组列表。
    * 生成权限组密钥 (Privilege Key/Token)。
    * 修改服务器组权限。
    * 删除服务器组。
* **封禁管理**:
    * 查看封禁列表。
    * 添加封禁（支持 IP、UID、昵称正则）。
    * 解除封禁 / 清空封禁列表。
* **安全认证**:
    * 基于 JWT 的身份验证。
    * 管理员注册与登录。

## 🛠️ 技术栈

### 后端 (Backend)
* **语言**: Go 1.25+
* **Web 框架**: [Gin](https://github.com/gin-gonic/gin)
* **ORM**: [GORM](https://gorm.io/) (支持 PostgreSQL & SQLite)
* **TS3 库**: [ts3-go](https://github.com/jkesh/ts3-go)（基于 v2.0.0 标签对应实现）
* **配置管理**: [Viper](https://github.com/spf13/viper)
* **鉴权**: JWT (JSON Web Tokens)

### 前端 (Frontend)
* **框架**: Vue 3 (Composition API)
* **构建工具**: Vite
* **UI 组件库**: Element Plus
* **路由**: Vue Router
* **状态管理**: Pinia
* **HTTP 客户端**: Axios
* **实时通信**: Server-Sent Events (SSE)

## 📋 环境要求

* **Go**: 版本 1.20 或更高
* **Node.js**: 版本 16 或更高 (推荐使用 LTS)
* **数据库**: PostgreSQL (推荐) 或 SQLite
* **TeamSpeak 3 Server**: 需开启 ServerQuery 功能

## 🚀 快速开始

### 1. 后端部署

1.  进入后端目录：
    ```bash
    cd backend
    ```

2.  下载依赖：
    ```bash
    go mod download
    ```

3.  配置应用：
    修改 `config.yaml` 文件，填入你的数据库信息和 TS3 ServerQuery 凭据。

    ```yaml
    # backend/config.yaml
    app:
      jwt_secret: "your_secure_jwt_secret" # 修改为随机字符串
      port: ":8080"

    ts3:
      protocol: "tcp"       # "tcp" | "ssh" | "webquery"
      host: "127.0.0.1"     # TS3 服务器 IP
      port: 10011           # tcp:10011, ssh:10022, webquery:10080/10443
      server_id: 1          # 虚拟服务器 ID (use sid)
      user: "serveradmin"   # tcp/ssh Query 账号
      password: "your_password" # tcp/ssh Query 密码
      api_key: ""           # webquery 模式使用
      https: false          # webquery + TLS 设为 true
      base_path: ""         # 可选，如 "query"

    database:
      driver: "postgres"    # "postgres" 或 "sqlite"
      host: "127.0.0.1"
      port: 5432
      user: "postgres"
      password: "db_password"
      dbname: "ts3panel"
      sslmode: "disable"
      timezone: "Asia/Shanghai"
    ```

4.  运行后端：
    ```bash
    go run main.go
    ```
    后端服务将在 `http://localhost:8080` 启动。

### 2. 前端部署

1.  进入前端目录：
    ```bash
    cd frontend
    ```

2.  安装依赖：
    ```bash
    npm install
    ```

3.  开发模式运行：
    ```bash
    npm run dev
    ```
    访问终端输出的地址（通常是 `http://localhost:5173`）。

4.  生产环境构建：
    ```bash
    npm run build
    ```
    构建生成的文件位于 `frontend/dist` 目录。你可以将这些文件部署到 Nginx，或者配置后端 Go 服务来托管静态文件。


## 📝 注意事项

* **首次使用**: 启动项目后，请访问 `/register` 页面注册第一个管理员账号。
* **权限映射**: 前端 `utils/permMap.js` 包含了常用权限 ID 到名称的映射。如果遇到权限名显示为 ID 的情况，可以在此文件中补充。
* **ServerQuery 限制**: 请确保你的 TS3 服务器 `ip_whitelist.txt` 中包含了运行本面板的服务器 IP，以避免因请求过多被封禁 (Flood Ban)。

## 🤝 贡献

欢迎提交 Issue 或 Pull Request 来改进 Ts3Panel！

## 📄 许可证

MIT License







