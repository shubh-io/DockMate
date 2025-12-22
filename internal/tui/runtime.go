package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	runtimeTitleStyle        = lipgloss.NewStyle().MarginLeft(2)
	runtimeItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	runtimeSelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	runtimePaginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	runtimeHelpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type runtimeItem string

func (i runtimeItem) FilterValue() string { return "" }

type runtimeItemDelegate struct{}

func (d runtimeItemDelegate) Height() int                             { return 1 }
func (d runtimeItemDelegate) Spacing() int                            { return 0 }
func (d runtimeItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d runtimeItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(runtimeItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := runtimeItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return runtimeSelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// RuntimeSelectionModel handles runtime selection (Docker/Podman)
type RuntimeSelectionModel struct {
	list   list.Model
	choice string
}

func NewRuntimeSelectionModel() RuntimeSelectionModel {
	items := []list.Item{
		runtimeItem("Docker (default)"),
		runtimeItem("Podman"),
	}

	const defaultWidth = 20

	l := list.New(items, runtimeItemDelegate{}, defaultWidth, listHeight)
	l.Title = "Which runtime do you want to use in DockMate🐳? (don't worry, you can change this later in settings)"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = runtimeTitleStyle
	l.Styles.PaginationStyle = runtimePaginationStyle
	l.Styles.HelpStyle = runtimeHelpStyle

	return RuntimeSelectionModel{list: l}
}

func (m RuntimeSelectionModel) Init() tea.Cmd {
	return nil
}

func (m RuntimeSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(runtimeItem)

			if ok {
				if strings.Contains(string(i), "Docker") {
					m.choice = "docker"
				} else {
					m.choice = strings.ToLower(string(i))
				}
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m RuntimeSelectionModel) View() string {
	return "\n" + m.list.View()
}

// GetChoice returns the selected runtime
func (m RuntimeSelectionModel) GetChoice() string {
	return m.choice
}
