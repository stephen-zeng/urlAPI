# 项目功能计划
+ [x] 后台管理
+ [x] 安全检测
+ [x] 文字生成
+ [x] 图片生成
+ [x] 随机图片
+ [x] 网页缩略图

# 安全检测项目
+ 防盗链 (Referer) 检测
+ IP滥用检测
+ 后台登录IP限制

# API本体使用
## 文字生成
格式：https://api.example.com/txt? + 参数1=值1&参数2=值2&参数3=值3，参数如下
| 参数名 | 参数值 | 说明 |
| ---- | ---- | ---- |
| prompt（必填） | laugh, poem, sentence, 其他提示词 | laught为笑话，poem是诗句（创作），sentence是鸡汤 |
| format（选填） | json, txt | 默认直接跳转图片，json为输出JSON格式 |
| api（选填） | alibaba/openai/deepseek/otherapi | 在后台可以设置默认API |
| model（选填） | * | 在后台可以设置默认模型 |
| more（选填） | * | 添加更多信息，便于记录和筛选 |

+ 图片上使用的字体是得意黑:)

## 图片生成
格式：https://api.example.com/img? + 参数1=值1&参数2=值2&参数3=值3，参数如下
| 参数名 | 参数值 | 说明 |
| ---- | ---- | ---- |
| prompt（必填） | / | 提示词 |
| format（选填） | json | 默认直接跳转，json为输出JSON格式 |
| api（选填） | alibaba, openai | 在后台可以设置默认API |
| size（选填） | / | 后台可以设置默认大小，自定义的话参见对应的文档填入 |
| more（选填） | * | 添加更多信息，便于记录和筛选 |

## 随机图片
格式：https://api.example.com/rand? + 参数1=值1&参数2=值2&参数3=值3，参数如下
| 参数名 | 参数值 | 说明 |
| ---- | ---- | ---- |
| api（选填） | github, gitee | 使用的仓库 |
| user（必填） | / | 仓库主人的用户名 |
| repo（必填） | / | 仓库名称 |
| format（选填） | json | 默认直接跳转，json为输出JSON |
| more（选填） | * | 添加更多信息，便于记录和筛选 |

**请注意，仓库内不要出现文件夹，否则会添加失败**

## 网页缩略图
格式：https://api.example.com/web? + 参数1=值1&参数2=值2&参数3=值3，参数如下
| 参数名 | 参数值 | 说明 |
| ---- | ---- | ---- |
| img（必填） | URL | 目标页面地址（需URL编码） |
| format（选填） | json | 默认返回图片，json为返回JSON |
| more（选填） | * | 添加更多信息，便于记录和筛选 |

支持的地址格式：
+ Bilibili：`https://www.bilibili.com/video/BV.../` 或 `.../video/av...`
+ YouTube：`https://www.youtube.com/watch?v=...`
+ arXiv：`https://arxiv.org/abs/<论文编号>`
+ GitHub / Gitee 仓库：`https://github.com/<用户>/<仓库>`、`https://gitee.com/<用户>/<仓库>`
+ IT之家文章（需要启用文字总结）

`format=json` 时返回的 `url` 为绝对地址（见下方 `URLAPI_PUBLIC_URL`），否则 302 跳转。

## 参考文档
+ OpenAI - https://platform.openai.com/docs/api-reference/images/create
+ 阿里巴巴（通义万象） - https://help.aliyun.com/zh/model-studio/developer-reference/image-generation-wanx/

# 后台管理
地址：https://api.example.com/dash

## 登录密码
+ 新安装首次启动时，如果设置了 `URLAPI_ADMIN_PASSWORD` 则使用它作为初始密码；否则会随机生成一个密码并**只在启动日志中打印一次**（Docker 可用 `docker logs urlapi` 查看），请登录后立即修改。
+ 密码以 Argon2id 加盐哈希保存。旧版本保存的 SHA-256 密码仍然可以登录，并会在第一次成功登录后自动升级为 Argon2id，无需任何操作。
+ 忘记密码时运行 `urlAPI repwd`：会生成一个新的随机密码并在命令行打印一次，同时清除所有登录凭证（不再重置为 `123456`）。
+ “允许登录后台的IP”现在会被真正校验（默认 `*` 不限制）。如果把自己锁在外面，运行 `urlAPI clear_ip_restriction`。
## 访问情况
+ 上半部分为访问情况，单击记录可以查看具体信息，具体信息中单击对应信息可复制
+ 下半部分为分类，单击分类中具体项目可以筛选，刷新以重置

## 接口设置
这里设置对接OpenAI，阿里巴巴等的接口，包括默认模型选择，API Key等等。

