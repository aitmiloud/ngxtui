package styles

import "github.com/charmbracelet/lipgloss"

// Color palette - modern dark theme with vibrant accents
var (
	// Primary colors
	AccentPrimary   = lipgloss.Color("#00E5FF") // Bright cyan
	AccentSecondary = lipgloss.Color("#8B5CF6") // Vivid purple
	AccentSuccess   = lipgloss.Color("#10B981") // Emerald green
	AccentWarning   = lipgloss.Color("#FBBF24") // Amber
	AccentDanger    = lipgloss.Color("#F43F5E") // Rose red
	AccentInfo      = lipgloss.Color("#3B82F6") // Blue

	// Text colors
	TextPrimary   = lipgloss.Color("#F8FAFC") // Almost white
	TextSecondary = lipgloss.Color("#CBD5E1") // Light slate
	TextMuted     = lipgloss.Color("#94A3B8") // Slate
	TextDim       = lipgloss.Color("#64748B") // Dark slate

	// Background colors
	BgPrimary   = lipgloss.Color("#0F172A") // Deep slate
	BgSecondary = lipgloss.Color("#1E293B") // Slate 800
	BgTertiary  = lipgloss.Color("#334155") // Slate 700
	BgHighlight = lipgloss.Color("#475569") // Slate 600

	// Border colors
	BorderSubtle  = lipgloss.Color("#475569") // Slate 600
	BorderDefault = lipgloss.Color("#64748B") // Slate 500
	BorderBright  = lipgloss.Color("#94A3B8") // Slate 400

	// Gradient colors for charts
	GradientStart = lipgloss.Color("#8B5CF6") // Purple
	GradientMid   = lipgloss.Color("#3B82F6") // Blue
	GradientEnd   = lipgloss.Color("#00E5FF") // Cyan
)

// Base styles
var (
	Base = lipgloss.NewStyle().
		Foreground(TextPrimary).
		Background(BgPrimary)

	// Header with gradient-like effect
	Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(TextPrimary).
		Background(lipgloss.Color("#1E293B")).
		Padding(1, 3).
		MarginBottom(1).
		Border(lipgloss.Border{
			Bottom: "─",
		}).
		BorderForeground(AccentPrimary)

	// Tab styles with modern look
	ActiveTab = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("14")). // Bright cyan
			Padding(0, 4).
			MarginRight(1).
			Border(lipgloss.Border{
			Bottom: "▔",
		}).
		BorderForeground(lipgloss.Color("14"))

	InactiveTab = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")). // Bright yellow - guaranteed visible
			Padding(0, 4).
			MarginRight(1)

	// Enhanced panel with subtle shadow effect
	Panel = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderSubtle).
		Padding(1, 2).
		MarginBottom(1)

	PanelTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(AccentPrimary).
			MarginBottom(1).
			Underline(true)

	// Modern badge styles with icons
	EnabledBadge = lipgloss.NewStyle().
			Foreground(AccentSuccess).
			Padding(0, 2).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(AccentSuccess)

	DisabledBadge = lipgloss.NewStyle().
			Foreground(TextMuted).
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderSubtle)

	// Action menu styles
	ActionItem = lipgloss.NewStyle().
			Foreground(TextSecondary).
			Padding(0, 2).
			MarginLeft(2)

	SelectedAction = lipgloss.NewStyle().
			Bold(true).
			Foreground(AccentPrimary).
			Padding(0, 3).
			MarginLeft(1).
			Border(lipgloss.Border{
			Left: "▐",
		}).
		BorderForeground(AccentPrimary)

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

	// Card styles for metrics and stats
	Card = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderDefault).
		Padding(1, 2).
		Margin(0, 1).
		Background(BgSecondary)

	CardTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(AccentPrimary).
			MarginBottom(1)

	CardValue = lipgloss.NewStyle().
			Bold(true).
			Foreground(TextPrimary).
			Align(lipgloss.Center)

	CardSubtext = lipgloss.NewStyle().
			Foreground(TextMuted).
			Align(lipgloss.Center)

	// Status indicator styles
	StatusSuccess = lipgloss.NewStyle().
			Foreground(AccentSuccess).
			Bold(true)

	StatusError = lipgloss.NewStyle().
			Foreground(AccentDanger).
			Bold(true)

	StatusWarning = lipgloss.NewStyle().
			Foreground(AccentWarning).
			Bold(true)

	StatusInfo = lipgloss.NewStyle().
			Foreground(AccentInfo).
			Bold(true)

	// Table header style
	TableHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(AccentPrimary).
			Padding(0, 1)

	// Divider styles
	Divider = lipgloss.NewStyle().
		Foreground(BorderSubtle)

	DividerBright = lipgloss.NewStyle().
			Foreground(BorderBright)

	// Footer style
	Footer = lipgloss.NewStyle().
		Foreground(TextMuted).
		Padding(0, 2).
		MarginTop(1).
		Border(lipgloss.Border{
			Top: "─",
		}).
		BorderForeground(BorderSubtle)
)
