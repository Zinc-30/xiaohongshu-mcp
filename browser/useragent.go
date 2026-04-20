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

// GenerateUserAgent 根据当前日期推算合理的 Chrome 主版本号
func GenerateUserAgent() string {
	weeks := int(time.Since(chromeBaseDate).Hours() / (24 * 7))
	majorVersion := chromeBaseVersion + weeks/4

	macOS := macOSVersions[rand.Intn(len(macOSVersions))]

	return fmt.Sprintf(
		"Mozilla/5.0 (Macintosh; Intel Mac OS X %s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.0.0.0 Safari/537.36",
		macOS, majorVersion,
	)
}