## 功能设置
有安全，文字，图像，随机图片，网页五个选项，具体讲解如图（新版本可能会有一点差异）：
![](https://raw.githubusercontent.com/stephen-zeng/urlAPI/master/guide/1.png)
<img src="https://raw.githubusercontent.com/stephen-zeng/urlAPI/master/guide/2.png" width="300px"/>

# 启动设置
启动参数（可选）
+ clear - 清空任务
+ logout - 清空登录凭证
+ repwd - 生成新的随机后台密码并打印一次
+ port .... - 将端口设置为....，默认端口是2233
+ clear_ip_restriction - 清除后台登录IP限制

服务收到 `SIGINT`/`SIGTERM` 时会等待正在处理的请求结束后再关闭数据库。

## 环境变量
| 变量 | 说明 |
| ---- | ---- |
| `URLAPI_SECRET_KEY` | 32 字节主密钥（base64 或 64 位 hex），用于以 AES-256-GCM 加密保存接口 API Key 与 GitHub/YouTube Token。生成：`openssl rand -base64 32`。**请妥善备份，丢失后已加密的密钥无法恢复。** |
| `URLAPI_SECRET_KEY_FILE` | 从文件读取主密钥（适用于 Docker/systemd secrets），`URLAPI_SECRET_KEY` 优先 |
| `URLAPI_ADMIN_PASSWORD` | 仅在新安装首次启动时使用的初始后台密码 |
| `URLAPI_TRUSTED_PROXIES` | 允许设置 `X-Forwarded-For`/`X-Real-IP`/`X-Forwarded-Proto` 的反向代理 IP/CIDR，逗号分隔；默认信任本机与内网地址（`127.0.0.0/8, ::1, 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16, fc00::/7`），设为 `none` 则不信任任何代理 |
| `URLAPI_PUBLIC_URL` | 对外访问地址，例如 `https://api.example.com`，用于生成 JSON 中的绝对地址；未设置时使用经过校验的 `Host` 请求头 |
| `URLAPI_DASHBOARD_ORIGINS` | 允许跨域调用 `/session` 的后台来源（逗号分隔），一般无需设置 |

主密钥是可选的：未设置时程序照常运行，但 API Key 仍按旧格式（Base64）保存，并在日志中给出警告。设置后下次启动会自动把已有的密钥加密保存；如果之后主密钥缺失或错误，已加密的值会原样保留（不会被删除），只是暂时无法使用，恢复正确的主密钥即可。

## 反向代理
+ 前端容器（`urlapi-frontend`）或本机 Nginx 位于默认信任的内网/本机地址中，无需额外设置；其他部署请设置 `URLAPI_TRUSTED_PROXIES`，否则客户端 IP 将显示为代理地址。
+ 推荐同时设置 `URLAPI_PUBLIC_URL`。

# Demo
+ demo地址: 在简介里面
+ dash密码是123456
+ 运用了该项目的文章demo：https://www.qwqwq.com.cn/test/urlapi/

# 一些可能生成错误的原因
+ 服务器与上游API的连接问题，比如国内服务器不经过特殊手段无法连接OpenAI的服务器。
+ 还有出现的问题欢迎Issue

# Docker部署

可以使用项目中的`Dockerfile`自行构建，也可以使用下面的命令来运行
```bash
# 前后端统一部署
docker run -d --name urlapi -p %EXPOSE_PORT%:2233 -v %LOCAL_DATA_PLACE%:/app/assets -e URLAPI_SECRET_KEY=%SECRET_KEY% 0w0w0/urlapi:latest

# e.g
docker run -d --name urlapi -p 8080:2233 -v /home/stephenzeng/dockerData/urlAPI:/app/assets \
  -e URLAPI_SECRET_KEY="$(openssl rand -base64 32)" -e URLAPI_PUBLIC_URL=https://api.example.com 0w0w0/urlapi:latest

# 前端独立部署
docker run -d --name urlapi-frontend -p %EXPOSE_PORT%:80 -e BACKEND_URL=http://your-backend.url 0w0w0/urlapi-frontend:latest

# e.g
docker run -d --name urlapi-frontend -p 8080:80 -e BACKEND_URL=http://your-backend.url 0w0w0/urlapi-frontend:latest
```
+ 镜像目前`latest`和具体版本号两种tag，建议使用`latest`。
+ arm版本的镜像为`0w0w0/urlapi-arm`
+ 示例中的 `$(openssl rand -base64 32)` 只是演示，请生成一次并保存好，每次启动都要使用同一个主密钥。
+ 数据库使用 SQLite WAL 模式，`assets` 目录中会出现 `database.db-wal`、`database.db-shm`，备份时请停止服务或一并复制。数据目录请放在本地磁盘，不要使用 NFS/SMB。

# 升级说明
从旧版本升级无需手动迁移数据库：
+ 旧的后台密码继续有效，首次登录后自动升级为 Argon2id。
+ 设置 `URLAPI_SECRET_KEY` 后启动一次，已保存的 API Key/Token 会自动加密。
+ 后台登录 IP 白名单开始生效；通过反向代理访问时请确认 `URLAPI_TRUSTED_PROXIES` 覆盖了代理地址。
+ `/session` 不再允许任意跨域来源访问。
+ `/download` 只接受生成图片的 UUID 或 `empty`，其他参数返回 400。

# 二进制文件
可以到action里面寻找最新的build，也可以去release里面寻找稳定的build

# 注册为systemctl服务
下面是`/etc/systemd/system/urlAPI.service`模板
```
[Unit]
Description = urlAPI
After = network.target syslog.target
Wants = network.target

[Service]
Type = simple
WorkingDirectory = /root/urlAPI/
ExecStart = /root/urlAPI/urlAPI port 2233
Restart = on-failure


[Install]
WantedBy = multi-user.target
```
之后运行`systemctl daemon-reload && systemctl enable --now urlAPI`即可注册为开机启动的服务
