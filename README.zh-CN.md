# HA Climate

一个 macOS 菜单栏应用，用来控制 Home Assistant 里的空调（`climate.*` 实体）。它用 Go 编写，基于 [MyGo](https://github.com/egoist/mygo) 和 [Onevra](https://github.com/oneo-me/onevra)。

菜单栏图标旁会显示目标温度。左键点击图标打开控制面板，右键点击可以刷新、打开设置或退出。

## 配置

1. 在 Home Assistant 中打开个人资料页，创建一个**长期访问令牌**。
2. 打开 **设置…**，填写地址（例如 `http://homeassistant.local:8123`）和令牌。
3. 点击 **连接**，然后选择空调。

地址和实体保存在 `~/Library/Application Support/ha-climate/config.json`，令牌保存在系统钥匙串中。

## 面板功能

| 控件 | Home Assistant 服务 |
| --- | --- |
| 电源 | `climate.turn_off` / `climate.set_hvac_mode` |
| 温度滑块、− 和 + | `climate.set_temperature`（滑块停下后才发送） |
| 模式 | `climate.set_hvac_mode` |
| 风速 | `climate.set_fan_mode` |
| 风向 | `climate.set_swing_mode` |
| 节能模式、睡眠 | `climate.set_preset_mode`，设为 `eco` / `sleep`，关闭时恢复为 `none` |

实体不支持的控件会显示为禁用。面板打开时每 3 秒读取一次状态，关闭时每 20 秒读取一次。

## 开发

```sh
cd source
go tool mygo dev      # 运行，带热重载
go test ./...         # 无界面测试，不会打开窗口
HA_SNAPSHOT_DIR=/tmp go test -run Snapshot .   # 把面板渲染为 PNG
go tool mygo build    # 在 ../dist 生成 .app 和 .dmg
```
