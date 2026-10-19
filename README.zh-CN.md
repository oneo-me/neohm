# HA Climate

[English](README.md) | **简体中文**

HA Climate 是一个 macOS 菜单栏应用，用来控制 Home Assistant 中的空调（`climate.*` 实体）。它用 Go 编写，基于 [MyGo](https://github.com/egoist/mygo) 和 [Onevra](https://github.com/oneo-me/onevra)。

左键点击菜单栏图标打开面板。面板列出所有空调，选择其中一台即可打开它的控制界面。设置也在同一个面板中，从齿轮按钮打开。右键点击图标可以退出。

## 功能

- 列出所有空调，并显示模式、室内温度和目标温度。
- 控制电源、温度、模式、风速、风向、节能模式和睡眠，以每台空调支持的功能为准。
- 在菜单栏图标旁显示最近打开的空调的目标温度。

## 安装

目前还没有发布版本，请按[开发](#开发)中的命令从源码构建。需要 macOS 13 或更高版本、Go 1.27 或更高版本，以及 Mac 能访问到的 Home Assistant。

## 快速开始

1. 在 Home Assistant 中打开个人资料页，创建一个**长期访问令牌**。
2. 启动 HA Climate。完成 Home Assistant 配置之前，面板会直接打开设置；之后可以从面板的齿轮按钮或用 ⌘, 打开设置。
3. 填写地址（例如 `http://homeassistant.local:8123`）和令牌，然后点击**连接**。
4. 关闭设置，从列表中选择一台空调。

地址和最近打开的空调保存在 `~/.ha-go/config.json`，令牌保存在系统钥匙串中。

## 开发

Go 模块和 MyGo 打包配置都在 `source/` 中，请在该目录下运行 Go 和 MyGo 命令：

```sh
cd source
go tool mygo dev      # 运行，带热重载
go test ./...         # 无界面测试，不会打开窗口
go tool mygo build    # 在 ../dist 生成 .app 和 .dmg
```

代码布局见 [ARCHITECTURE.md](ARCHITECTURE.md)，协作规则见 [AGENTS.md](AGENTS.md)。
