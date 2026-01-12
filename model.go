package main

import (
	"fmt"
	// "path/filepath"
	// // "strings"
	// "time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type View int

const (
	titleView View = iota
	fileView
	searchView
	actionView
	settingsView
	zipActionView
)



// item implements list.Item for our files
type item struct {
	node *Node
}

func (i item) Title() string { 
	if i.node.metadata.IsDir {
		return "▸ " + i.node.metadata.Name
	}
	return "  " + i.node.metadata.Name 
}

func (i item) Description() string {
	size := formatSize(i.node.metadata.Size)
	if i.node.metadata.IsDir {
		size = "Directory"
	}
	modTime := i.node.metadata.ModTime.Format("Jan 02 15:04")
	return fmt.Sprintf("%s • %s", size, modTime)
}

func (i item) FilterValue() string { return i.node.metadata.Name }

// actionItem for the action menu
type actionItem struct {
	title, desc string
	actionID    string
}

func (i actionItem) Title() string       { return i.title }
func (i actionItem) Description() string { return i.desc }
func (i actionItem) FilterValue() string { return i.title }

// settingItem for settings
type settingItem struct {
	feature Feature
}

func (i settingItem) Title() string {
	status := "○"
	switch i.feature.Status {
	case "done":
		status = "●"
	case "in-progress":
		status = "◐"
	}
	return fmt.Sprintf("%s %s", status, i.feature.Name)
}
func (i settingItem) Description() string { return i.feature.Description }
func (i settingItem) FilterValue() string { return i.feature.Name }


type fileModel struct {
	list list.Model
}

type searchModel struct {
	input textinput.Model
	list  list.Model
}

type actionModel struct {
	list list.Model
}

type settingsModel struct {
	list list.Model
}

type zipModel struct {
	input textinput.Model
	chosenPath string
}

type model struct {
	currentView View
	engine      *Engine
	compressingEngine *CompressEngine
	views Stack[View]
	file    fileModel
	search  searchModel
	actions actionModel
	settings settingsModel
	zip     zipModel

	width, height int
}

func NewModel() model {
	// Initialize Engine
	
	engine := NewEngine(dir);

	// File List
	fileList := list.New(nodesToItems(engine.current.children), list.NewDefaultDelegate(), 0, 0)
	fileList.Title = "File Explorer"
	fileList.SetShowHelp(false)

	// Search
	ti := textinput.New()
	ti.Placeholder = "Search files..."
	ti.CharLimit = 156
	ti.Width = 20

	searchList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	searchList.Title = "Search Results"
	searchList.SetShowHelp(false)

	// Actions
	actionItems := []list.Item{
		actionItem{title: "Rename", desc: "Rename the selected file", actionID: "rename"},
		actionItem{title: "Delete", desc: "Delete the selected file", actionID: "delete"},
		actionItem{title: "Copy Path", desc: "Copy full path to clipboard", actionID: "copypath"},
		actionItem{title: "Properties", desc: "Show file properties", actionID: "props"},
		actionItem{title: "Compress to Zip", desc: "Create a zip archive", actionID: "zip"},
	}
	actionList := list.New(actionItems, list.NewDefaultDelegate(), 0, 0)
	actionList.Title = "Actions"
	actionList.SetShowHelp(false)

	// Settings
	// Initialize features list
	features := []Feature{
		{Name: "Open Files", Description: "Open files with default Windows app", Status: "todo", Priority: "high"},
		{Name: "Copy/Paste", Description: "Ctrl+C, Ctrl+V file operations", Status: "todo", Priority: "high"},
		{Name: "Fuzzy Search", Description: "Fast fuzzy file matching", Status: "todo", Priority: "high"},
		{Name: "Multi-Select", Description: "Space to select, bulk operations", Status: "todo", Priority: "high"},
		{Name: "Sort Options", Description: "Sort by name/size/date/type", Status: "todo", Priority: "medium"},
	}
	
	settingsItems := make([]list.Item, len(features))
	for i, f := range features {
		settingsItems[i] = settingItem{feature: f}
	}
	settingsList := list.New(settingsItems, list.NewDefaultDelegate(), 0, 0)
	settingsList.Title = "Roadmap / Settings"

	// Zip
	zipInput := textinput.New()
	zipInput.Placeholder = "archive.zip"

	return model{
		currentView: titleView,
		engine:      engine,
		compressingEngine: NewCompressEngine(4),
		file:        fileModel{list: fileList},
		search:      searchModel{input: ti, list: searchList},
		actions:     actionModel{list: actionList},
		settings:    settingsModel{list: settingsList},
		zip:         zipModel{input: zipInput},
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Resize lists
		h, v := docStyle.GetFrameSize()
		m.file.list.SetSize(msg.Width-h, msg.Height-v)
		m.search.list.SetSize(msg.Width-h, msg.Height-v-4) // -4 for input height roughly
		m.actions.list.SetSize(msg.Width-h, msg.Height-v)
		m.settings.list.SetSize(msg.Width-h, msg.Height-v)
	}

	switch m.currentView {
	case titleView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter":
				m.views.Push(m.currentView)
				m.currentView = fileView
				return m, nil
			case "s":
				m.views.Push(m.currentView)
				m.currentView = searchView
				m.search.input.Focus()
				return m, textinput.Blink
			case "?":
				m.views.Push(m.currentView)
				m.currentView = settingsView
				return m, nil
			}
		}
	case fileView:
		m.file, cmd = m.updateFileView(msg)
		cmds = append(cmds, cmd)
	case searchView:
		m.search, cmd = m.updateSearchView(msg)
		cmds = append(cmds, cmd)
	case actionView:
		// Navigation handled here for view switching
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" {
				view,poss := m.views.Pop();
				if(poss==true){
					m.currentView=view;
					return m, nil
				}
			}
			if msg.String() == "enter" {
				// Handle Action
				selectedItem := m.actions.list.SelectedItem()
				if selectedItem != nil {
					act := selectedItem.(actionItem)
					return m.handleAction(act)
				}
			}
		}
		newActionList, newCmd := m.actions.list.Update(msg)
		m.actions.list = newActionList
		cmds = append(cmds, newCmd)

	case settingsView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" {
				view,poss := m.views.Pop();
				if(poss==true){
					m.currentView=view;
					return m, nil
				}
			}
		}
		newSettingsList, newCmd := m.settings.list.Update(msg)
		m.settings.list = newSettingsList
		cmds = append(cmds, newCmd)

	case zipActionView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" {
				view,poss := m.views.Pop();
				if(poss==true){
					m.currentView=view;
					return m, nil
				}
			}
			if msg.String() == "enter" {
				// Execute Zip
				err := m.compressingEngine.CompressFileZip(m.zip.chosenPath, m.zip.input.Value())
				if err != nil {
					// In real app, show error msg
				}
				m.currentView = fileView
				return m, nil
			}
		}
		newInput, newCmd := m.zip.input.Update(msg)
		m.zip.input = newInput
		cmds = append(cmds, newCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) updateFileView(msg tea.Msg) (fileModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m.file, tea.Quit
		case "s":
			m.views.Push(m.currentView)
			m.currentView = searchView
			m.search.input.Focus()
			return m.file, textinput.Blink
		case "a":
			m.views.Push(m.currentView)
			m.currentView = actionView
			return m.file, nil
		case "?":
			m.views.Push(m.currentView)
			m.currentView = settingsView
			return m.file, nil
		case "enter":
			// Navigate into directory
			selected := m.file.list.SelectedItem()
			if selected != nil {
				itm := selected.(item)
				if itm.node.metadata.IsDir {
					loadChildren(itm.node)
					m.engine.ChangeDirectory(itm.node)
					// Update list items
					cmd = m.file.list.SetItems(nodesToItems(m.engine.current.children))
					m.file.list.ResetSelected()
				}
			}
			return m.file, cmd
		case "backspace", "left":
			// Go up
			if m.engine.current.parent != nil {
				m.engine.ChangeDirectory(m.engine.current.parent)
				cmd = m.file.list.SetItems(nodesToItems(m.engine.current.children))
				m.file.list.ResetSelected()
			}
			return m.file, cmd
		case "esc":
			view,poss := m.views.Pop();
			if(poss==true){
				m.currentView=view;
			}
			return m.file, cmd;
		}
	}
	
	newFileModel, newCmd := m.file.list.Update(msg)
	m.file.list = newFileModel
	return m.file, newCmd
}

