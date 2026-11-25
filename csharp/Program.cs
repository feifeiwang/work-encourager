using System;
using System.Drawing;
using System.Windows.Forms;
using System.IO;
using System.Text.RegularExpressions;
using System.Runtime.InteropServices;
using System.Globalization;

[assembly: System.Runtime.Versioning.TargetFramework(".NETFramework,Version=v4.0")]

namespace WorkEncourager
{
    // ---------- 配置 ----------
    public class Config
    {
        public double Salary = 4000;      // 月薪
        public double Hours = 8;          // 每日工作时长
        public string Start = "09:00";
        public string End = "18:00";
        public int Theme = 0;
    }

    // ---------- 配色主题 ----------
    public struct Theme
    {
        public string Name;
        public Color Bg, Main, Money, Sub, Body, Off, Idle;
        public Theme(string n, Color bg, Color main, Color money, Color sub, Color body, Color off, Color idle)
        {
            Name = n; Bg = bg; Main = main; Money = money; Sub = sub; Body = body; Off = off; Idle = idle;
        }
    }

    // ---------- 便宜好物 ----------
    public struct Item
    {
        public string Name, Unit;
        public double Price;
        public Item(string n, string u, double p) { Name = n; Unit = u; Price = p; }
    }

    public static class Program
    {
        public static Config cfg = new Config();
        public static Theme[] Themes = new Theme[] {
            new Theme("薄荷绿", C(0xF2,0xFA,0xF7), C(0x2E,0xA7,0x7D), C(0x1F,0xA0,0x74), C(0x7A,0x94,0x88), C(0x44,0x55,0x4E), C(0xE0,0x9A,0x50), C(0x9A,0xAA,0xA2)),
            new Theme("奶油黄", C(0xFF,0xF9,0xEE), C(0xD9,0xA4,0x41), C(0xC9,0x8A,0x2D), C(0x9C,0x8B,0x6E), C(0x5C,0x53,0x46), C(0xE0,0x7B,0x54), C(0xB5,0xA8,0x8E)),
            new Theme("樱花粉", C(0xFD,0xF3,0xF5), C(0xD9,0x7B,0x8F), C(0xC9,0x68,0x77), C(0xA2,0x84,0x8C), C(0x5C,0x4A,0x50), C(0xE0,0xA0,0x50), C(0xB7,0x9C,0xA5)),
            new Theme("雾霾蓝", C(0xF0,0xF5,0xFA), C(0x5B,0x8D,0xB8), C(0x41,0x79,0xA8), C(0x7E,0x93,0xA6), C(0x42,0x50,0x5C), C(0xC9,0x8A,0x5A), C(0x9A,0xA8,0xB5)),
            new Theme("紫罗兰", C(0xF6,0xF4,0xFA), C(0x8B,0x7F,0xC7), C(0x74,0x68,0xB5), C(0x8F,0x8A,0xA6), C(0x4C,0x48,0x59), C(0xC9,0x8A,0x5A), C(0xA8,0xA3,0xBC)),
        };

        public static Item[] Items = new Item[] {
            new Item("豆浆","杯",3), new Item("茶叶蛋","个",2.5), new Item("冰可乐","罐",3.5), new Item("雪糕","支",5),
            new Item("辣条","包",5), new Item("柠檬水","杯",6), new Item("红牛","罐",6), new Item("烤串","串",5),
            new Item("煎饼果子","份",8), new Item("肠粉","份",10), new Item("汉堡","个",12), new Item("炒饭","份",12),
            new Item("毛巾","条",12), new Item("牛肉面","碗",15), new Item("瑞幸咖啡","杯",15), new Item("沙县小吃","顿",15),
            new Item("螺蛳粉","份",16), new Item("奶茶","杯",18), new Item("饺子","份",20), new Item("网吧包夜","次",20),
            new Item("黄焖鸡米饭","份",22), new Item("麻辣烫","顿",25), new Item("外卖","顿",25), new Item("炸鸡","份",25),
            new Item("视频会员","个月",25), new Item("爆米花","桶",30), new Item("西瓜","个",30), new Item("一本小说","本",30),
            new Item("一束花","束",30), new Item("多肉植物","盆",15), new Item("抽纸","提",15), new Item("鸡蛋","盒",15),
            new Item("麦当劳套餐","份",35), new Item("星巴克","杯",35), new Item("洗发水","瓶",40), new Item("车厘子","斤",40),
            new Item("电影票","张",45), new Item("抱枕","个",45), new Item("一箱啤酒","箱",50), new Item("T恤","件",59),
            new Item("牛奶","箱",65), new Item("水果","箱",88), new Item("衬衫","件",99), new Item("充电宝","个",99),
        };

