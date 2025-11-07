package git

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"os"
	"time"
)

// Operator 封装 Git 操作的结构体
type Operator struct {
	repoPath string          // 仓库目录路径
	repo     *git.Repository // 仓库实例
	worktree *git.Worktree   // 工作区实例
}

// NewGitOperator 创建 Git 操作实例
// dir: 仓库目录（如 "." 表示当前目录）
func NewGitOperator(dir string) (*Operator, error) {
	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("目录不存在: %s", dir)
	}
	
	return &Operator{repoPath: dir}, nil
}

// Open 打开已存在的 Git 仓库（类似 cd 到仓库目录）
func (g *Operator) Open() error {
	repo, err := git.PlainOpen(g.repoPath)
	if err != nil {
		return fmt.Errorf("打开仓库失败（请确保是 Git 仓库）: %w", err)
	}
	g.repo = repo
	
	// 获取工作区
	wt, err := g.repo.Worktree()
	if err != nil {
		return fmt.Errorf("获取工作区失败: %w", err)
	}
	g.worktree = wt
	return nil
}

// Commit 提交修改（类似 git commit -m）
// message: 提交信息
func (g *Operator) Commit(message string) (string, error) {
	if g.worktree == nil || g.repo == nil {
		return "", fmt.Errorf("请先初始化或打开仓库")
	}
	
	// 从 Git 配置读取提交者信息
	conf, err := config.LoadConfig(config.GlobalScope)
	if err != nil {
		conf, err = config.LoadConfig(config.LocalScope)
		if err != nil {
			return "", fmt.Errorf("读取 Git 配置失败: %w", err)
		}
	}
	
	if err != nil {
		return "", fmt.Errorf("读取 Git 配置失败: %w", err)
	}
	
	// 提交者信息（默认值兜底）
	author := &object.Signature{
		Name:  "Unknown User",
		Email: "unknown@example.com",
		When:  time.Now(),
	}
	
	if conf.User.Name != "" {
		author.Name = conf.User.Name
	}
	if conf.User.Email != "" {
		author.Email = conf.User.Email
	}
	
	// 执行提交
	commitHash, err := g.worktree.Commit(message, &git.CommitOptions{
		Author: author,
	})
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			return "", fmt.Errorf("没有可提交的修改（工作区已干净）")
		}
		return "", err
	}
	
	return commitHash.String(), nil
}
