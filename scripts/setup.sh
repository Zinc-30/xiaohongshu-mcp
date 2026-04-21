#!/bin/bash
# xiaohongshu-mcp 安装脚本 — 从 GitHub Releases 下载预编译二进制
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
BIN_DIR="$PROJECT_DIR/bin"
REPO="Zinc-30/xiaohongshu-mcp"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

info()  { echo -e "${CYAN}[INFO]${NC} $*"; }
ok()    { echo -e "${GREEN}  ✓${NC} $*"; }
warn()  { echo -e "${YELLOW}  ⚠${NC} $*"; }
fail()  { echo -e "${RED}  ✗${NC} $*"; }

# 检测平台
detect_platform() {
    local os arch
    os="$(uname -s | tr '[:upper:]' '[:lower:]')"
    arch="$(uname -m)"

    case "$arch" in
        x86_64|amd64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *) fail "不支持的架构: $arch"; exit 1 ;;
    esac

    echo "${os}-${arch}"
}

PLATFORM=$(detect_platform)
ARCHIVE_NAME="xiaohongshu-mcp-${PLATFORM}.tar.gz"

info "安装 xiaohongshu-mcp (${PLATFORM})..."

# 获取最新 Release 下载 URL
info "查找最新版本..."
DOWNLOAD_URL=""

if command -v gh &>/dev/null; then
    # 优先用 gh CLI（已认证，不受 API 限流）
    DOWNLOAD_URL=$(gh release view --repo "${REPO}" --json assets \
        --jq ".assets[] | select(.name == \"${ARCHIVE_NAME}\") | .url" 2>/dev/null || echo "")
fi

if [ -z "$DOWNLOAD_URL" ]; then
    # 回退到 GitHub API
    DOWNLOAD_URL=$(curl -sf "https://api.github.com/repos/${REPO}/releases/latest" \
        | grep "browser_download_url.*${ARCHIVE_NAME}" \
        | head -1 \
        | cut -d '"' -f 4 || echo "")
fi

if [ -z "$DOWNLOAD_URL" ]; then
    DOWNLOAD_URL=$(curl -sf "https://api.github.com/repos/${REPO}/releases" \
        | grep "browser_download_url.*${ARCHIVE_NAME}" \
        | head -1 \
        | cut -d '"' -f 4 || echo "")
fi

if [ -z "$DOWNLOAD_URL" ]; then
    warn "未找到预编译二进制，尝试本地编译..."
    if ! command -v go &>/dev/null; then
        fail "未找到预编译二进制且 Go 未安装"
        fail "请先安装 Go (brew install go) 或等待 GitHub Release 构建完成后重试"
        exit 1
    fi
    mkdir -p "$BIN_DIR"
    cd "$PROJECT_DIR"
    go build -o "$BIN_DIR/xiaohongshu-mcp" .
    ok "本地编译成功"
    exit 0
fi

info "下载: $DOWNLOAD_URL"

mkdir -p "$BIN_DIR"
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

curl -sL "$DOWNLOAD_URL" -o "$TMP_DIR/$ARCHIVE_NAME"
tar xzf "$TMP_DIR/$ARCHIVE_NAME" -C "$TMP_DIR"

cp "$TMP_DIR/xiaohongshu-mcp-${PLATFORM}" "$BIN_DIR/xiaohongshu-mcp"
chmod +x "$BIN_DIR/xiaohongshu-mcp"

# 如果有登录工具也拷贝
if [ -f "$TMP_DIR/xiaohongshu-login-${PLATFORM}" ]; then
    cp "$TMP_DIR/xiaohongshu-login-${PLATFORM}" "$BIN_DIR/xiaohongshu-login"
    chmod +x "$BIN_DIR/xiaohongshu-login"
fi

ok "安装成功: $BIN_DIR/xiaohongshu-mcp"
