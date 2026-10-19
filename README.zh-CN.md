# NeoHM

[English](README.md) | **简体中文**

NeoHM 是一个 macOS 菜单栏应用，用来控制 Home Assistant 中的设备，并为空调提供完整的控制界面。它用 Go 编写，基于 [MyGo](https://github.com/egoist/mygo) 和 [Onevra](https://github.com/oneo-me/onevra)。

左键点击菜单栏图标打开面板。面板按房间以方块显示设备，设备开启时方块会亮起。点击方块可以开关设备，方块角上的箭头会打开设备的页面；不能开关的设备点击方块即打开它的页面。设置也在同一个面板中，从齿轮按钮打开。右键点击图标可以退出 NeoHM，面板打开时也可以按 ⌘Q 退出。

## 功能

- 按房间以方块显示设备，每个方块显示设备的主要状态以及是否开启。
- 在方块上开关灯、插座、风扇、窗帘等设备，并在设备的页面中调光、调色、调节位置、播放、锁门或启动设备，一键选择灯光颜色或窗帘的四分位置。
- 控制空调的电源、温度、模式、风速、风向、节能模式和睡眠，以每台空调支持的功能为准。
- 以英文、简体中文、繁体中文、日文或韩文显示界面，跟随 macOS 或在设置中选择的语言。

## 安装

目前还没有发布版本，请按[开发](#开发)中的命令从源码构建。需要 macOS 13 或更高版本、Go 1.27 或更高版本，以及 Mac 能访问到的 Home Assistant。

## 快速开始

1. 在 Home Assistant 中打开个人资料页，创建一个**长期访问令牌**。
2. 启动 NeoHM。完成 Home Assistant 配置之前，面板会直接打开设置；之后可以从面板的齿轮按钮或用 ⌘, 打开设置。
3. 填写地址（例如 `http://homeassistant.local:8123`）和令牌，然后点击**连接**。
4. 关闭设置即可看到你的设备。点击方块可以开关设备，点击它的箭头可以打开设备的页面。

地址、令牌、所选语言和最近打开的空调保存在 `~/.neohm/config.json`，只有你的用户账户可以读取。

## 开发

Go 模块和 MyGo 打包配置都在 `source/` 中，请在该目录下运行 Go 和 MyGo 命令：

```sh
cd source
go tool mygo dev      # 运行，带热重载
go test ./...         # 无界面测试，不会打开窗口
go tool mygo build    # 在 ../dist 生成 .app 和 .dmg
```

代码布局见 [ARCHITECTURE.md](ARCHITECTURE.md)，协作规则见 [AGENTS.md](AGENTS.md)。
