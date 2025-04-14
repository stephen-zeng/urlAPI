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

+ 图片上使用的字体是得意黑:)

## 图片生成
格式：https://api.example.com/img? + 参数1=值1&参数2=值2&参数3=值3，参数如下
| 参数名 | 参数值 | 说明 |
| ---- | ---- | ---- |
| prompt（必填） | / | 提示词 |
| format（选填） | json | 默认直接跳转，json为输出JSON格式 |
| api（选填） | alibaba, openai | 在后台可以设置默认API |
| size（选填） | / | 后台可以设置默认大小，自定义的话参见对应的文档填入 |

## 随机图片
格式：https://api.example.com/rand? + 参数1=值1&参数2=值2&参数3=值3，参数如下
| 参数名 | 参数值 | 说明 |
| ---- | ---- | ---- |
| api（选填） | github, gitee | 使用的仓库 |
| user（必填） | / | 仓库主人的用户名 |
| repo（必填） | / | 仓库名称 |
| format（选填） | json | 默认直接跳转，json为输出JSON |

**请注意，仓库内不要出现文件夹，否则会添加失败**

## 网页缩略图
格式：https://api.example.com/rand? + 参数1=值1&参数2=值2&参数3=值3，参数如下
| 参数名 | 参数值 | 说明 |
| ---- | ---- | ---- |
| img（必填） | URL | 网页缩略图 |
| format（选填） | json | 默认返回图片，json为返回JSON |

## 参考文档
+ OpenAI - https://platform.openai.com/docs/api-reference/images/create
+ 阿里巴巴（通义万象） - https://help.aliyun.com/zh/model-studio/developer-reference/image-generation-wanx/

# 后台管理
地址：https://api.example.com/dash
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
+ repwd - 重置后台登录密码
+ port .... - 将端口设置为....，默认端口是2233
+ clear_ip_restriction - 清除后台登录IP限制
反向代理的时候注意将发送域名设置为`$http_host`

注意，初始密码为`123456`，请注意及时修改

# Demo
+ demo地址: 在简介里面
+ dash密码是123456
+ 运用了该项目的文章demo：https://www.qwqwq.com.cn/test/urlapi/

# 一些可能生成错误的原因
+ 服务器与上游API的连接问题，比如国内服务器不经过特殊手段无法连接OpenAI的服务器。
+ 还有出现的问题欢迎Issue
