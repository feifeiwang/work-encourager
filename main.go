package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

// ---------- 配置 ----------

type Config struct {
	MonthlySalary float64 `json:"月薪"`     // 每月到手收入
	DailyHours    float64 `json:"每日工作时长"` // 每天实际工作小时数
	StartTime     string  `json:"上班时间"`   // HH:MM
	EndTime       string  `json:"下班时间"`   // HH:MM
	Theme         int     `json:"主题"`     // 配色编号
}

func defaultConfig() Config {
	return Config{MonthlySalary: 4000, DailyHours: 8, StartTime: "09:00", EndTime: "18:00"}
}

func configPath() string {
	if exe, err := os.Executable(); err == nil {
		return filepath.Join(filepath.Dir(exe), "config.json")
	}
	return "config.json"
}

func loadConfig() Config {
	c := defaultConfig()
	if b, err := os.ReadFile(configPath()); err == nil {
		json.Unmarshal(b, &c)
	}
	return c
}

func saveConfig(c Config) {
	b, _ := json.MarshalIndent(c, "", "  ")
	os.WriteFile(configPath(), b, 0644)
}

// ---------- 商品语录清单（便宜好物，一天工资买得起） ----------

var items = []struct {
	name  string
	unit  string
	price float64
}{
	{"豆浆", "杯", 3}, {"茶叶蛋", "个", 2.5}, {"冰可乐", "罐", 3.5}, {"雪糕", "支", 5},
	{"辣条", "包", 5}, {"柠檬水", "杯", 6}, {"红牛", "罐", 6}, {"烤串", "串", 5},
	{"煎饼果子", "份", 8}, {"肠粉", "份", 10}, {"汉堡", "个", 12}, {"炒饭", "份", 12},
	{"毛巾", "条", 12}, {"牛肉面", "碗", 15}, {"瑞幸咖啡", "杯", 15}, {"沙县小吃", "顿", 15},
	{"螺蛳粉", "份", 16}, {"奶茶", "杯", 18}, {"饺子", "份", 20}, {"网吧包夜", "次", 20},
	{"黄焖鸡米饭", "份", 22}, {"麻辣烫", "顿", 25}, {"外卖", "顿", 25}, {"炸鸡", "份", 25},
	{"视频会员", "个月", 25}, {"爆米花", "桶", 30}, {"西瓜", "个", 30}, {"一本小说", "本", 30},
	{"一束花", "束", 30}, {"多肉植物", "盆", 15}, {"抽纸", "提", 15}, {"鸡蛋", "盒", 15},
	{"麦当劳套餐", "份", 35}, {"星巴克", "杯", 35}, {"洗发水", "瓶", 40}, {"车厘子", "斤", 40},
	{"电影票", "张", 45}, {"抱枕", "个", 45}, {"一箱啤酒", "箱", 50}, {"T恤", "件", 59},
	{"牛奶", "箱", 65}, {"水果", "箱", 88}, {"衬衫", "件", 99}, {"充电宝", "个", 99},
}

// ---------- 配色（换肤） ----------

type Theme struct {
	name                 string
	bg, main, money, sub, body, off, idle walk.Color
}

var themes = []Theme{
	{name: "薄荷绿", bg: walk.RGB(0xf2, 0xfa, 0xf7), main: walk.RGB(0x2e, 0xa7, 0x7d), money: walk.RGB(0x1f, 0xa0, 0x74), sub: walk.RGB(0x7a, 0x94, 0x88), body: walk.RGB(0x44, 0x55, 0x4e), off: walk.RGB(0xe0, 0x9a, 0x50), idle: walk.RGB(0x9a, 0xaa, 0xa2)},
	{name: "奶油黄", bg: walk.RGB(0xff, 0xf9, 0xee), main: walk.RGB(0xd9, 0xa4, 0x41), money: walk.RGB(0xc9, 0x8a, 0x2d), sub: walk.RGB(0x9c, 0x8b, 0x6e), body: walk.RGB(0x5c, 0x53, 0x46), off: walk.RGB(0xe0, 0x7b, 0x54), idle: walk.RGB(0xb5, 0xa8, 0x8e)},
	{name: "樱花粉", bg: walk.RGB(0xfd, 0xf3, 0xf5), main: walk.RGB(0xd9, 0x7b, 0x8f), money: walk.RGB(0xc9, 0x68, 0x77), sub: walk.RGB(0xa2, 0x84, 0x8c), body: walk.RGB(0x5c, 0x4a, 0x50), off: walk.RGB(0xe0, 0xa0, 0x50), idle: walk.RGB(0xb7, 0x9c, 0xa5)},
	{name: "雾霾蓝", bg: walk.RGB(0xf0, 0xf5, 0xfa), main: walk.RGB(0x5b, 0x8d, 0xb8), money: walk.RGB(0x41, 0x79, 0xa8), sub: walk.RGB(0x7e, 0x93, 0xa6), body: walk.RGB(0x42, 0x50, 0x5c), off: walk.RGB(0xc9, 0x8a, 0x5a), idle: walk.RGB(0x9a, 0xa8, 0xb5)},
	{name: "紫罗兰", bg: walk.RGB(0xf6, 0xf4, 0xfa), main: walk.RGB(0x8b, 0x7f, 0xc7), money: walk.RGB(0x74, 0x68, 0xb5), sub: walk.RGB(0x8f, 0x8a, 0xa6), body: walk.RGB(0x4c, 0x48, 0x59), off: walk.RGB(0xc9, 0x8a, 0x5a), idle: walk.RGB(0xa8, 0xa3, 0xbc)},
}

