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

// ---------- 商品语录清单 ----------

var items = []struct {
	name  string
	price float64
}{
	{"一瓶冰可乐", 3.5}, {"一支雪糕", 5}, {"一杯柠檬水", 6}, {"一份豆浆油条", 7},
	{"一个煎饼果子", 8}, {"一份肠粉", 10}, {"一碗牛肉面", 15}, {"一份螺蛳粉", 16},
	{"一杯瑞幸咖啡", 15}, {"一顿沙县小吃", 15}, {"一杯奶茶", 18}, {"一份饺子", 20},
	{"一顿黄焖鸡米饭", 22}, {"一份麦当劳套餐", 35}, {"一杯星巴克", 35}, {"一顿外卖", 25},
	{"一次电影票", 45}, {"一桶爆米花", 30}, {"一箱牛奶", 65}, {"一箱水果", 88},
	{"一顿火锅", 130}, {"一顿烤肉自助", 98}, {"一顿小龙虾", 150}, {"一顿日料", 180},
	{"一场剧本杀", 88}, {"一场KTV", 100}, {"一张演唱会门票", 680}, {"一张球赛门票", 300},
	{"一件T恤", 59}, {"一件衬衫", 99}, {"一条牛仔裤", 149}, {"一件卫衣", 199},
	{"一双运动鞋", 400}, {"一支口红", 220}, {"一瓶香水", 350}, {"一套护肤品", 500},
	{"一个充电宝", 99}, {"一副蓝牙耳机", 299}, {"一个机械键盘", 350}, {"一支电动牙刷", 199},
	{"一个吹风机", 299}, {"一个篮球", 129}, {"一张健身房月卡", 200}, {"一盒猫粮", 120},
	{"一个月话费", 100}, {"一个月房租", 1800}, {"一张回家的机票", 800}, {"一部新手机", 5999},
	{"一台平板电脑", 2999}, {"一台Switch", 2100}, {"一台PS5", 3899}, {"一台笔记本电脑", 7999},
	{"一台空调", 2999}, {"一台洗衣机", 1500}, {"一辆电动车", 2500}, {"一个金镯子", 18000},
	{"一辆五菱宏光", 45000}, {"一个车位", 80000},
}

// ---------- 配色 ----------

var (
	// 清新薄荷配色
	cMint = walk.RGB(0x2e, 0xa7, 0x7d) // 主色：薄荷绿
	cMoney = walk.RGB(0x1f, 0xa0, 0x74) // 金额：清爽绿
	cSub   = walk.RGB(0x7a, 0x94, 0x88) // 次要文字：灰绿
	cBody  = walk.RGB(0x44, 0x55, 0x4e) // 正文：墨绿灰
	cOff   = walk.RGB(0xe0, 0x9a, 0x50) // 下班：柔和橙
	cIdle  = walk.RGB(0x9a, 0xaa, 0xa2) // 未上班：浅灰
	cBg    = walk.RGB(0xf2, 0xfa, 0xf7) // 背景：极浅薄荷
)

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
	statusL  *walk.Label
	cdL      *walk.Label
	moneyL   *walk.Label
	speedL   *walk.Label
	phraseL  *walk.Label
	progL    *walk.Label
	prog     *walk.ProgressBar
	monthL   *walk.Label
	rateL    *walk.Label
	fullL    *walk.Label
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
	dateL.SetText(vNow.Format("2006-01-02 ") + wd + "  " + vNow.Format("15:04:05"))

	nMin := vNow.Hour()*60 + vNow.Minute()
	var worked float64 // 0~1
	var status string
	switch {
	case nMin < start:
		worked, status = 0, "还没上班"
	case nMin >= end:
		worked, status = 1, "下班啦 🎉"
	default:
		worked = float64(nMin-start) / float64(end-start)
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
		cdL.SetText(fmt.Sprintf("距离下班还有 %s", fmtDur(float64((end-nMin))*60)))
	} else if status == "下班啦 🎉" {
		cdL.SetText("今天的班上完了，好好休息！")
	} else {
		remain := float64(end*60-(nMin*60+vNow.Second()))
		cdL.SetText(fmt.Sprintf("距离下班还有 %s", fmtDur(remain)))
	}

	moneyL.SetText(fmt.Sprintf("¥%.4f", earned))
	speedL.SetText(fmt.Sprintf("当前赚钱速度 ¥%.4f / 秒", speed))

	it := items[itemIdx]
	if remainMoney := it.price - earned; remainMoney <= 0 {
		phraseL.SetText(fmt.Sprintf("今天已经挣回%s，还有富余！", it.name))
	} else {
		phraseL.SetText(fmt.Sprintf("再过%s，%s就挣回来了。", fmtDur(remainMoney/speed), it.name))
	}

	progL.SetText(fmt.Sprintf("今日工时进度 %.1f%%\n已工作 %.2f / %.2f 小时", worked*100, workedH, cfg.DailyHours))
	prog.SetValue(int(worked * 100))

	monthL.SetText(fmt.Sprintf("¥%.4f", monthEarn))
	rateL.SetText(fmt.Sprintf("时薪 ¥%.2f", hourly))
	fullL.SetText(fmt.Sprintf("今日满勤可得 ¥%.2f", dailyEarn))

	if status == "上班中" && running {
		statusL.SetTextColor(cMint)
	} else if status == "下班啦 🎉" {
		statusL.SetTextColor(cOff)
	} else {
		statusL.SetTextColor(cIdle)
	}
	statusL.SetText(status)
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
	vNow, lastReal = time.Now(), time.Now()
	itemIdx = rnd.Intn(len(items))

	MainWindow{
		Title:      "上班鼓励器",
		AssignTo:   &mw,
		MinSize:    Size{Width: 250, Height: 400},
		MaxSize:    Size{Width: 250, Height: 400},
		Background: SolidColorBrush{Color: cBg},
		Font:       Font{Family: "Microsoft YaHei", PointSize: 8},
		Layout:     VBox{Margins: Margins{Left: 10, Top: 6, Right: 10, Bottom: 10}, Spacing: 5},
		Children: []Widget{
			Composite{
				Layout: HBox{Spacing: 4},
				Children: []Widget{
					Label{Text: "🔥 上班鼓励器", Font: Font{PointSize: 10, Bold: true}, TextColor: cMint},
					HSpacer{},
					PushButton{AssignTo: &pinBtn, Text: "📌 置顶", OnClicked: togglePin},
					PushButton{Text: "⚙ 设置", OnClicked: openSettings},
				},
			},
			Label{AssignTo: &dateL, TextColor: cSub},
			Label{AssignTo: &statusL, Text: "上班中", Font: Font{PointSize: 12, Bold: true}},
			Label{AssignTo: &cdL, Text: "距离下班还有 --", TextColor: cSub},
			VSpacer{Size: 2},
			Label{Text: "今日已赚", TextColor: cSub},
			Label{AssignTo: &moneyL, Text: "¥0.0000", Font: Font{PointSize: 15, Bold: true}, TextColor: cMoney},
			Label{AssignTo: &speedL, TextColor: cSub},
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
			Label{Text: "已过的日子按满勤补算。本月只计已赚，不含亏损。", Font: Font{PointSize: 6}, TextColor: cIdle},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{AssignTo: &rateL, TextColor: cSub},
					HSpacer{},
					Label{AssignTo: &fullL, TextColor: cSub},
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
