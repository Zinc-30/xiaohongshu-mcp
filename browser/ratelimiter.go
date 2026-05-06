package browser

import (
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// RateLimiter 全局请求频率限制器，确保连续请求之间有随机间隔
type RateLimiter struct {
	mu       sync.Mutex
	lastCall time.Time
	// 最小间隔 (ms)
	minInterval int
	// 最大间隔 (ms)
	maxInterval int
}

var (
	globalLimiter     *RateLimiter
	globalLimiterOnce sync.Once
)

// GetRateLimiter 获取全局请求限制器
func GetRateLimiter() *RateLimiter {
	globalLimiterOnce.Do(func() {
		globalLimiter = &RateLimiter{
			minInterval: 500,   // 最少等 0.5 秒
			maxInterval: 1500,  // 最多等 1.5 秒
		}
	})
	return globalLimiter
}

// Wait 在发起请求前调用，自动等待随机间隔以模拟人类节奏
func (r *RateLimiter) Wait() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.lastCall.IsZero() {
		r.lastCall = time.Now()
		return
	}

	// 计算需要等待的随机间隔
	interval := time.Duration(r.minInterval+rand.Intn(r.maxInterval-r.minInterval)) * time.Millisecond
	elapsed := time.Since(r.lastCall)

	if elapsed < interval {
		waitTime := interval - elapsed
		logrus.Debugf("rate limiter: 等待 %dms", waitTime.Milliseconds())
		time.Sleep(waitTime)
	}

	r.lastCall = time.Now()
}
