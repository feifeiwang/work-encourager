# 上班鼓励器 💰

一个用 Go 写的 Windows 桌面小工具：把你的月薪折算成**每秒进账**，看着钱一点点挣回来，上班更有盼头。

复刻自「上班鼓励器」应用，**兼容 Windows 7 / 8 / 10 / 11**，单文件 exe，双击即开，无需安装任何运行库。

## 界面预览

```
┌─────────────────────────────────┐
│ 🔥 上班鼓励器    [📌 置顶] [⚙ 设置] │
│ 2026-11-21 周五  09:32:15        │
│ 上班中                           │
│ 距离下班还有 8小时27分             │
│                                 │
│ 今日已赚                         │
│ ¥27.9104                        │
│ 当前赚钱速度 ¥0.1543 / 秒          │
│ 再过6小时53分，一台Switch就挣回来了。│
│                                 │
│ [ ⏸ 关闭计时器 ]                 │
│ 今日工时进度 20.4%               │
│ 已工作 1.63 / 8.00 小时  ▓▓░░░░░ │
│ ──────────────────────────────  │
│ 本月累计             ¥1258.7315  │
│ 时薪 ¥24.68      今日满勤可得 ¥197.48│
└─────────────────────────────────┘
```

## 功能

| 功能 | 说明 |
|---|---|
| 🕐 实时时钟 | 日期、星期、秒级时钟，下班倒计时 |
| 💰 实时挣钱 | 今日已赚按秒跳动，显示赚钱速度（元/秒），按 `月薪 ÷ 21.75 天 ÷ 每日工时` 折算时薪 |
| 🗣 随机鼓励语 | 内置 **57 条**商品语录（一瓶冰可乐 → 一部 iPhone → 一个金镯子），每 12 秒随机换一条，按你的赚钱速度实时算出「再过 X 小时 X 分，XX 就挣回来了」 |
| 📌 窗口置顶 | 一键置顶 / 取消置顶（原生 Win32 `SetWindowPos`），挂在屏幕角落不挡事 |
| ⏸ 暂停计时 | 关闭/开启计时器，暂停时时间冻结，恢复后接着算，不丢失进度 |
| 📊 工时进度 | 今日工时进度条 + 已工作小时数 |
| 📅 本月累计 | 已过工作日按满勤补算（周末不计薪） |
| ⚙ 自定义设置 | 月薪、每天工作时长、上班/下班时间，保存到 exe 同目录 `config.json` |

## 下载使用

到 [Releases](../../releases/latest) 下载对应版本：

| 文件 | 适用系统 |
|---|---|
| `work-encourager-win64.exe` | 64 位 Windows |
| `work-encourager-win32.exe` | 32 位 Windows |

> - 单文件、免安装，双击即开
> - 请放在**有写权限的目录**（如桌面、D 盘），否则设置无法保存
> - 删除 `config.json` 即恢复默认值

### 默认配置

| 配置项 | 默认值 |
|---|---|
| 月薪 | ¥4000 |
| 每天工作时长 | 8 小时 |
| 上班时间 | 09:00 |
| 下班时间 | 18:00 |

保存设置后，exe 同目录会生成 `config.json`：

```json
{
  "月薪": 5000,
  "每日工作时长": 8,
  "上班时间": "09:00",
  "下班时间": "18:00"
}
```

## 从源码构建

⚠️ 必须使用 **Go 1.20 及以下版本**编译（Go 1.21+ 生成的 exe 无法在 Win7 运行）。

```bash
# 1. 拉取依赖
go mod tidy

# 2. 生成 manifest 资源（需先安装 rsrc：go install github.com/akavel/rsrc@latest）
rsrc -manifest app.manifest -arch amd64 -o rsrc_windows_amd64.syso
rsrc -manifest app.manifest -arch 386  -o rsrc_windows_386.syso

# 3. 交叉编译
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w -H windowsgui" -o 上班鼓励器.exe
GOOS=windows GOARCH=386  go build -trimpath -ldflags "-s -w -H windowsgui" -o 上班鼓励器_32位.exe
```

`-ldflags "-H windowsgui"` 用于隐藏控制台窗口；syso 文件同时嵌入 Win7 兼容声明和现代视觉样式。

## 技术栈

- Go 1.20 + [lxn/walk](https://github.com/lxn/walk)（纯 Win32 GUI，无运行库依赖）
- [akavel/rsrc](https://github.com/akavel/rsrc) 嵌入 manifest（Win7 支持声明 + 隐藏控制台）
- 配置持久化：JSON（`config.json`）

## 常见问题

**Q：双击没反应？**
A：检查是否被杀毒软件拦截（单文件 exe 常见误报），添加信任即可。

**Q：设置保存不了？**
A：exe 所在目录无写权限（如 `C:\Program Files`），移到桌面或其他普通目录。

**Q：Win7 运行报错「不是有效的 Win32 应用程序」？**
A：你的系统是 32 位，请下载 `work-encourager-win32.exe`。

## License

[MIT](./LICENSE)
