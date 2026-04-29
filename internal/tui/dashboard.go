package tui

import (
	"fmt"
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MetricsData holds the data displayed in the dashboard.
type MetricsData struct {
	QPS            float64
	AvgLatency     float64
	P99Latency     float64
	TotalRequests  int64
	Status2xx      int64
	Status4xx      int64
	Status5xx      int64
	InFlight       int
	Upstreams      []UpstreamStatus
	CircuitStates  map[string]string
}

// UpstreamStatus represents an upstream server's health.
type UpstreamStatus struct {
	Address string
	Healthy bool
	Latency float64
}

// Dashboard is the TUI model.
type Dashboard struct {
	data     MetricsData
	quitting bool
	width    int
	height   int
}

// NewDashboard creates a new TUI dashboard.
func NewDashboard() *Dashboard {
	return &Dashboard{
		width:  80,
		height: 24,
		data: MetricsData{
			CircuitStates: make(map[string]string),
		},
	}
}

func (d *Dashboard) Init() tea.Cmd {
	return tea.Batch(
		tick(),
		tea.WindowSize(),
	)
}

func (d *Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height
		return d, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			d.quitting = true
			return d, tea.Quit
		}

	case tickMsg:
		// In production, fetch from Prometheus /metrics endpoint
		// For now, simulate data
		d.data = d.fetchMetrics()
		return d, tick()
	}

	return d, nil
}

func (d *Dashboard) View() string {
	if d.quitting {
		return "Dashboard closed.\n"
	}

	var b strings.Builder

	// Header
	header := titleStyle.Render("  AEGIS API GATEWAY - MONITOR  ")
	b.WriteString(header)
	b.WriteString("\n\n")

	// Stats cards
	qpsCard := statCard("QPS", fmt.Sprintf("%.1f", d.data.QPS))
	latencyCard := statCard("Avg Latency", fmt.Sprintf("%.1fms", d.data.AvgLatency))
	p99Card := statCard("P99 Latency", fmt.Sprintf("%.1fms", d.data.P99Latency))
	totalCard := statCard("Total Requests", fmt.Sprintf("%d", d.data.TotalRequests))
	inflightCard := statCard("In Flight", fmt.Sprintf("%d", d.data.InFlight))

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, qpsCard, latencyCard, p99Card, totalCard, inflightCard))
	b.WriteString("\n\n")

	// Status codes
	b.WriteString(sectionTitle.Render("  HTTP Status Codes"))
	b.WriteString("\n")
	s2xx := statusBar("2xx", d.data.Status2xx, d.data.TotalRequests, successColor)
	s4xx := statusBar("4xx", d.data.Status4xx, d.data.TotalRequests, warnColor)
	s5xx := statusBar("5xx", d.data.Status5xx, d.data.TotalRequests, errorColor)
	b.WriteString(s2xx + "\n")
	b.WriteString(s4xx + "\n")
	b.WriteString(s5xx + "\n\n")

	// Upstreams
	b.WriteString(sectionTitle.Render("  Upstream Nodes"))
	b.WriteString("\n")
	for _, u := range d.data.Upstreams {
		b.WriteString(upstreamRow(u) + "\n")
	}
	if len(d.data.Upstreams) == 0 {
		b.WriteString(dimStyle.Render("  No upstreams configured") + "\n")
	}
	b.WriteString("\n")

	// Circuit breakers
	b.WriteString(sectionTitle.Render("  Circuit Breakers"))
	b.WriteString("\n")
	for route, state := range d.data.CircuitStates {
		b.WriteString(circuitRow(route, state) + "\n")
	}
	if len(d.data.CircuitStates) == 0 {
		b.WriteString(dimStyle.Render("  No circuit breakers active") + "\n")
	}
	b.WriteString("\n")

	// Footer
	b.WriteString(dimStyle.Render("  Press 'q' to quit"))

	return b.String()
}

func (d *Dashboard) fetchMetrics() MetricsData {
	// Simulate metrics for demo
	return MetricsData{
		QPS:           1250.5,
		AvgLatency:    12.3,
		P99Latency:    45.8,
		TotalRequests: d.data.TotalRequests + 1250,
		Status2xx:     d.data.Status2xx + 1200,
		Status4xx:     d.data.Status4xx + 40,
		Status5xx:     d.data.Status5xx + 10,
		InFlight:      3,
		Upstreams: []UpstreamStatus{
			{Address: "http://localhost:9001", Healthy: true, Latency: 8.5},
			{Address: "http://localhost:9002", Healthy: true, Latency: 12.1},
			{Address: "http://localhost:9003", Healthy: false, Latency: 0},
		},
		CircuitStates: map[string]string{
			"/api/users":          "closed",
			"/api/users/{id}":     "closed",
			"/api/files/*":        "open",
		},
	}
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Bold(true).
			Padding(0, 2)

	sectionTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("62")).
			Bold(true)

	statLabel = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	statValue = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Bold(true)

	successColor = lipgloss.Color("42")
	warnColor    = lipgloss.Color("214")
	errorColor   = lipgloss.Color("196")
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

func statCard(label, value string) string {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 2).
		Width(18).
		Align(lipgloss.Center).
		Render(statLabel.Render(label) + "\n" + statValue.Render(value))
}

func statusBar(label string, count, total int64, color lipgloss.TerminalColor) string {
	pct := float64(0)
	if total > 0 {
		pct = float64(count) / float64(total) * 100
	}
	barWidth := 30
	filled := int(math.Round(pct / 100 * float64(barWidth)))
	if filled > barWidth {
		filled = barWidth
	}

	bar := lipgloss.NewStyle().Foreground(color).Render(strings.Repeat("█", filled))
	empty := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(strings.Repeat("░", barWidth-filled))

	return fmt.Sprintf("  %-4s %s%s %6d (%5.1f%%)", label, bar, empty, count, pct)
}

func upstreamRow(u UpstreamStatus) string {
	status := lipgloss.NewStyle().Foreground(successColor).Render("● UP")
	if !u.Healthy {
		status = lipgloss.NewStyle().Foreground(errorColor).Render("● DOWN")
	}
	lat := fmt.Sprintf("%.1fms", u.Latency)
	if !u.Healthy {
		lat = "-"
	}
	return fmt.Sprintf("  %-30s %s  %s", u.Address, status, lat)
}

func circuitRow(route, state string) string {
	var color lipgloss.TerminalColor
	switch state {
	case "closed":
		color = successColor
	case "open":
		color = errorColor
	case "half-open":
		color = warnColor
	}
	stateStr := lipgloss.NewStyle().Foreground(color).Render(strings.ToUpper(state))
	return fmt.Sprintf("  %-30s %s", route, stateStr)
}
