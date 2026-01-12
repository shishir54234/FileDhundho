package main

import (
	"github.com/charmbracelet/lipgloss"
)

// NeoVim-inspired palette
var (
	colorGreen      = lipgloss.Color("#98c379") // NeoVim green
	colorBlue       = lipgloss.Color("#61afef") // Soft blue
	colorYellow     = lipgloss.Color("#e5c07b") // Warm yellow
	colorOrange     = lipgloss.Color("#d19a66") // Orange accent
	colorPurple     = lipgloss.Color("#c678dd") // Purple
	colorCyan       = lipgloss.Color("#56b6c2") // Cyan
	colorRed        = lipgloss.Color("#e06c75") // Red
	colorGray       = lipgloss.Color("#5c6370") // Comment gray
	colorDimGray    = lipgloss.Color("#3e4451") // Dimmer gray
	colorLightGray  = lipgloss.Color("#abb2bf") // Light gray text
	colorWhite      = lipgloss.Color("#dcdfe4") // Off-white text
	colorBgDark     = lipgloss.Color("#282c34") // Dark background
	colorBgSelected = lipgloss.Color("#3e4451") // Selected bg
)

var (
	headerStyle = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorGray).
			Italic(true)

	glowStyle = lipgloss.NewStyle().
			Foreground(colorDimGray)

	listStyle = lipgloss.NewStyle().
			Padding(0, 1)

	itemStyle = lipgloss.NewStyle().
			Foreground(colorWhite)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(colorWhite).
				Background(colorBgSelected).
				Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(colorDimGray).
			MarginTop(1)

	searchBarStyle = lipgloss.NewStyle().
			Foreground(colorWhite).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorDimGray).
			Padding(0, 1).
			Width(60)

	searchIconStyle = lipgloss.NewStyle().
			Foreground(colorGreen)

	placeholderStyle = lipgloss.NewStyle().
				Foreground(colorDimGray).
				Italic(true)

	focusedStyle = lipgloss.NewStyle().
			Foreground(colorWhite).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorGreen).
			Padding(0, 1).
			Width(60)

	resultsStyle = lipgloss.NewStyle().
			Foreground(colorWhite).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorDimGray).
			Padding(1, 1).
			MarginTop(1).
			Width(60)

	resultSelectedStyle = lipgloss.NewStyle().
				Foreground(colorGreen).
				Bold(true)

	resultCountStyle = lipgloss.NewStyle().
				Foreground(colorGray).
				MarginTop(1)

	scrollIndicatorStyle = lipgloss.NewStyle().
				Foreground(colorDimGray)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorGray).
			MarginTop(1)

	pathStyle = lipgloss.NewStyle().
			Foreground(colorBlue)

	folderIconStyle = lipgloss.NewStyle().
			Foreground(colorYellow)

	fileIconStyle = lipgloss.NewStyle().
			Foreground(colorLightGray)

	dividerStyle = lipgloss.NewStyle().
			Foreground(colorDimGray)

	// NeoVim-style elements
	accentStyle = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true)

	mutedStyle = lipgloss.NewStyle().
			Foreground(colorGray)

	successStyle = lipgloss.NewStyle().
			Foreground(colorGreen)

	// Logo style for ASCII art
	logoStyle = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true)

	// Settings view styles
	warningStyle = lipgloss.NewStyle().
			Foreground(colorYellow)

	highPriorityStyle = lipgloss.NewStyle().
				Foreground(colorRed).
				Bold(true)

	mediumPriorityStyle = lipgloss.NewStyle().
				Foreground(colorYellow)

	lowPriorityStyle = lipgloss.NewStyle().
				Foreground(colorGray)
)
var docStyle = lipgloss.NewStyle().Margin(1, 2)
var (
	titleLogoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#98c379")).
			Bold(true)

	titleAccentStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#98c379")).
				Bold(true)

	titleMutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5c6370"))

	titlePathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#61afef"))

	titleDividerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#3e4451"))
)