package cookies

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Cookier interface {
	LoadCookies() ([]byte, error)
	SaveCookies(data []byte) error
	DeleteCookies() error
}

type localCookie struct {
	path string
}

func NewLoadCookie(path string) Cookier {
	if path == "" {
		panic("path is required")
	}

	return &localCookie{
		path: path,
	}
}

// LoadCookies 从文件中加载 cookies。
func (c *localCookie) LoadCookies() ([]byte, error) {

	data, err := os.ReadFile(c.path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read cookies from tmp file")
	}

	return data, nil
}

// SaveCookies 保存 cookies 到文件中。
func (c *localCookie) SaveCookies(data []byte) error {
	return os.WriteFile(c.path, data, 0644)
}

// DeleteCookies 删除 cookies 文件。
func (c *localCookie) DeleteCookies() error {
	if _, err := os.Stat(c.path); os.IsNotExist(err) {
		// 文件不存在，返回 nil（认为已经删除）
		return nil
	}
	return os.Remove(c.path)
}

// GetCookiesFilePath 获取 cookies 文件路径。
// 优先级：$COOKIES_PATH 环境变量 > ~/.xiaohongshu-mcp/cookies.json
func GetCookiesFilePath() string {
	// 环境变量优先（Docker/自定义场景）
	if path := os.Getenv("COOKIES_PATH"); path != "" {
		return path
	}

	// 默认路径：~/.xiaohongshu-mcp/cookies.json
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "cookies.json"
	}

	cookieDir := filepath.Join(homeDir, ".xiaohongshu-mcp")
	_ = os.MkdirAll(cookieDir, 0700)

	return filepath.Join(cookieDir, "cookies.json")
}
