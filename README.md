# 项目功能计划
+ [x] 后台管理
+ [x] 安全检测
+ [x] 文字生成
+ [x] 图片生成
+ [x] 随机图片
+ [ ] 网页总结
+ [ ] 网页缩略图

# 项目短期计划
+ [ ] 美化文字生成的返回图片
+ [ ] 文字返回错误图片
+ [ ] 测试平台

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
| format（选填） | json, txt | 默认直接输出图片，json为输出JSON格式, txt为输出文字（需要使用iframe） |
| api（选填） | alibaba/openai/deepseek/otherapi | 在后台可以设置默认API |

+ 图片上使用的字体是得意黑:)
+ v1.2及以上的版本支持图片输出以及参数值`txt`，v1.1及以下版本默认直接输出文字。

## 图片生成
格式：https://api.example.com/txt? + 参数1=值1&参数2=值2&参数3=值3，参数如下
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
| api（必填） | github, gitee | 使用的仓库 |
| user（必填） | / | 仓库主人的用户名 |
| repo（必填） | / | 仓库名称 |
| format（选填） | json | 默认直接跳转，json为输出JSON格式 |

**请注意，仓库内不要出现文件夹，否则会添加失败**

## 网页总结
格式：https://api.example.com/rand? + 参数1=值1&参数2=值2&参数3=值3，参数如下
| 参数名 | 参数值 | 说明 |
| ---- | ---- | ---- |
| img（选填） | URL | 网页缩略图 |
| sum（选填，暂未完成） | URL | 网页总结 |
| format（选填） | json | 默认返回图片，json为返回JSON |

**img和sum同时填写的话默认为img**

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
+ restore - 清空登录凭证，初始化设置
+ repwd - 重置后台登录密码
+ port .... - 将端口设置为....，默认端口是2233
+ clear_ip_restriction - 清除后台登录IP限制
+ update - 更新数据库（版本更新时必须运行一次）
反向代理的时候注意将发送域名设置为`$http_host`

注意，初始密码为`123456`，请注意及时修改

# Demo
+ demo地址: http://urlapi.asdfhjkl.cn
+ dash密码是123456
+ 运用了该项目的文章demo：https://www.qwqwq.com.cn/test/urlapi/
+ P.S. 项目release的时候本人的OpenAI余额用完了，然后充钱渠道暂时不可用，所以OpenAI的接口均为理论可行，实际未测试状态。使用demo的时候请注意不要留下任何敏感信息。

# 一些可能生成错误的原因
+ 同一个页面同时请求太多次，导致上游API发出`429 Too Many Requests`，请求过于频繁，包括上面的文章demo页面。解决办法：刷新重新发起请求请求
+ 服务器与上游API的连接问题，比如国内服务器不经过特殊手段无法连接OpenAI的服务器。

# 其他
+ 前端Vue的项目文件在`internal/router/template`目录下