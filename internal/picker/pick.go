package picker

import (
	"errors"
	"fmt"
	"git-msg-picker/internal/conf"
	"github.com/charmbracelet/bubbletea"
	"sync"
)

var (
	pickModelIns *PickModel
	pickOnce     sync.Once
)

func GetPickerModel() PickModel {
	prefixes, err := conf.LoadPrefixes()
	if err != nil {
		panic(err)
	}
	pickOnce.Do(func() {
		pickModelIns = &PickModel{
			picker: Picker{
				options:  prefixes,
				selected: make(map[int]struct{}),
			},
			cursor: 0,
		}
	})
	return *pickModelIns
}

type Picker struct {
	options  []conf.CommitPrefix
	selected map[int]struct{}
}

type PickModel struct {
	picker Picker
	cursor int
}

func (pm PickModel) Init() tea.Cmd {
	return nil
}

func (pm PickModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return pm, tea.Quit
		}
		
		switch msg.String() {
		case "q":
			return pm, tea.Quit
		
		case "j", "down":
			pm.cursor++
			if pm.cursor >= pm.picker.Len() {
				pm.cursor = 0
			}
		
		case "k", "up":
			pm.cursor--
			if pm.cursor < 0 {
				pm.cursor = pm.picker.Len() - 1
			}
		
		case "enter":
			if _, ok := pm.picker.selected[pm.cursor]; ok {
				return GetCommitModel(
					"请输入commit内容",
					pm.picker.options[pm.cursor].Prefix,
				), nil
			}
			
			pm.picker.selected = make(map[int]struct{})
			err := pm.picker.Pick(pm.cursor)
			if err != nil {
				return pm, tea.Quit
			}
			
		}
	}
	
	return pm, nil
}

func (pm PickModel) View() string {
	var s string
	
	for i, item := range pm.picker.options {
		cursor := " "
		if pm.cursor == i {
			cursor = "→"
		}
		
		checked := " "
		if _, ok := pm.picker.selected[i]; ok {
			checked = "✓"
		}
		
		s += fmt.Sprintf("%s [%s]%s: %s\n", cursor, checked, item.Prefix, item.Description)
	}
	
	return s
}

func (p Picker) Len() int {
	return len(p.options)
}

func (p Picker) Pick(idx int) error {
	if idx < 0 || idx > p.Len() {
		return errors.New("超出选项范围")
	}
	
	p.selected[idx] = struct{}{}
	
	return nil
}

func (p Picker) UnPick(idx int) error {
	if idx < 0 || idx > p.Len() {
		return errors.New("超出选项范围")
	}
	
	delete(p.selected, idx)
	
	return nil
}