// 当前生效的配色（由 setTheme 切换）
var (
	cMint, cMoney, cSub, cBody, cOff, cIdle, cBg walk.Color
	bgBrush *walk.SolidColorBrush
)

func setTheme(i int) {
	cfg.Theme = i % len(themes)
	t := themes[cfg.Theme]
	cBg, cMint, cMoney, cSub, cBody, cOff, cIdle = t.bg, t.main, t.money, t.sub, t.body, t.off, t.idle
	if mw != nil {
		if bgBrush != nil {
			bgBrush.Dispose()
		}
		b, err := walk.NewSolidColorBrush(cBg)
		if err != nil {
			return
		}
		bgBrush = b
		mw.SetBackground(bgBrush)
		update()
	}
}

// ---------- 全局状态 ----------

var (
	cfg    Config
	running = true
	topmost bool
	vNow    time.Time // 虚拟当前时间（暂停时冻结）
	lastReal time.Time
	itemIdx int
	rnd     = rand.New(rand.NewSource(time.Now().UnixNano()))

	mw       *walk.MainWindow
	dateL    *walk.Label
	clockL   *walk.Label
	statusL  *walk.Label
	cdL      *walk.Label
	moneyL   *walk.Label
	speedL   *walk.Label
	phraseL  *walk.Label
	progL    *walk.Label
	prog     *walk.ProgressBar
	monthL   *walk.Label
	titleL   *walk.Label
	capL     *walk.Label
	skinBtn  *walk.PushButton
	toggleBtn *walk.PushButton
	pinBtn   *walk.PushButton
)

// ---------- 工具函数 ----------

func parseHM(s string) int {
	p := strings.Split(strings.TrimSpace(s), ":")
	if len(p) != 2 {
		return -1
	}
	h, e1 := strconv.Atoi(strings.TrimSpace(p[0]))
	m, e2 := strconv.Atoi(strings.TrimSpace(p[1]))
	if e1 != nil || e2 != nil || h < 0 || h > 23 || m < 0 || m > 59 {
		return -1
	}
	return h*60 + m
}

func fmtDur(sec float64) string {
	m := int(sec / 60)
	if m < 1 {
		return "不到 1 分钟"
	}
	if m < 60 {
		return fmt.Sprintf("%d 分钟", m)
	}
	if m < 60*48 {
		return fmt.Sprintf("%d 小时 %d 分", m/60, m%60)
	}
	return fmt.Sprintf("%d 天 %d 小时", m/1440, (m%1440)/60)
}

// ---------- 界面刷新 ----------