        static Color C(int r, int g, int b) { return Color.FromArgb(r, g, b); }

        public static string ConfigPath()
        {
            return Path.Combine(Path.GetDirectoryName(Application.ExecutablePath), "config.json");
        }

        public static Config LoadConfig()
        {
            Config c = new Config();
            try
            {
                string s = File.ReadAllText(ConfigPath());
                Match m;
                m = Regex.Match(s, "\"月薪\"\\s*:\\s*([0-9.]+)");
                if (m.Success) c.Salary = double.Parse(m.Groups[1].Value, CultureInfo.InvariantCulture);
                m = Regex.Match(s, "\"每日工作时长\"\\s*:\\s*([0-9.]+)");
                if (m.Success) c.Hours = double.Parse(m.Groups[1].Value, CultureInfo.InvariantCulture);
                m = Regex.Match(s, "\"上班时间\"\\s*:\\s*\"([^\"]+)\"");
                if (m.Success) c.Start = m.Groups[1].Value;
                m = Regex.Match(s, "\"下班时间\"\\s*:\\s*\"([^\"]+)\"");
                if (m.Success) c.End = m.Groups[1].Value;
                m = Regex.Match(s, "\"主题\"\\s*:\\s*(\\d+)");
                if (m.Success) c.Theme = int.Parse(m.Groups[1].Value);
            }
            catch { }
            if (c.Hours <= 0) c.Hours = 8;
            if (c.Salary <= 0) c.Salary = 4000;
            if (ParseHM(c.Start) < 0) c.Start = "09:00";
            if (ParseHM(c.End) < 0) c.End = "18:00";
            if (c.Theme < 0 || c.Theme >= Themes.Length) c.Theme = 0;
            return c;
        }

        public static void SaveConfig(Config c)
        {
            try
            {
                string s = "{\n  \"月薪\": " + c.Salary.ToString("0.00", CultureInfo.InvariantCulture)
                    + ",\n  \"每日工作时长\": " + c.Hours.ToString("0.00", CultureInfo.InvariantCulture)
                    + ",\n  \"上班时间\": \"" + c.Start + "\""
                    + ",\n  \"下班时间\": \"" + c.End + "\""
                    + ",\n  \"主题\": " + c.Theme + "\n}";
                File.WriteAllText(ConfigPath(), s);
            }
            catch { }
        }

        public static int ParseHM(string s)
        {
            if (s == null) return -1;
            string[] p = s.Split(':');
            if (p.Length != 2) return -1;
            int h, m;
            if (!int.TryParse(p[0].Trim(), out h) || !int.TryParse(p[1].Trim(), out m)) return -1;
            if (h < 0 || h > 23 || m < 0 || m > 59) return -1;
            return h * 60 + m;
        }

        public static string FmtDur(double sec)
        {
            int m = (int)(sec / 60);
            if (m < 1) return "不到 1 分钟";
            if (m < 60) return m + " 分钟";
            if (m < 60 * 48) return (m / 60) + " 小时 " + (m % 60) + " 分";
            return (m / 1440) + " 天 " + ((m % 1440) / 60) + " 小时";
        }