func (m *model) updateSearchView(msg tea.Msg) (searchModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			view,poss := m.views.Pop();
			if(poss==true){
				m.currentView=view;
				m.search.input.Blur()
			}
			return m.search, nil
		case "enter":
			if m.search.input.Focused() {
				// Perform search
				query := m.search.input.Value()
				results, _ := m.engine.Search(query)
				m.search.list.SetItems(nodesToItems(results))
				m.search.input.Blur()
			} else {
				// Navigate to result
				selected := m.search.list.SelectedItem()
				if selected != nil {
					itm := selected.(item)
					if itm.node.metadata.IsDir {
						loadChildren(itm.node)
						m.engine.ChangeDirectory(itm.node)
						m.file.list.SetItems(nodesToItems(m.engine.current.children))
						m.views.Push(m.currentView)
						m.currentView = fileView
					}
				}
			}
			return m.search, nil
		case "tab":
			if m.search.input.Focused() {
				m.search.input.Blur()
			} else {
				m.search.input.Focus()
			}
			return m.search, nil
		}
	}

	if m.search.input.Focused() {
		m.search.input, cmd = m.search.input.Update(msg)
	} else {
		m.search.list, cmd = m.search.list.Update(msg)
	}
	return m.search, cmd
}

func (m *model) handleAction(act actionItem) (tea.Model, tea.Cmd) {
	selected := m.file.list.SelectedItem()
	if selected == nil {
		m.views.Push(m.currentView)
		m.currentView = fileView
		return m, nil
	}
	fileItem := selected.(item)

	switch act.actionID {
	case "zip":
		m.zip.chosenPath = fileItem.node.metadata.Path
		m.views.Push(m.currentView)
		m.currentView = zipActionView
		m.zip.input.Focus()
		return m, textinput.Blink
	case "delete":
		// Implement delete logic interaction or command
	}
	m.views.Push(m.currentView)
	m.currentView = fileView
	return m, nil
}


