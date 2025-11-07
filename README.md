# git-msg-picker

`git-msg-picker` 是一个基于 Go 开发的交互式 Git 提交信息工具，通过终端界面快速选择规范的提交前缀（如 `feat`、`fix` 等），快速输入提交信息，简化 `git commit` 流程。工具名缩写为 `cip`（Commitmit Input Picker），方便高频使用。


## 功能特点
- 交互式终端界面，支持 `j/k` 键或方向键上下移动选择提交前缀
- 可视化高亮选中项，直观清晰
- 自动拼接前缀与提交信息，执行 `git commit -m "前缀: 信息"`


## 库依赖
| 库名                | 地址                                      | 作用说明                                                                 |
|---------------------|-------------------------------------------|--------------------------------------------------------------------------|
| bubbletea           | https://github.com/charmbracelet/bubbletea | 终端交互框架，处理键盘输入、UI 渲染和状态管理，是工具的核心交互引擎       |
| go-git              | https://github.com/go-git/go-git/v5        | 纯 Go 实现的 Git 操作库，用于在代码中直接执行 Git 提交命令，无需调用系统命令 |
| termenv             | https://github.com/muesli/termenv         | 终端样式处理库，负责选中项高亮、颜色适配等 UI 美化功能                   |
| urfave/cli/v2       | https://github.com/urfave/cli/v2          | 命令行参数解析库，支持工具的参数、选项和帮助信息生成（如 `cip --help`）  |