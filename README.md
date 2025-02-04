# 项目计划
+ [x] 后台管理
+ [x] 安全检测
+ [x] 文字生成
+ [x] 图片生成
+ [ ] 网页总结
+ [ ] 网页缩略图

# 安全检测项目
+ 防盗链 (Referer) 检测
+ IP滥用检测
+ 后台登录IP限制

# API本体使用
## 文字生成
格式：https://api.example.com/txt? + 参数1=值1&参数2=值2&参数3=值3，参数如下
| 参数名 | 参数值 | 说明 |
| ---- | ---- | ---- |
| format（选填） | json | 默认直接输出，json为输出JSON格式 |
| api（选填） | alibaba/openai/deepseek/otherapi | 在后台可以设置默认API |
| prompt（必填） | laugh, poem, sentence, 其他提示词 | laught为笑话，poem是诗句（创作），sentence是鸡汤 |

## 图片生成
格式：https://api.example.com/txt? + 参数1=值1&参数2=值2&参数3=值3，参数如下
| 参数名 | 参数值 | 说明 |
| ---- | ---- | ---- |
| format（选填） | json | 默认直接输出，json为输出JSON格式 |
| api（选填） | alibaba/openai | 在后台可以设置默认API |
| prompt（必填） | / | 提示词 |
| size（选填） | / | 后台可以设置默认大小，自定义的话参见对应的文档填入 |

参考文档：
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
有安全，文字，图像，网页四个选项，具体讲解如图：
![](https://raw.githubusercontent.com/stephen-zeng/urlAPI/master/guide.png)

# 启动设置
启动参数（可选）
+ clear - 清空任务
+ restore - 清空登录凭证，初始化设置
+ repwd - 重置后台登录密码
+ port .... - 将端口设置为....，默认端口是2233
反向代理的时候注意将发送域名设置为`$http_host`

# Demo
+ demo地址: http://urlapi.asdfhjkl.cn
+ dash密码是123456
+ 运用了该项目的文章demo：https://www.qwqwq.com.cn/urlapi
+ P.S. 项目release的时候本人的OpenAI余额用完了，然后充钱渠道暂时不可用，所以OpenAI的接口均为理论可行，实际未测试状态。使用demo的时候请注意不要留下任何敏感信息。