        [STAThread]
        static void Main()
        {
            Application.EnableVisualStyles();
            Application.SetCompatibleTextRenderingDefault(false);
            cfg = LoadConfig();
            Application.Run(new MainForm());
        }
    }

    public class MainForm : Form
    {
        [DllImport("user32.dll")]
        static extern bool SetWindowPos(IntPtr hWnd, IntPtr hWndInsertAfter, int X, int Y, int cx, int cy, uint uFlags);

        static readonly IntPtr HWND_TOPMOST = new IntPtr(-1);
        static readonly IntPtr HWND_NOTOPMOST = new IntPtr(-2);
        const uint SWP_NOSIZE = 0x0001, SWP_NOMOVE = 0x0002, SWP_NOACTIVATE = 0x0010;

        Label titleL, dateL, clockL, statusL, cdL, capL, moneyL, speedL, phraseL, progL, monthL;
        Button skinBtn, pinBtn, setBtn, toggleBtn;
        ProgressBar bar;
        Font moneyFont14, moneyFont12, moneyFont10, moneyFont8;
        int moneyFontSize = 14;
        Timer timer;
        DateTime vNow, lastReal;
        bool running = true, topmost = false;
        int itemIdx, tick;
        Random rnd = new Random();

        public MainForm()
        {
            Text = "上班鼓励器";
            FormBorderStyle = FormBorderStyle.FixedSingle;
            MaximizeBox = false;
            StartPosition = FormStartPosition.CenterScreen;
            // 固定像素布局：不让系统字体/DPI 缩放控件尺寸，避免别人电脑上字体偏大时布局错乱
            AutoScaleMode = AutoScaleMode.None;
            ClientSize = new Size(270, 335);
            try
            {
                Stream s = typeof(Program).Assembly.GetManifestResourceStream("app.ico");
                if (s != null) Icon = new Icon(s);
            }
            catch { }
            ApplyTheme(Program.cfg.Theme);

            // 整体字号比上一版再降一档，留足余量防止字体偏大的电脑显示不全
            Font f7 = new Font("Microsoft YaHei", 7f);
            Font f8 = new Font("Microsoft YaHei", 8f);
            Font f9 = new Font("Microsoft YaHei", 9f, FontStyle.Bold);
            Font f11 = new Font("Microsoft YaHei", 11f, FontStyle.Bold);
            Font f12 = new Font("Microsoft YaHei", 12f, FontStyle.Bold);
            moneyFont14 = new Font("Microsoft YaHei", 14f, FontStyle.Bold);
            moneyFont12 = new Font("Microsoft YaHei", 12f, FontStyle.Bold);
            moneyFont10 = new Font("Microsoft YaHei", 10f, FontStyle.Bold);
            moneyFont8 = new Font("Microsoft YaHei", 8f, FontStyle.Bold);

            titleL = L("上班鼓励器", 8, 10, 82, 20, f9, Program.Themes[Program.cfg.Theme].Main);
            skinBtn = B("换肤", 92, 8, 80, 26, f8); skinBtn.Click += delegate { ApplyTheme(Program.cfg.Theme + 1); Program.SaveConfig(Program.cfg); };
            pinBtn = B("置顶", 176, 8, 42, 26, f8); pinBtn.Click += delegate { TogglePin(); };
            setBtn = B("设置", 222, 8, 40, 26, f8); setBtn.Click += delegate { SettingsForm sf = new SettingsForm(); sf.ShowDialog(this); };

            // 左栏：日期 / 时间 / 状态 / 倒计时（贴左，宽度留余量）
            dateL = L("", 8, 44, 124, 15, f7, Program.Themes[Program.cfg.Theme].Sub);
            clockL = L("00:00:00", 8, 60, 124, 26, f12, Program.Themes[Program.cfg.Theme].Body);
            statusL = L("上班中", 8, 90, 124, 20, f11, Program.Themes[Program.cfg.Theme].Main);
            cdL = L("距下班 --", 8, 114, 124, 15, f7, Program.Themes[Program.cfg.Theme].Sub);

            // 右栏：今日已赚 / 金额 / 速度（贴右，中间留白）
            capL = L("今日已赚", 138, 44, 124, 15, f7, Program.Themes[Program.cfg.Theme].Sub);
            capL.TextAlign = ContentAlignment.TopRight;
            moneyL = L("¥0.0000", 130, 58, 132, 30, moneyFont14, Program.Themes[Program.cfg.Theme].Money);
            moneyL.TextAlign = ContentAlignment.TopRight;
            speedL = L("", 130, 114, 132, 15, f7, Program.Themes[Program.cfg.Theme].Sub);
            speedL.TextAlign = ContentAlignment.TopRight;

            phraseL = L("今天的价钱能买 1 杯奶茶啦！", 8, 138, 254, 34, f8, Program.Themes[Program.cfg.Theme].Body);

            toggleBtn = B("关闭计时器", 8, 178, 254, 28, f9);
            toggleBtn.Click += delegate
            {
                running = !running;
                toggleBtn.Text = running ? "关闭计时器" : "开启计时器";
            };

            progL = L("", 8, 212, 254, 28, f7, Program.Themes[Program.cfg.Theme].Sub);
            bar = new ProgressBar();
            bar.Location = new Point(8, 244); bar.Size = new Size(254, 14);
            bar.Minimum = 0; bar.Maximum = 100;
            Controls.Add(bar);

            Panel sep = new Panel();
            sep.Location = new Point(8, 268); sep.Size = new Size(254, 1);
            sep.BackColor = Program.Themes[Program.cfg.Theme].Sub;
            Controls.Add(sep);

            Label ml = L("本月累计", 8, 278, 80, 20, f8, Program.Themes[Program.cfg.Theme].Sub);
            monthL = L("", 130, 277, 132, 24, f9, Program.Themes[Program.cfg.Theme].Money);
            monthL.TextAlign = ContentAlignment.TopRight;

            vNow = DateTime.Now;
            lastReal = DateTime.Now;
            itemIdx = rnd.Next(Program.Items.Length);
            UpdateUI();

            timer = new Timer();
            timer.Interval = 1000;
            timer.Tick += delegate(object sender, EventArgs e)
            {
                DateTime now = DateTime.Now;
                if (running) vNow = vNow.Add(now.Subtract(lastReal));
                lastReal = now;
                tick++;
                if (tick % 12 == 0) itemIdx = rnd.Next(Program.Items.Length);
                UpdateUI();
            };
            timer.Start();
        }

