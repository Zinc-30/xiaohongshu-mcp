package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sirupsen/logrus"
	"github.com/xpzouying/xiaohongshu-mcp/configs"
)

func main() {
	var (
		headless bool
		binPath  string
		port     string
		stdio    bool
	)
	flag.BoolVar(&headless, "headless", true, "是否无头模式")
	flag.StringVar(&binPath, "bin", "", "浏览器二进制文件路径")
	flag.StringVar(&port, "port", ":18060", "端口")
	flag.BoolVar(&stdio, "stdio", false, "使用 stdio 传输模式（供 Claude Code 等客户端自动管理进程）")
	flag.Parse()

	if len(binPath) == 0 {
		binPath = os.Getenv("ROD_BROWSER_BIN")
	}

	configs.InitHeadless(headless)
	configs.SetBinPath(binPath)

	xiaohongshuService := NewXiaohongshuService()

	if stdio {
		// stdio 模式：通过 stdin/stdout 通信，由客户端管理进程生命周期
		logrus.SetOutput(os.Stderr)
		appServer := NewAppServer(xiaohongshuService)

		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		if err := appServer.mcpServer.Run(ctx, &mcp.StdioTransport{}); err != nil {
			logrus.Fatalf("stdio server exited: %v", err)
		}
	} else {
		// HTTP 模式：启动 HTTP 服务器
		appServer := NewAppServer(xiaohongshuService)
		if err := appServer.Start(port); err != nil {
			logrus.Fatalf("failed to run server: %v", err)
		}
	}
}
