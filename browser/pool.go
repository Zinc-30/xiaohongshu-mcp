package browser

import (
	"math/rand"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/sirupsen/logrus"
	"github.com/xpzouying/headless_browser"
	"github.com/xpzouying/xiaohongshu-mcp/configs"
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
	mu      sync.Mutex
	browser *headless_browser.Browser
	ua      string

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
func (p *BrowserPool) ensureBrowser() {
	if p.browser != nil {
		return
	}

	logrus.Info("创建浏览器实例...")

	opts := []headless_browser.Option{
		headless_browser.WithHeadless(configs.IsHeadless()),
		headless_browser.WithUserAgent(p.ua),
	}

	if binPath := configs.GetBinPath(); binPath != "" {
		opts = append(opts, headless_browser.WithChromeBinPath(binPath))
	}

	p.browser = NewBrowser(configs.IsHeadless(), WithBinPath(configs.GetBinPath()))
}

// GetPage 从复用的浏览器实例中获取一个新 page，已配置视口、DPR、时区
func (p *BrowserPool) GetPage() *rod.Page {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.ensureBrowser()
	page := p.browser.NewPage()

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
		p.browser.Close()
		p.browser = nil
		logrus.Info("浏览器实例已关闭")
	}
}
