/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var titleStyle = func() lipgloss.Style {
	b := lipgloss.RoundedBorder()
	b.Right = "├"
	return lipgloss.NewStyle().Border(b).Padding(0, 1)
}()

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "View and filter Android and iOS logs",
	Long: `The 'log' command allows you to stream, view, and filter logs from Android and iOS devices or simulators/emulators.

You can use this tool to debug applications by capturing device logs in real-time, applying filters based on keywords, log levels, platforms, or other criteria.

Examples:
  log                             # View all logs
  log --platform android          # Show only Android logs
  log --filter "Error"            # Filter logs containing the word 'Error'
  log --platform ios --level debug # Show iOS logs at debug level

Supports both physical and virtual devices.
Useful for mobile developers and testers needing better insight into runtime behavior.`,
	Run: func(cmd *cobra.Command, args []string) {

		platform, err := cmd.Flags().GetString("platform")

		if err != nil {
			fmt.Println("You must provide a platform to watch logs", err)
			return
		}

		switch platform {
		case "android":
			now := time.Now().Format("01-02 15:04:05.000")
			logCmd := exec.Command("adb", "logcat", "-T", now)
			watchLogs("Android", logCmd)
		case "ios":
			logCmd := exec.Command("xcrun", "simctl", "spawn", "booted", "log", "stream", "--style", "syslog")
			watchLogs("IOS", logCmd)
		}

	},
}

type model struct {
	selectedPlatform string
	viewport         viewport.Model
	textInput        textinput.Model
	logs             []string
	filterText       string
	logChan          chan string
	err              error
	focused          string
	paused           bool
}

func (m model) helpView() string {
	return "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(
		"j/k: scroll • /: focus filter • esc: blur filter • r: resume • c: clear • q: quit",
	)
}

func initialModel(platform string) model {
	ti := textinput.New()
	ti.Placeholder = "Type to filter logs..."

	return model{
		selectedPlatform: platform,
		viewport:         viewport.New(0, 0),
		textInput:        ti,
		logs:             []string{},
		logChan:          make(chan string),
		focused:          "viewport", // New field to track focus
	}
}

func (m model) header() string {
	title := titleStyle.Render(m.selectedPlatform)
	line := strings.Repeat("-", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		waitForLogs(m.logChan),
	)
}

func highlightText(text, substr string) string {
	if substr == "" {
		return text
	}
	highlightStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("205")).
		Foreground(lipgloss.Color("0"))

	return strings.ReplaceAll(
		text,
		substr,
		highlightStyle.Render(substr),
	)
}

func wrapLine(line string, width int) []string {
	var wrapped []string
	for len(line) > width {
		wrapped = append(wrapped, line[:width])
		line = line[width:]
	}
	if len(line) > 0 {
		wrapped = append(wrapped, line)
	}
	return wrapped
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "/":
			if m.focused == "viewport" {
				m.focused = "filter"
				m.textInput.Focus()
			}
		case "esc":
			if m.focused == "filter" {
				m.focused = "viewport"
				m.textInput.Blur()
			}
		case "r":
			m.paused = false
			m.viewport.GotoBottom()
		}

		if m.focused == "viewport" {
			switch msg.String() {
			case "up", "k":
				m.paused = true
				m.viewport.ScrollUp(1)
			case "down", "j":
				m.viewport.ScrollDown(1)
				if m.viewport.AtBottom() {
					m.paused = false
				}
			case "pgup", "ctrl+b":
				m.paused = true
				m.viewport.HalfPageUp()
			case "pgdown", "ctrl+f":
				m.viewport.HalfPageDown()
				if m.viewport.AtBottom() {
					m.paused = false
				}
			case "home", "g":
				m.paused = true
				m.viewport.GotoTop()
			case "end", "G":
				m.paused = false
				m.viewport.GotoBottom()
			case "c":
				m.viewport.SetContent("")
				m.logs = []string{}
			}
		}

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 5 // Adjusted for help view
		m.textInput.Width = msg.Width - 4

	case logMsg:
		m.logs = append(m.logs, string(msg))
		m.updateViewport()
		return m, waitForLogs(m.logChan)

	case error:
		m.err = msg
		return m, nil
	}

	if m.focused == "filter" {
		m.textInput, cmd = m.textInput.Update(msg)
		m.filterText = m.textInput.Value()
		m.updateViewport()
	}

	return m, cmd
}

func (m *model) updateViewport() {
	filteredLogs := []string{}
	for _, log := range m.logs {
		if m.filterText == "" || strings.Contains(strings.ToLower(log), strings.ToLower(m.filterText)) {
			highlightedLog := highlightText(log, m.filterText)
			wrappedLines := wrapLine(highlightedLog, m.viewport.Width)
			filteredLogs = append(filteredLogs, wrappedLines...)
		}
	}
	m.viewport.SetContent(strings.Join(filteredLogs, "\n"))
	if !m.paused {
		m.viewport.GotoBottom()
	}
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n%s\n%s%s",
		m.header(),
		m.viewport.View(),
		m.textInput.View(),
		m.helpView(),
	)
}

type logMsg string

func waitForLogs(logChan chan string) tea.Cmd {
	return func() tea.Msg {
		return logMsg(<-logChan)
	}
}

func watchLogs(platform string, execCmd *exec.Cmd) {
	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating StdoutPipe", err)
		return
	}

	if err := execCmd.Start(); err != nil {
		fmt.Println("Error starting command", err)
		return
	}

	m := initialModel(platform)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			m.logChan <- scanner.Text()
		}
	}()

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		return
	}

	if err := execCmd.Wait(); err != nil {
		fmt.Println("Error waiting for command", err)
	}
}

func init() {
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	logCmd.Flags().StringP("platform", "p", "", "Choose android or ios")
}
