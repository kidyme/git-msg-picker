package picker

import (
	"errors"
	"fmt"
	"git-msg-picker/internal/git"
	"strings"
	"sync"
	
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
)

const DIR = "."

var (
	commitModelIns *CommitModel
	commitOnce     sync.Once
)

func GetCommitModel(prompt, prefix string) CommitModel {
	commitOnce.Do(func() {
		ti := textinput.New()
		ti.Placeholder = prompt
		ti.CharLimit = 200
		ti.Width = 50
		ti.Prompt = ""
		ti.Focus()
		
		commitModelIns = &CommitModel{
			prefix:    prefix,
			textInput: ti,
			state:     stateInputting,
		}
	})
	
	return *commitModelIns
}

const (
	stateInputting  = iota // 输入中
	stateProcessing        // 提交处理中
	stateSuccess           // 提交成功
	stateError             // 提交失败
)

type CommitModel struct {
	prefix     string
	textInput  textinput.Model
	state      int    // 当前状态
	commitHash string // 成功时存储Commit Hash
	errorMsg   string // 失败时存储错误信息
}

type commitResultMsg struct {
	hash string
	err  error
}

func (cm CommitModel) Init() tea.Cmd {
	return textinput.Blink
}

func (cm CommitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if cm.state != stateInputting {
			return cm, tea.Quit
		}
		
		switch msg.Type {
		case tea.KeyEnter:
			cm.state = stateProcessing
			return cm, cm.CommitCmd()
		
		case tea.KeyEscape:
			cm.textInput.SetValue("")
			updatedTI, cmd := cm.textInput.Update(msg)
			cm.textInput = updatedTI
			return cm, cmd
		
		case tea.KeyCtrlQ:
			return GetPickerModel(), nil
		
		case tea.KeyCtrlC:
			return cm, tea.Quit
		
		default:
			updatedTI, cmd := cm.textInput.Update(msg)
			cm.textInput = updatedTI
			return cm, cmd
		}
	
	case commitResultMsg:
		if msg.err != nil {
			cm.state = stateError
			cm.errorMsg = msg.err.Error()
		} else {
			cm.state = stateSuccess
			cm.commitHash = msg.hash
		}
		return cm, tea.Quit
	
	default:
		updatedTI, cmd := cm.textInput.Update(msg)
		cm.textInput = updatedTI
		return cm, cmd
	}
}

func (cm CommitModel) View() string {
	switch cm.state {
	case stateInputting:
		return fmt.Sprintf(
			"%s: %s\n",
			cm.prefix,
			cm.textInput.View(),
		)
	
	case stateProcessing:
		return "正在提交代码...\n"
	
	case stateSuccess:
		return fmt.Sprintf(
			"✅ 提交成功！\nCommit Hash: %s\n",
			cm.commitHash,
		)
	
	case stateError:
		return fmt.Sprintf(
			"❌ 提交失败！\nError: %s\n",
			wrapText(cm.errorMsg, 100), // 错误信息自动换行，避免超出终端宽度
		)
	
	default:
		return ""
	}
}

func (cm CommitModel) CommitCmd() tea.Cmd {
	return func() tea.Msg {
		inputMsg := strings.TrimSpace(cm.textInput.Value())
		if inputMsg == "" {
			return commitResultMsg{err: errors.New("请输入commit message")}
		}
		
		operator, err := git.NewGitOperator(DIR)
		if err != nil {
			return commitResultMsg{err: err}
		}
		
		err = operator.Open()
		if err != nil {
			return commitResultMsg{err: err}
		}
		
		commitHash, err := operator.Commit(cm.Message())
		if err != nil {
			return commitResultMsg{err: err}
		}
		
		return commitResultMsg{hash: commitHash}
	}
}

func (cm CommitModel) Message() string {
	return fmt.Sprintf("%s: %s", cm.prefix, cm.textInput.Value())
}

func wrapText(text string, width int) string {
	lines := strings.Split(text, "\n")
	var wrapped []string
	for _, line := range lines {
		if len(line) <= width {
			wrapped = append(wrapped, line)
			continue
		}
		for i := 0; i < len(line); i += width {
			end := i + width
			if end > len(line) {
				end = len(line)
			}
			wrapped = append(wrapped, line[i:end])
		}
	}
	return strings.Join(wrapped, "\n")
}
