## MiHome

![favicon.ico](favicon.ico)

> 小爱同学控制电脑网易云音乐的播放

### 使用步骤

1. 进入[巴法](https://cloud.bemfa.com/tcp/topic.html?did=wyy005&v=1)官网注册账号并登录
2. 点击MQTT设备云
3. 新建主题，主题名为：自定义名字+005。例（wyy005）
4. 双击miHome.exe，第一次启动所在目录会创建config.properties文件
5. 配置config.properties文件cloudmusicPath为网易云音乐程序路径，key为巴法平台的私钥，topicName为主题名里的自定义名字
6. 配置完后再次启动miHome.exe，启动成功右下角会出现程序图标
7. 软件适配的是巴法平台的空调设备，打开手机上的米家，添加设备，选择巴法
8. 进入小爱同学训练计划自行添加
   - 关闭空调：暂停网易云播放
   - 打开空调：开始网易云播放
   - 空调上下扫风：网易云音量调大
   - 空调停止上下扫风：网易云音量调小
   - 空调左右扫风：下一首歌
   - 空调停止左右扫风：上一首歌
