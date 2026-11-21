# 上班鼓励器

一个用 Go 写的 Windows 桌面小工具：把你的月薪折算成每秒进账，看着钱一点点挣回来。
复刻自「上班鼓励器」应用，兼容 **Windows 7 / 8 / 10 / 11**，单文件 exe，双击即开，无需安装任何运行库。

## 功能

- 🕐 实时时钟 + 下班倒计时 + 「上班中 / 下班啦」状态
- 💰 今日已赚金额按秒跳动，显示赚钱速度（元/秒），按月薪 ÷ 21.75 天 ÷ 每日工时折算时薪
- 🗣 随机鼓励语：内置 57 条商品语录（一瓶冰可乐 → 一部 iPhone → 一个金镯子），每 12 秒随机换一条，按你的赚钱速度实时算出「再过 X 小时 X 分，XX 就挣回来了」
- 📌 窗口置顶 / 取消置顶（原生 Win32 SetWindowPos）
- ⏸ 关闭计时器：暂停/恢复，暂停时时间冻结，恢复后接着算
- 📊 今日工时进度条 + 已工作小时数
- 📅 本月累计：已过工作日按满勤补算
- ⚙ 设置面板：月薪、每天工作时长、上班/下班时间，保存到 exe 同目录 `config.json`

## 下载使用

在 [Releases](../../releases) 中下载对应版本：

| 文件 | 适用系统 |
|---|---|
| `上班鼓励器.exe` | 64 位 Windows |
| `上班鼓励器_32位.exe` | 32 位 Windows |

> 请放在有写权限的目录（如桌面、D 盘），否则设置无法保存。
> 删除 `config.json` 即恢复默认值（默认：月薪 4000 元、每天 8 小时、09:00–18:00）。

## 从源码构建

⚠️ 必须使用 **Go 1.20 及以下版本**编译（Go 1.21+ 生成的 exe 无法在 Win7 运行）。

```bash
go mod tidy
rsrc -manifest app.manifest -arch amd64 -o rsrc_windows_amd64.syso
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w -H windowsgui" -o 上班鼓励器.exe
```

## 技术栈

- Go + [lxn/walk](https://github.com/lxn/walk)（纯 Win32 GUI，无依赖）
- rsrc 嵌入 manifest（声明 Win7 支持 + 隐藏控制台）

## License

MIT