func update() {
	hourly := cfg.MonthlySalary / (21.75 * cfg.DailyHours)
	if hourly <= 0 {
		hourly = 1
	}
	speed := hourly / 3600.0
	dailyEarn := hourly * cfg.DailyHours

	start, end := parseHM(cfg.StartTime), parseHM(cfg.EndTime)
	if start < 0 || end < 0 || end <= start {
		start, end = 540, 1080
	}

	wd := [...]string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}[int(vNow.Weekday())]
	dateL.SetText(vNow.Format("2006-01-02 ") + wd)
	clockL.SetText(vNow.Format("15:04:05"))

	nSec := vNow.Hour()*3600 + vNow.Minute()*60 + vNow.Second() // 秒级精度，金额才能按秒跳动
	startSec, endSec := start*60, end*60
	var worked float64 // 0~1
	var status string
	switch {
	case nSec < startSec:
		worked, status = 0, "还没上班"
	case nSec >= endSec:
		worked, status = 1, "下班啦 🎉"
	default:
		worked = float64(nSec-startSec) / float64(endSec-startSec)
		status = "上班中"
	}

	workedH := worked * cfg.DailyHours
	earned := workedH * hourly
	monthEarn := earned
	for d := time.Date(vNow.Year(), vNow.Month(), 1, 0, 0, 0, 0, vNow.Location()); d.Before(vNow); d = d.AddDate(0, 0, 1) {
		if w := d.Weekday(); w != time.Saturday && w != time.Sunday {
			monthEarn += dailyEarn
		}
	}

	if status == "还没上班" {
		cdL.SetText("距上班还有 " + fmtDur(float64(startSec-nSec)))
	} else if status == "下班啦 🎉" {
		cdL.SetText("今天下班啦")
	} else {
		cdL.SetText("距下班 " + fmtDur(float64(endSec-nSec)))
	}

	moneyL.SetText(fmt.Sprintf("¥%.4f", earned))
	speedL.SetText(fmt.Sprintf("当前赚钱速度 ¥%.4f / 秒", speed))

	it := items[itemIdx]
	if n := int(earned / it.price); n >= 1 {
		phraseL.SetText(fmt.Sprintf("今天的价钱能买 %d %s%s啦！", n, it.unit, it.name))
	} else {
		phraseL.SetText(fmt.Sprintf("再过%s，今天的价钱就能买 1 %s%s了。",
			fmtDur((it.price-earned)/speed), it.unit, it.name))
	}

	progL.SetText(fmt.Sprintf("今日工时进度 %.1f%%\n已工作 %.2f / %.2f 小时", worked*100, workedH, cfg.DailyHours))
	prog.SetValue(int(worked * 100))

	monthL.SetText(fmt.Sprintf("¥%.4f", monthEarn))

	if status == "上班中" && running {
		statusL.SetTextColor(cMint)
	} else if status == "下班啦 🎉" {
		statusL.SetTextColor(cOff)
	} else {
		statusL.SetTextColor(cIdle)
	}
	statusL.SetText(status)
	titleL.SetTextColor(cMint)
	capL.SetTextColor(cSub)
	skinBtn.SetText("🎨 " + themes[cfg.Theme].name)
}

// ---------- 置顶 ----------

func togglePin() {
	topmost = !topmost
	flag := uint32(win.SWP_NOMOVE | win.SWP_NOSIZE | win.SWP_NOACTIVATE)
	hwnd := win.HWND_TOPMOST
	if !topmost {
		hwnd = win.HWND_NOTOPMOST
	}
	win.SetWindowPos(mw.Handle(), hwnd, 0, 0, 0, 0, flag)
	if topmost {
		pinBtn.SetText("取消置顶")
	} else {
		pinBtn.SetText("📌 置顶")
	}
}

// ---------- 设置对话框 ----------

func openSettings() {
	var dlg *walk.Dialog
	var okPB *walk.PushButton
	var neSalary, neHours *walk.NumberEdit
	var leStart, leEnd *walk.LineEdit

	Dialog{
		Title:    "设置",
		AssignTo: &dlg,
		MinSize:  Size{Width: 280, Height: 200},
		Layout:   VBox{},
		DefaultButton: &okPB,
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2, Spacing: 10},
				Children: []Widget{
					Label{Text: "月薪（元）"},
					NumberEdit{AssignTo: &neSalary, Value: cfg.MonthlySalary, MinValue: 1, MaxValue: 10000000, Decimals: 2},
					Label{Text: "每天工作时长（小时）"},
					NumberEdit{AssignTo: &neHours, Value: cfg.DailyHours, MinValue: 1, MaxValue: 24, Decimals: 2},
					Label{Text: "上班时间（HH:MM）"},
					LineEdit{AssignTo: &leStart, Text: cfg.StartTime, CueBanner: "09:00"},
					Label{Text: "下班时间（HH:MM）"},
					LineEdit{AssignTo: &leEnd, Text: cfg.EndTime, CueBanner: "18:00"},
				},
			},
			VSpacer{},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{AssignTo: &okPB, Text: "保存", OnClicked: func() {
						s, e := parseHM(leStart.Text()), parseHM(leEnd.Text())
						if s < 0 || e < 0 || e <= s {
							walk.MsgBox(dlg, "提示", "时间格式应为 HH:MM，且下班时间要晚于上班时间", walk.MsgBoxIconWarning)
							return
						}
						cfg.MonthlySalary = neSalary.Value()
						cfg.DailyHours = neHours.Value()
						cfg.StartTime = strings.TrimSpace(leStart.Text())
						cfg.EndTime = strings.TrimSpace(leEnd.Text())
						saveConfig(cfg)
						dlg.Accept()
					}},
					PushButton{Text: "取消", OnClicked: func() { dlg.Cancel() }},
				},
			},
		},
	}.Create(mw)
	dlg.Run()
}