func (m model) View() string {
	switch m.currentView {
	case titleView:
		return m.renderTitleView()
	case fileView:
		return docStyle.Render(m.file.list.View())
	case searchView:
		return docStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left, 
				m.search.input.View(),
				"We are currently in " + m.engine.root.metadata.Path,
				m.search.list.View(),
			),
		)
	case actionView:
		return docStyle.Render(m.actions.list.View())
	case settingsView:
		return docStyle.Render(m.settings.list.View())
	case zipActionView:
		return docStyle.Render(fmt.Sprintf(
			"Compressing %s\n\nEnter output filename:\n%s",
			m.zip.chosenPath,
			m.zip.input.View(),
		))
	}
	return "Unknown View"
}


func nodesToItems(nodes []*Node) []list.Item {
	items := make([]list.Item, len(nodes))
	for i, n := range nodes {
		items[i] = item{node: n}
	}
	return items
}

// Feature represents a potential feature to implement
type Feature struct {
	Name        string
	Description string
	Status      string // "todo", "in-progress", "done"
	Priority    string // "high", "medium", "low"
}

// formatSize converts bytes to human readable format
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func (m model) renderTitleView() string {
	logo := []string{
	"  ______ _ _      _____  _                     _       ",
	" |  ____(_) |    |  __ \\| |                   | |      ",
	" | |__   _| | ___| |  | | |__  _   _ _ __   __| | ___  ",
	" |  __| | | |/ _ \\ |  | | '_ \\| | | | '_ \\ / _` |/ _ \\ ",
	" | |    | | |  __/ |__| | | | | |_| | | | | (_| | (_) |",
	" |_|    |_|_|\\___|_____/|_| |_|\\__,_|_| |_|\\__,_|\\___/ ",
}

	var s string
	s += "\n"
	for _, line := range logo {
		s += titleLogoStyle.Render(line) + "\n"
	}
	s += "\n"

	// Tagline
	s += titleMutedStyle.Render("  Terminal File Browser") + titlePathStyle.Render(" v0.1.0") + "\n\n"

	// Quick Actions
	s += titleMutedStyle.Render("  Quick Actions") + "\n"
	s += titleDividerStyle.Render("  ────────────────────────────────────────") + "\n\n"

	actions := []struct {
		key  string
		desc string
	}{
		{"enter", "Browse files"},
		{"s", "Search files"},
		{"?", "Settings / Roadmap"},
		{"q", "Quit"},
	}

	for _, a := range actions {
		s += fmt.Sprintf("    %s %s\n", titleAccentStyle.Render(fmt.Sprintf("%-8s", a.key)), a.desc)
	}

	s += "\n"

	// Current directory
	s += titleDividerStyle.Render("  ────────────────────────────────────────") + "\n"
	s += titleMutedStyle.Render("   ") + titlePathStyle.Render(m.engine.current.metadata.Path) + "\n"

	// Hints
	s += "\n"
	s += titleAccentStyle.Render("  enter") + titleMutedStyle.Render(" browse  ")
	s += titleAccentStyle.Render("s") + titleMutedStyle.Render(" search  ")
	s += titleAccentStyle.Render("?") + titleMutedStyle.Render(" settings  ")
	s += titleAccentStyle.Render("q") + titleMutedStyle.Render(" quit") + "\n"

	return s
}
