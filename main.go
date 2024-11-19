package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

const (
	appName    = "GoNpKill"
	appVersion = "1.0.0"
)

type item struct {
	path string
	size int64
}

func (i item) Title() string       { return i.path }
func (i item) Description() string { return humanize.Bytes(uint64(i.size)) }
func (i item) FilterValue() string { return i.path }

type model struct {
	list         list.Model
	spinner      spinner.Model
	scanning     bool
	selectedPath string
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		list:     list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		spinner:  s,
		scanning: true,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, scanDirectories)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.scanning {
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			m.selectedPath = m.list.SelectedItem().(item).path
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := msg.Width, msg.Height
		m.list.SetSize(h, v-3)
	case scanResultMsg:
		m.scanning = false
		items := []list.Item{}
		for _, dir := range msg {
			items = append(items, dir)
		}
		m.list.SetItems(items)
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.scanning {
		return fmt.Sprintf(
			"\n\n   %s Scanning for node_modules directories...\n\n",
			m.spinner.View(),
		)
	}

	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render(fmt.Sprintf("%s v%s", appName, appVersion))

	return title + "\n\n" + m.list.View()
}

type scanResultMsg []item

func scanDirectories() tea.Msg {
	var dirs []item
	root := "."

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && info.Name() == "node_modules" {
			size, _ := getDirSize(path)
			dirs = append(dirs, item{path: path, size: size})
		}
		return nil
	})

	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].size > dirs[j].size
	})

	return scanResultMsg(dirs)
}

func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	if m, ok := m.(model); ok && m.selectedPath != "" {
		fmt.Printf("Deleting %s\n", m.selectedPath)
		err := os.RemoveAll(m.selectedPath)
		if err != nil {
			fmt.Printf("Error deleting directory: %v\n", err)
		} else {
			fmt.Println("Directory deleted successfully.")
		}
	}
}
