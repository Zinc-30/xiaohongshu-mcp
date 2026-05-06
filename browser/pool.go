package browser

import (
	"encoding/json"
	"math/rand"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/sirupsen/logrus"
	"github.com/xpzouying/xiaohongshu-mcp/configs"
	"github.com/xpzouying/xiaohongshu-mcp/cookies"
)

// 常见桌面分辨率
var viewportPresets = []struct {
	Width, Height int
}{
	{1920, 1080},
	{2560, 1440},
	{1680, 1050},
	{1440, 900},
}

// BrowserPool 复用单个 Chromium 实例，避免每次请求创建新进程
type BrowserPool struct {
	mu       sync.Mutex
	browser  *rod.Browser
	launcher *launcher.Launcher
	ua       string

	// 同一实例生命周期内保持不变的视口参数
	viewport struct {
		width, height int
	}
}

var (
	globalPool     *BrowserPool
	globalPoolOnce sync.Once
)

// GetPool 获取全局浏览器池单例
func GetPool() *BrowserPool {
	globalPoolOnce.Do(func() {
		globalPool = &BrowserPool{}
		globalPool.ua = GenerateUserAgent()
		vp := viewportPresets[rand.Intn(len(viewportPresets))]
		globalPool.viewport.width = vp.Width
		globalPool.viewport.height = vp.Height
		logrus.Infof("浏览器池初始化: UA=%s, 视口=%dx%d", globalPool.ua, vp.Width, vp.Height)
	})
	return globalPool
}

// ensureBrowser 确保浏览器实例存在，不存在则创建
// 不使用 --no-sandbox，避免被反爬系统检测
func (p *BrowserPool) ensureBrowser() {
	if p.browser != nil {
		return
	}

	logrus.Info("创建浏览器实例...")

	l := launcher.New().
		Headless(configs.IsHeadless()).
		Set("user-agent", p.ua).
		// 移除反爬可检测的标记
		Delete("--enable-automation").
		Set("--disable-blink-features", "AutomationControlled")

	if binPath := configs.GetBinPath(); binPath != "" {
		l = l.Bin(binPath)
	}

	url := l.MustLaunch()

	browser := rod.New().ControlURL(url).MustConnect()

	// 加载 cookies
	cookiePath := cookies.GetCookiesFilePath()
	cookieLoader := cookies.NewLoadCookie(cookiePath)
	if data, err := cookieLoader.LoadCookies(); err == nil {
		var cks []*proto.NetworkCookie
		if err := json.Unmarshal(data, &cks); err == nil {
			browser.MustSetCookies(cks...)
		}
	}

	p.browser = browser
	p.launcher = l
}

// GetPage 从复用的浏览器实例中获取一个新 page，已配置视口、DPR、时区
// 调用前自动执行 rate limiter 等待，模拟人类请求节奏
func (p *BrowserPool) GetPage() *rod.Page {
	// 请求频率限制：确保连续调用间有随机间隔
	GetRateLimiter().Wait()

	p.mu.Lock()
	defer p.mu.Unlock()

	p.ensureBrowser()
	page := p.browser.MustPage("")

	// 设置视口 + Retina DPR
	err := proto.EmulationSetDeviceMetricsOverride{
		Width:             p.viewport.width,
		Height:            p.viewport.height,
		DeviceScaleFactor: 2,
		Mobile:            false,
	}.Call(page)
	if err != nil {
		logrus.Warnf("设置视口失败: %v", err)
	}

	// 设置中国时区
	err = proto.EmulationSetTimezoneOverride{
		TimezoneID: "Asia/Shanghai",
	}.Call(page)
	if err != nil {
		logrus.Warnf("设置时区失败: %v", err)
	}

	return page
}

// ReleasePage 关闭 page 但保留浏览器实例
func (p *BrowserPool) ReleasePage(page *rod.Page) {
	if page != nil {
		_ = page.Close()
	}
}

// Close 关闭浏览器实例（仅在服务停止时调用）
func (p *BrowserPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.browser != nil {
		p.browser.MustClose()
		p.browser = nil
		logrus.Info("浏览器实例已关闭")
	}
}