        Label L(string text, int x, int y, int w, int h, Font f, Color c)
        {
            Label l = new Label();
            l.Text = text; l.Location = new Point(x, y); l.Size = new Size(w, h);
            // 关键：关闭自动缩放，否则别人电脑上系统字体偏大时 Label 会撑大、挤占相邻控件导致显示不全
            l.AutoSize = false;
            l.Font = f; l.ForeColor = c; l.BackColor = Color.Transparent;
            Controls.Add(l);
            return l;
        }

        Button B(string text, int x, int y, int w, int h, Font f)
        {
            Button b = new Button();
            b.Text = text; b.Location = new Point(x, y); b.Size = new Size(w, h);
            b.Font = f;
            Controls.Add(b);
            return b;
        }

        void ApplyTheme(int index)
        {
            Program.cfg.Theme = index % Program.Themes.Length;
            Theme t = Program.Themes[Program.cfg.Theme];
            BackColor = t.Bg;
            if (titleL == null) return;
            titleL.ForeColor = t.Main; dateL.ForeColor = t.Sub; clockL.ForeColor = t.Body;
            statusL.ForeColor = t.Main; cdL.ForeColor = t.Sub; capL.ForeColor = t.Sub;
            moneyL.ForeColor = t.Money; speedL.ForeColor = t.Sub; phraseL.ForeColor = t.Body;
            progL.ForeColor = t.Sub; monthL.ForeColor = t.Money;
            skinBtn.Text = "换肤(" + t.Name + ")";
        }

