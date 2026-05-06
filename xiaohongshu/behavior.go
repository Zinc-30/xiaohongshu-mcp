package xiaohongshu

import (
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/sirupsen/logrus"
)

// simulateHumanPresence 模拟人类在页面上的基本存在感：随机鼠标移动 + 短暂停留
func simulateHumanPresence(page *rod.Page) {
	viewWidth := page.MustEval(`() => window.innerWidth`).Int()
	viewHeight := page.MustEval(`() => window.innerHeight`).Int()

	// 1-2 次随机鼠标移动
	moves := 1 + rand.Intn(2)
	for i := 0; i < moves; i++ {
		x := float64(100 + rand.Intn(viewWidth-200))
		y := float64(100 + rand.Intn(viewHeight-200))
		page.Mouse.MustMoveTo(x, y)
		time.Sleep(time.Duration(80+rand.Intn(150)) * time.Millisecond)
	}

	logrus.Debugf("人类存在模拟: %d 次鼠标移动", moves)
}

// simulatePageBrowsing 模拟人类浏览页面：小幅滚动 + 鼠标轨迹 + 随机停顿
func simulatePageBrowsing(page *rod.Page) {
	simulateHumanPresence(page)

	// 1 次小幅滚动
	delta := 100 + rand.Intn(200)
	page.MustEval(`(d) => window.scrollBy(0, d)`, delta)
	time.Sleep(time.Duration(200+rand.Intn(300)) * time.Millisecond)

	// 滚动后一次鼠标移动
	viewWidth := page.MustEval(`() => window.innerWidth`).Int()
	x := float64(200 + rand.Intn(viewWidth-400))
	y := float64(200 + rand.Intn(300))
	page.Mouse.MustMoveTo(x, y)
	time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)

	// 回到顶部
	page.MustEval(`() => window.scrollTo(0, 0)`)
	time.Sleep(time.Duration(150+rand.Intn(250)) * time.Millisecond)

	// 短暂阅读停顿
	readDelay := 500 + rand.Intn(800)
	time.Sleep(time.Duration(readDelay) * time.Millisecond)

	logrus.Debugf("页面浏览模拟: 阅读延迟 %dms", readDelay)
}
