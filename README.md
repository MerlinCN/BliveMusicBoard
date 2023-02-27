# 哔哩哔哩直播点歌板

这是一个点歌板项目，可用于OBS中，支持用户点歌、删除歌曲、查看历史歌单以及管理员设置点歌状态等功能。此项目包含了前端和后端两部分，其中前端使用了 Vue.js 框架实现，后端则是使用 Go 语言。

## 功能

- 点歌：用户可以通过弹幕来点歌，并将歌曲添加到歌曲列表中。
- 删除歌曲：管理员可以删除歌曲列表中的歌曲。
- 同步歌曲列表：所有客户端将实时同步歌曲列表。

## 安装

1. 安装 Go 和 Git
2. 克隆此仓库：`git clone https://github.com/MerlinCN/BliveMusicBoard.git`
3. 下载编译好的前端文件
4. 进入项目目录：`cd BliveMusicBoard`
5. 安装依赖：`go mod tidy`

## 运行

1. 启动服务器：`go run main.go`
2. 在浏览器中打开 `http://localhost:4220`

## 代码结构

- `main.go`：入口文件，定义了 HTTP 服务器和 WebSocket 服务器。

- ```
  pkg
  ```

   目录：包含了项目的核心代码。

  - `song.go`：定义了 `Song` 结构体和相关的方法，用于表示歌曲。
  - `utils.go`：定义了一些辅助结构体和常量。
  - `websocket.go`：定义了 WebSocket 相关的方法，用于与客户端通信。

## 技术栈

- Go
- HTML/CSS/JavaScript
- Vue.js
- Gorilla WebSocket

## 许可证

本项目使用AGPL-3.0 许可证。详情请参阅 LICENSE 文件。
