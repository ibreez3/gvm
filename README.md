gvm 是一个使用 Go 编写的 Go 版本管理器，支持远程安装、版本切换、本地与远程列表、版本搜索，以及自动初始化 shell 环境。

特性
- 安装指定版本：`gvm install <version>`（支持 `1.25.5` 或 `go1.25.5`）
- 切换当前版本：`gvm use <version>`（建立 `~/.gvm/goroot` 软链指向 `~/.gvm/go<version>`）
- 查看本地版本：`gvm list -l`
- 查看远程版本：`gvm list -r -n 30`（默认返回最新 30 个稳定版，过滤 rc/beta）
- 查看当前使用：`gvm current`
- 搜索子版本：`gvm search 1.21`（返回 `1.21.x` 的所有稳定子版本）

安装与构建
- 从源码构建：
  - `go build -o gvm ./cmd/gvm`
  - 首次使用前执行：`./gvm init`
- 通过 Release 安装：在 GitHub Releases 页面下载对应平台的压缩包，解压得到 `gvm` 可执行文件并加入 PATH。

用法示例
- 初始化：`gvm init`
- 安装并切换：
  - `gvm install 1.24.5`
  - `gvm use 1.24.5`
- 列表：
  - 本地：`gvm list -l`
  - 远程：`gvm list -r -n 30`
- 当前版本：`gvm current`
- 搜索：`gvm search 1.21`

目录结构与环境
- 安装目录：`$HOME/.gvm`，每个版本为一个目录：`~/.gvm/go<version>`
- 当前版本：`~/.gvm/goroot -> ~/.gvm/go<version>`（软链接）
- 环境文件：`~/.gvm/.gvmrc`，内容：
  - `export GOROOT=$HOME/.gvm/goroot`
  - `export PATH=$PATH:$GOROOT/bin`
  - `export GOPATH=$HOME/go`
  - `export GOBIN=$GOPATH/bin`
  - `export GOPROXY=https://goproxy.cn,direct`
- Shell 注入：自动检测 `~/.zshrc` 或 `~/.bashrc` 并追加：
  - `# gvm shell setup`
  - `if [ -f "$HOME/.gvm/.gvmrc" ]; then`
  - `    source "$HOME/.gvm/.gvmrc"`
  - `fi`

Release 工作流
- 推送符合 `v*.*.*` 的 tag 将自动触发构建与发布（多平台：Linux、macOS、Windows；架构：amd64、arm64）。
- 工作流使用 Go 官方环境与 GoReleaser 构建压缩包并上传到 GitHub Release。
- 手动触发也支持：在 GitHub Actions 中选择 `Release` 工作流，点击 `Run workflow`。

遥测与下载源
- 远程版本列表与安装包源来自官方 `https://go.dev/dl`，安全且稳定。

致谢
- Go 语言具有快速编译、并发原生支持与强大的标准库，适用于 CLI 工具与发布流程。
