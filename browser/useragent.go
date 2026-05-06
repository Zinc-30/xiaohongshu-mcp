package browser

import (
	"fmt"
	"math/rand"
	"time"
)

// Chrome 124 发布于 2024-04-16，此后每 4 周一个大版本
var chromeBaseVersion = 124
var chromeBaseDate = time.Date(2024, 4, 16, 0, 0, 0, 0, time.UTC)

var macOSVersions = []string{
	"10_15_7",
	"13_6_7",
	"14_4_1",
	"14_5",
	"14_6_1",
	"15_0",
	"15_1",
	"15_2",
}

// 真实 Chrome 小版本号范围（从公开 release 数据中采样）
// 格式: Chrome/{major}.0.{build}.{patch}
// build 通常在 6000-7000 范围，patch 在 0-300 范围
var chromeBuildRanges = [][2]int{
	{6613, 6723}, // Chrome 124-128 range
	{6778, 6845}, // Chrome 129-132 range
	{6900, 6998}, // Chrome 133-136 range
}

// GenerateUserAgent 根据当前日期推算合理的 Chrome 主版本号，生成真实格式的 UA
func GenerateUserAgent() string {
	weeks := int(time.Since(chromeBaseDate).Hours() / (24 * 7))
	majorVersion := chromeBaseVersion + weeks/4

	macOS := macOSVersions[rand.Intn(len(macOSVersions))]

	// 选择合适的 build 范围
	rangeIdx := (majorVersion - chromeBaseVersion) / 4
	if rangeIdx >= len(chromeBuildRanges) {
		rangeIdx = len(chromeBuildRanges) - 1
	}
	buildRange := chromeBuildRanges[rangeIdx]
	build := buildRange[0] + rand.Intn(buildRange[1]-buildRange[0]+1)
	patch := rand.Intn(200)

	return fmt.Sprintf(
		"Mozilla/5.0 (Macintosh; Intel Mac OS X %s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.0.%d.%d Safari/537.36",
		macOS, majorVersion, build, patch,
	)
}
