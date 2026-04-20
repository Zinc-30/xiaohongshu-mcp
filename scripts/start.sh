#!/bin/bash
# xiaohongshu-mcp 启动脚本（stdio 模式）
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
BIN_PATH="$PROJECT_DIR/bin/xiaohongshu-mcp"

# 如果二进制不存在，先构建
if [ ! -f "$BIN_PATH" ]; then
    bash "$SCRIPT_DIR/setup.sh" >&2
fi

export COOKIES_PATH="${COOKIES_PATH:-$PROJECT_DIR/cookies.json}"

exec "$BIN_PATH" -headless=true -stdio