// ---------- 主程序 ----------

func main() {
	cfg = loadConfig()
	setTheme(cfg.Theme)
	vNow, lastReal = time.Now(), time.Now()
	itemIdx = rnd.Intn(len(items))

	MainWindow{
		Title:      "上班鼓励器",
		AssignTo:   &mw,
		MinSize:    Size{Width: 245, Height: 330},
		MaxSize:    Size{Width: 245, Height: 330},
		Size:       Size{Width: 245, Height: 330}, // 打开即为此大小，不要用户手动调整
		Background: SolidColorBrush{Color: cBg},
		Font:       Font{Family: "Microsoft YaHei", PointSize: 8},
		Layout:     VBox{Margins: Margins{Left: 8, Top: 6, Right: 8, Bottom: 8}, Spacing: 5},
		Children: []Widget{
			Composite{
				Layout: HBox{Spacing: 4},
				Children: []Widget{
					Label{AssignTo: &titleL, Text: "上班鼓励器", Font: Font{PointSize: 9, Bold: true}, TextColor: cMint},
					HSpacer{},
					PushButton{AssignTo: &skinBtn, Text: "🎨 换肤", OnClicked: func() {
						setTheme(cfg.Theme + 1)
						saveConfig(cfg)
					}},
					PushButton{AssignTo: &pinBtn, Text: "📌 置顶", OnClicked: togglePin},
					PushButton{Text: "⚙", OnClicked: openSettings},
				},
			},
			Composite{
				// 左栏（贴左）：日期/时间/状态/倒计时；右栏（贴右）：今日已赚/金额/速度
				Layout: HBox{Margins: Margins{}},
				Children: []Widget{
					Composite{
						Layout: VBox{Spacing: 3, Margins: Margins{}},
						Children: []Widget{
							Label{AssignTo: &dateL, TextColor: cSub, Font: Font{PointSize: 8}},
							Label{AssignTo: &clockL, Text: "15:04:05", Font: Font{PointSize: 13, Bold: true}, TextColor: cBody},
							Label{AssignTo: &statusL, Text: "上班中", Font: Font{PointSize: 12, Bold: true}},
							Label{AssignTo: &cdL, Text: "距下班 --", TextColor: cSub, Font: Font{PointSize: 8}},
						},
					},
					HSpacer{},
					Composite{
						Layout: VBox{Spacing: 3, Margins: Margins{}},
						Children: []Widget{
							Label{AssignTo: &capL, Text: "今日已赚", TextColor: cSub, Font: Font{PointSize: 8}},
							Label{AssignTo: &moneyL, Text: "¥0.0000", Font: Font{PointSize: 16, Bold: true}, TextColor: cMoney},
							VSpacer{},
							Label{AssignTo: &speedL, TextColor: cSub, Font: Font{PointSize: 8}},
						},
					},
				},
			},
			Label{AssignTo: &phraseL, TextColor: cBody},
			VSpacer{Size: 2},
			PushButton{AssignTo: &toggleBtn, Text: "⏸ 关闭计时器", MinSize: Size{Width: 0, Height: 28}, OnClicked: func() {
				running = !running
				if running {
					toggleBtn.SetText("⏸ 关闭计时器")
				} else {
					toggleBtn.SetText("▶ 开启计时器")
				}
			}},
			VSpacer{Size: 2},
			Label{AssignTo: &progL, TextColor: cSub},
			ProgressBar{AssignTo: &prog, MinValue: 0, MaxValue: 100},
			HSeparator{},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{Text: "本月累计", TextColor: cSub},
					HSpacer{},
					Label{AssignTo: &monthL, Font: Font{PointSize: 10, Bold: true}, TextColor: cMoney},
				},
			},
		},
	}.Create()

	go func() {
		tick := 0
		for range time.Tick(time.Second) {
			now := time.Now()
			if running {
				vNow = vNow.Add(now.Sub(lastReal))
			}
			lastReal = now
			tick++
			if tick%12 == 0 {
				itemIdx = rnd.Intn(len(items))
			}
			mw.Synchronize(update)
		}
	}()

	mw.Run()
}
