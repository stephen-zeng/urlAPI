# urlAPI 课程作业完整文档

**作者：武汉大学开源软件与技术课程 2026**

## 1. 项目概述

urlAPI 是一个面向网页与接口调用场景的多功能 API 服务。项目以 Go 语言实现后端服务，使用 SQLite 持久化配置与任务记录，并通过 Vue/Vite 构建后台管理前端。系统提供文本生成、图片生成、随机图片、网页缩略图、下载中转、后台管理和安全防护能力。

## 2. 功能清单

- 文本生成：根据预置类型或自定义提示词调用外部模型接口生成文本。
- 图片生成：根据提示词调用 OpenAI 兼容或阿里云图像接口生成图片。
- 随机图片：从 GitHub 或 Gitee 仓库读取图片列表并随机返回图片。
- 网页缩略图：根据目标 URL 生成网页缩略图或文章封面图。
- 下载中转：对生成图片进行下载代理，屏蔽原始资源地址。
- 后台管理：提供登录、任务查询、功能开关、接口配置和安全配置能力。
- 安全防护：提供 Referer 检测、IP 限流、后台登录 IP 限制与异常过滤。

## 3. 技术栈

- 后端：Go、Gin、GORM、SQLite。
- 前端：Vue 3、Vite、MDUI。
- 存储：SQLite 数据库与本地图片资源目录。
- 部署：二进制运行或 Docker/Compose 部署。

## 4. 仓库结构

```text
cmd/                    命令行入口与管理命令
file/                   内嵌配置、字体和静态资源
internal/bootstrap/     服务依赖初始化与释放
internal/database/      数据库连接、模型和配置持久化
internal/llm/           外部大模型服务客户端
internal/op/            核心业务操作与任务生命周期
internal/server/        HTTP 路由、处理器与中间件
static/                 后台管理前端源码
util/                   工具函数、绘图、外部接口封装
whuthesis/              课程论文和演示材料
api&db.md               API 与数据库说明
OPEN_SOURCE_COMPLIANCE.md 开源合规说明
NOTICE                  项目声明与第三方组件提示
LICENSE                 GPL-3.0 许可证
```

## 5. 运行环境

- Go 1.26.2 或兼容版本。
- Node.js 与 npm，用于构建 `static` 前端。
- SQLite 运行依赖，Go 模块会通过依赖自动引入驱动。

## 6. 构建与运行

```bash
go mod download
go test ./...
go build ./...
go run . port 2233
```

前端构建命令如下：

```bash
cd static
npm ci
npm run build
```

服务默认监听 `2233` 端口。初始后台密码为 `123456`，首次部署后应立即修改。

## 7. 命令行参数

- `start` 或 `server`：启动 HTTP 服务。
- `port <端口>`：设置监听端口，默认 `2233`。
- `admin repwd` 或 `repwd`：重置后台密码为默认值并清空会话。
- `admin clear` 或 `clear`：清空任务记录。
- `admin logout` 或 `logout`：清空登录会话。
- `admin clear_ip_restriction` 或 `clear_ip_restriction`：清除后台登录 IP 限制。

## 8. API 概览

- `GET /txt`：文本生成接口，核心参数为 `prompt`、`api`、`model`、`format`、`more`。
- `GET /img`：图片生成接口，核心参数为 `prompt`、`api`、`model`、`size`、`format`、`more`。
- `GET /rand`：随机图片接口，核心参数为 `api`、`user`、`repo`、`format`、`more`。
- `GET /web`：网页缩略图接口，核心参数为 `img` 或 `url`、`format`、`more`。
- `GET /download`：图片下载中转接口，核心参数为 `img`。
- `POST /session`：后台会话与管理操作接口。

完整接口与数据库说明见 `api&db.md`。

## 9. 数据设计

系统使用 SQLite 保存服务配置、外部 API 配置、提示词、任务、会话、仓库缓存和 API Key 使用记录。启动时执行自动迁移，并将常用配置加载到内存以提升请求处理效率。

主要数据包括：

- 任务记录：保存访问来源、任务类型、状态、目标、结果、模型、尺寸等信息。
- 应用设置：保存安全策略、功能开关和默认模型配置。
- 会话记录：保存后台登录 token 与过期时间。
- 仓库缓存：保存随机图片仓库内容，减少重复拉取。
- API Key 与用量：保存接口密钥与使用统计。

## 10. 安全设计

- Referer 白名单：限制公开接口可被哪些站点引用。
- IP 频率限制：降低短时间高频访问带来的滥用风险。
- 后台 IP 限制：限制后台管理访问来源。
- 登录会话：使用 token 维护后台登录态，支持长期与短期会话。
- 异常过滤：允许对特定域名或信息跳过入库，减少无效数据。

## 11. 注释规范

本仓库新增和整理的关键 Go 文件采用 Doxygen 风格注释：

- 文件头使用 `@file`、`@brief`、`@author`、`@copyright`。
- 函数注释使用 `@brief` 描述用途。
- 参数使用 `@param` 描述含义。
- 返回值使用 `@return` 描述错误或结果语义。

项目作者标识统一为：**武汉大学开源软件与技术课程 2026**。

## 12. 测试与验证

当前仓库可通过以下命令验证：

```bash
go test ./...
go build ./...
cd static && npm ci && npm run build
```

验证范围包括 Go 后端编译、Go 单元测试、前端依赖安装和 Vite 前端构建。

## 13. 开源许可证

本项目使用 GPL-3.0 许可证，完整文本见 `LICENSE`。选择 GPL-3.0 的原因是该项目包含完整可运行服务，强 copyleft 许可证可以保证派生版本在分发时继续提供源代码与相同自由。

分发源码、二进制、Docker 镜像或前端构建产物时，应保留 `LICENSE`、`NOTICE` 和必要的第三方组件许可证说明。详细合规说明见 `OPEN_SOURCE_COMPLIANCE.md`。

## 14. Git 仓库地址递交

课程作业提交时递交 Git 仓库地址即可。建议提交格式如下：

```text
项目名称：urlAPI
作者：武汉大学开源软件与技术课程 2026
仓库地址：<填写实际 Git 仓库 URL，例如 https://github.com/<owner>/<repo>>
许可证：GPL-3.0
完整文档：COURSE_REPORT.md
```

提交前建议确认远程地址：

```bash
git remote -v
```

## 15. 后续改进方向

- 为核心业务处理器补充更多单元测试和集成测试。
- 在发布流程中自动生成第三方依赖许可证清单。
- 增加 API 文档自动生成流程，使接口说明与代码保持同步。
- 增加数据库备份与迁移版本管理机制。
