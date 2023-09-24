## ZxwyWebSite/GoogTransTTS
### 简介
+ 使用谷歌搜索翻译Api的TTS工具
+ 开发中，部分功能暂未完善

### 使用
```shell
./googtranstts \
    -lang zh-CN \ # 源语言
    -proxy socks5://user:pswd@addr:port \ # 使用代理
    -text example/text3.txt \ # 要翻译的文本
    -usetxt \ # 从text参数解析文本文档
    -name 第三章_隐藏巫师有一把石锤。 \ # 输出文件名
    -format # 混合分段音频 (需要FFmpeg命令)
```

### 注意
+ 换行不会自动断句，请自行在末尾添加 "." (英文句号)

### 其它
+ 暂无

### 更新
#### 2023-09-24 v0.1
+ 上传项目
