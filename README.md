# aui - AI Agent User Interface

Orchestrate multiple AI agents from your terminal. Compare responses, manage contexts, and track costs across Claude, GPT-4, and Gemini.

## What is aui?

`aui` is a terminal-based tool for managing multiple AI assistants simultaneously. Instead of juggling browser tabs and losing context between sessions, you get a unified interface designed for developers who live in the terminal.

## Key Features

### 🤖 Multi-Agent Orchestration
Run multiple AI models in parallel. Send the same prompt to Claude and GPT-4, compare their responses side-by-side, and choose the best solution.

### 📁 Smart Context Management
Build context from your codebase intelligently. Select files, track token usage, and reuse contexts across sessions. No more copy-pasting files into chat windows.

### 💰 Cost Tracking
See exactly how much each query costs across different providers. Track usage patterns and optimize your AI spend.

### ⚡ Built for Speed
Keyboard-driven interface with vim-like navigation. Stream responses in real-time. Everything happens in your terminal.

## Installation

### Pre-built Binaries
Coming soon - check [Releases](https://github.com/yourusername/aui/releases)

### From Source
```bash
go install github.com/yourusername/aui/cmd/aui@latest
```

## Quick Start

1. **Configure your AI providers:**
```bash
aui config set anthropic.key "sk-ant-..."
aui config set openai.key "sk-..."
aui config set google.key "AIza..."
```

2. **Launch aui:**
```bash
aui
```

3. **Use keyboard shortcuts:**
- `Tab` / `l` - Next tab
- `Shift+Tab` / `h` - Previous tab  
- `a` - Add a new agent
- `c` - Create context from files
- `Enter` - Send prompt to agents
- `q` - Quit

## Common Workflows

### Compare Model Responses
1. Add Claude and GPT-4 agents
2. Load your code context
3. Ask both to solve the same problem
4. View responses side-by-side
5. Choose the best approach

### Debug with Context
1. Create context from your error logs and source files
2. Send to Claude with the error message
3. Get targeted suggestions based on your actual code

### Optimize Costs
1. Start with cheaper models (Gemini Flash, GPT-3.5)
2. Escalate to advanced models only when needed
3. Track token usage to identify expensive patterns

## Requirements

- Terminal with UTF-8 support
- API keys for AI providers you want to use
- macOS, Linux, or Windows

## Documentation

- [User Guide](docs/USER_GUIDE.md) - Detailed usage instructions
- [Configuration](docs/CONFIG.md) - Configuration options
- [Keyboard Shortcuts](docs/SHORTCUTS.md) - Complete shortcut reference

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Support

- [GitHub Issues](https://github.com/yourusername/aui/issues) - Bug reports and feature requests
- [Discussions](https://github.com/yourusername/aui/discussions) - Questions and community help

---

*aui - Bringing AI orchestration to your terminal*