        void TogglePin()
        {
            topmost = !topmost;
            SetWindowPos(Handle, topmost ? HWND_TOPMOST : HWND_NOTOPMOST, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE | SWP_NOACTIVATE);
            pinBtn.Text = topmost ? "取消置顶" : "置顶";
        }

        void UpdateUI()
        {
            Config c = Program.cfg;
            Theme t = Program.Themes[c.Theme];
            double hourly = c.Salary / (21.75 * c.Hours);
            if (hourly <= 0) hourly = 1;
            double speed = hourly / 3600.0;
            double dailyEarn = hourly * c.Hours;

            int start = Program.ParseHM(c.Start), end = Program.ParseHM(c.End);
            if (start < 0 || end < 0 || end <= start) { start = 540; end = 1080; }

            string[] wd = { "周日", "周一", "周二", "周三", "周四", "周五", "周六" };
            dateL.Text = vNow.ToString("yyyy-MM-dd ") + wd[(int)vNow.DayOfWeek];
            clockL.Text = vNow.ToString("HH:mm:ss");

            int nSec = vNow.Hour * 3600 + vNow.Minute * 60 + vNow.Second;
            int startSec = start * 60, endSec = end * 60;
            double worked;
            string status;
            if (nSec < startSec) { worked = 0; status = "还没上班"; }
            else if (nSec >= endSec) { worked = 1; status = "下班啦"; }
            else { worked = (double)(nSec - startSec) / (double)(endSec - startSec); status = "上班中"; }

            double workedH = worked * c.Hours;
            double earned = workedH * hourly;
            double monthEarn = earned;
            DateTime first = new DateTime(vNow.Year, vNow.Month, 1);
            for (DateTime d = first; d < vNow; d = d.AddDays(1))
                if (d.DayOfWeek != DayOfWeek.Saturday && d.DayOfWeek != DayOfWeek.Sunday) monthEarn += dailyEarn;

            if (status == "还没上班") cdL.Text = "距上班还有 " + Program.FmtDur(startSec - nSec);
            else if (status == "下班啦") cdL.Text = "今天下班啦";
            else cdL.Text = "距下班 " + Program.FmtDur(endSec - nSec);

            statusL.Text = status;
            if (status == "上班中" && running) statusL.ForeColor = t.Main;
            else if (status == "下班啦") statusL.ForeColor = t.Off;
            else statusL.ForeColor = t.Idle;

            string moneyText = "¥" + earned.ToString("0.0000", CultureInfo.InvariantCulture);
            moneyL.Text = moneyText;
            // 金额越长字号越小（14/12/10/8pt），在加宽后的 132px 内完整显示
            int wantSize = moneyText.Length <= 8 ? 14 : (moneyText.Length == 9 ? 12 : (moneyText.Length == 10 ? 10 : 8));
            if (wantSize != moneyFontSize)
            {
                moneyFontSize = wantSize;
                moneyL.Font = wantSize == 14 ? moneyFont14 : (wantSize == 12 ? moneyFont12 : (wantSize == 10 ? moneyFont10 : moneyFont8));
            }
            speedL.Text = "¥" + speed.ToString("0.0000", CultureInfo.InvariantCulture) + " / 秒";

            Item it = Program.Items[itemIdx];
            int n = (int)Math.Floor(earned / it.Price);
            if (n >= 1)
                phraseL.Text = "今天的价钱能买 " + n + " " + it.Unit + it.Name + "啦！";
            else
                phraseL.Text = "再过" + Program.FmtDur((it.Price - earned) / speed) + "，今天的价钱就能买 1 " + it.Unit + it.Name + "了。";

            progL.Text = "今日工时进度 " + (worked * 100).ToString("0.0", CultureInfo.InvariantCulture) + "%\n已工作 "
                + workedH.ToString("0.00", CultureInfo.InvariantCulture) + " / " + c.Hours.ToString("0.00", CultureInfo.InvariantCulture) + " 小时";
            bar.Value = (int)(worked * 100);
            monthL.Text = "¥" + monthEarn.ToString("0.0000", CultureInfo.InvariantCulture);
        }
    }

