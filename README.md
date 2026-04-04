# API Checker

一个用于批量检测和管理 API Key 的桌面工具。

## 功能特点

- **批量检测**: 支持批量验证 OpenAI 和 Anthropic API Key 的有效性
- **模型获取**: 自动获取每个有效 API Key 可用的模型列表
- **智能测试**: 测试模型的智能度并缓存结果
- **安全存储**: 使用机器码或用户自定义密码加密存储配置，确保敏感信息安全
- **聊天调试**: 内置聊天窗口，可直接与模型对话进行调试
- **快速复制**: 双击表格单元格快速复制 API Key 或 Base URL
- **多协议支持**: 支持 OpenAI、Anthropic 及兼容 OpenAI API 格式的自定义 API
- **提示词配置**: 可自定义测试提示词，用于智能度测试
- **启动密码保护**: 可设置启动密码，防止未授权访问 API Key

## 支持的协议

- OpenAI
- Anthropic
- 兼容 OpenAI API 、Anthropic API格式的自定义 API（如阿里云 DashScope、DeepSeek 等）

## 编译运行

### 前置要求

- Go 1.23+
- Node.js 22+
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### 开发模式

```bash
wails dev
```

### 生产构建

```bash
# Windows
wails build -platform windows/amd64

# Linux (使用 webkit2gtk-4.1)
wails build -tags webkit2_41 -platform linux/amd64

# macOS
wails build -platform darwin/universal
```

## 使用说明

### 添加 API Key

1. 点击底部的"添加"按钮
2. 选择协议类型（OpenAI/Anthropic/自定义）
3. 输入别名（用于标识）
4. 输入 API Key
5. 输入 Base URL（如需自定义）
6. 点击"确定"保存

### 检测 API Key

- **单个检测**: 点击 API Key 行的"检测"按钮
- **批量检测**: 点击底部的"批量检测"按钮

### 查看模型列表

1. 点击 API Key 行的"查看模型"按钮
2. 在弹出的窗口中可查看该 API Key 可用的模型列表
3. 可点击"测试提示"按钮测试提示词是否正确

### 聊天调试

1. 在底部选择已检测通过的 API Key
2. 选择模型
3. 在输入框中输入消息，按 Enter 发送

### 配置提示词

1. 点击顶部菜单的"设置"按钮
2. 在"提示词"输入框中输入测试提示词
3. 点击"保存"按钮

### 设置启动密码

1. 点击顶部菜单的"设置"按钮
2. 在"启动密码"输入框中输入密码
3. 点击"保存"按钮
4. 下次启动程序时需要输入密码才能查看 API Key

## 安全说明

API Key 使用 AES-GCM 加密算法进行加密存储，加密密钥基于为空时由机器码派生，只有在当前机器上才能解密查看 API Key 内容。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License
