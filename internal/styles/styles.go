package styles

import "github.com/charmbracelet/lipgloss"

// Color palette - dark theme
var (
	AccentPrimary   = lipgloss.Color("#00D9FF")
	AccentSecondary = lipgloss.Color("#7C3AED")
	AccentSuccess   = lipgloss.Color("#10B981")
	AccentWarning   = lipgloss.Color("#F59E0B")
	AccentDanger    = lipgloss.Color("#EF4444")
	TextPrimary     = lipgloss.Color("#F9FAFB")
	TextSecondary   = lipgloss.Color("#9CA3AF")
	TextMuted       = lipgloss.Color("#6B7280")
	BgPrimary       = lipgloss.Color("#111827")
	BgSecondary     = lipgloss.Color("#1F2937")
	BgTertiary      = lipgloss.Color("#374151")
	BorderSubtle    = lipgloss.Color("#4B5563")
)

// Base styles
var (
	Base = lipgloss.NewStyle().
		Foreground(TextPrimary)

	Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(AccentPrimary).
		Background(BgSecondary).
		Padding(0, 2).
		MarginBottom(1)

	ActiveTab = lipgloss.NewStyle().
			Bold(true).
			Foreground(TextPrimary).
			Background(AccentSecondary).
			Padding(0, 3).
			MarginRight(1)

	InactiveTab = lipgloss.NewStyle().
			Foreground(TextSecondary).
			Background(BgTertiary).
			Padding(0, 3).
			MarginRight(1)

	Panel = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderSubtle).
		Padding(1, 2).
		MarginBottom(1)

	PanelTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(AccentPrimary).
			MarginBottom(1)

	EnabledBadge = lipgloss.NewStyle().
			Foreground(TextPrimary).
			Background(AccentSuccess).
			Padding(0, 1).
			Bold(true)

	DisabledBadge = lipgloss.NewStyle().
			Foreground(TextPrimary).
			Background(TextMuted).
			Padding(0, 1)

	ActionItem = lipgloss.NewStyle().
			Foreground(TextSecondary).
			Padding(0, 2).
			MarginBottom(0)

	SelectedAction = lipgloss.NewStyle().
			Bold(true).
			Foreground(TextPrimary).
			Background(AccentSecondary).
			Padding(0, 2).
			MarginBottom(0)

	HelpKey = lipgloss.NewStyle().
		Foreground(AccentPrimary).
		Bold(true)

	HelpDesc = lipgloss.NewStyle().
			Foreground(TextSecondary)

	HelpSeparator = lipgloss.NewStyle().
			Foreground(TextMuted)

	ErrorText = lipgloss.NewStyle().
			Foreground(AccentDanger).
			Bold(true)

	SuccessText = lipgloss.NewStyle().
			Foreground(AccentSuccess).
			Bold(true)

	WarningText = lipgloss.NewStyle().
			Foreground(AccentWarning).
			Bold(true)

	InfoText = lipgloss.NewStyle().
			Foreground(AccentPrimary)

	MutedText = lipgloss.NewStyle().
			Foreground(TextMuted)
)