    public class SettingsForm : Form
    {
        NumericUpDown salary, hours;
        TextBox start, end;

        public SettingsForm()
        {
            Text = "设置";
            FormBorderStyle = FormBorderStyle.FixedSingle;
            MaximizeBox = false; MinimizeBox = false;
            StartPosition = FormStartPosition.CenterParent;
            AutoScaleMode = AutoScaleMode.None;
            ClientSize = new Size(300, 215);
            Font = new Font("Microsoft YaHei", 8f);
            BackColor = Program.Themes[Program.cfg.Theme].Bg;

            Label l1 = new Label(); l1.Text = "月薪（元）"; l1.Location = new Point(16, 18); l1.Size = new Size(120, 20); Controls.Add(l1);
            salary = new NumericUpDown(); salary.Location = new Point(150, 15); salary.Size = new Size(134, 24);
            salary.Minimum = 1; salary.Maximum = 10000000; salary.DecimalPlaces = 2;
            salary.Value = Math.Min(salary.Maximum, Math.Max(salary.Minimum, (decimal)Program.cfg.Salary)); Controls.Add(salary);

            Label l2 = new Label(); l2.Text = "每天工作时长（小时）"; l2.Location = new Point(16, 52); l2.Size = new Size(130, 20); Controls.Add(l2);
            hours = new NumericUpDown(); hours.Location = new Point(150, 49); hours.Size = new Size(134, 24);
            hours.Minimum = 1; hours.Maximum = 24; hours.DecimalPlaces = 2;
            hours.Value = Math.Min(hours.Maximum, Math.Max(hours.Minimum, (decimal)Program.cfg.Hours)); Controls.Add(hours);

            Label l3 = new Label(); l3.Text = "上班时间（HH:mm）"; l3.Location = new Point(16, 86); l3.Size = new Size(130, 20); Controls.Add(l3);
            start = new TextBox(); start.Location = new Point(150, 83); start.Size = new Size(134, 24); start.Text = Program.cfg.Start; Controls.Add(start);

            Label l4 = new Label(); l4.Text = "下班时间（HH:mm）"; l4.Location = new Point(16, 120); l4.Size = new Size(130, 20); Controls.Add(l4);
            end = new TextBox(); end.Location = new Point(150, 117); end.Size = new Size(134, 24); end.Text = Program.cfg.End; Controls.Add(end);

            // 固定 Label 尺寸，防止系统字体偏大时撑开
            foreach (Control ctl in Controls)
                if (ctl is Label) ((Label)ctl).AutoSize = false;

            Button ok = new Button(); ok.Text = "保存"; ok.Location = new Point(130, 168); ok.Size = new Size(72, 26);
            ok.Click += delegate
            {
                int s = Program.ParseHM(start.Text), e = Program.ParseHM(end.Text);
                if (s < 0 || e < 0 || e <= s)
                {
                    MessageBox.Show(this, "时间格式应为 HH:mm，且下班时间要晚于上班时间", "提示", MessageBoxButtons.OK, MessageBoxIcon.Warning);
                    return;
                }
                Program.cfg.Salary = (double)salary.Value;
                Program.cfg.Hours = (double)hours.Value;
                Program.cfg.Start = start.Text.Trim();
                Program.cfg.End = end.Text.Trim();
                Program.SaveConfig(Program.cfg);
                DialogResult = DialogResult.OK;
                Close();
            };
            Controls.Add(ok);
            Button cancel = new Button(); cancel.Text = "取消"; cancel.Location = new Point(210, 168); cancel.Size = new Size(72, 26);
            cancel.Click += delegate { DialogResult = DialogResult.Cancel; Close(); };
            Controls.Add(cancel);
        }
    }
}
