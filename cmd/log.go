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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		platform, err := cmd.Flags().GetString("platform")

		if err != nil {
			fmt.Println("You must provide a platform to watch logs", err)
			return
		}

		switch platform {
		case "android":
			watchAndroidLogs()
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
}

func initialModel(platform string) model {
	ti := textinput.New()
	ti.Placeholder = "Type to filter logs..."
	ti.Focus()

	return model{
		selectedPlatform: platform,
		viewport:         viewport.New(0, 0),
		textInput:        ti,
		logs:             []string{},
		logChan:          make(chan string),
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
	highlightStart := "\033[1;33m" // Bold yellow
	highlightEnd := "\033[0m"
	return strings.Replace(text, substr, highlightStart+substr+highlightEnd, -1)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 3
		m.textInput.Width = msg.Width - 4

	case logMsg:
		m.logs = append(m.logs, string(msg))
		m.updateViewport()
		return m, waitForLogs(m.logChan)

	case error:
		m.err = msg
		return m, nil
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	m.filterText = m.textInput.Value()
	m.updateViewport()

	return m, tea.Batch(cmds...)
}

func (m *model) updateViewport() {
	filteredLogs := []string{}
	for _, log := range m.logs {
		if m.filterText == "" || strings.Contains(strings.ToLower(log), strings.ToLower(m.filterText)) {
			highlightedLog := highlightText(log, m.filterText)
			filteredLogs = append(filteredLogs, highlightedLog)
		}
	}
	m.viewport.SetContent(strings.Join(filteredLogs, "\n"))
	m.viewport.GotoBottom()
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n%s\n%s",
		m.header(),
		m.viewport.View(),
		m.textInput.View(),
	)
}

type logMsg string

func waitForLogs(logChan chan string) tea.Cmd {
	return func() tea.Msg {
		return logMsg(<-logChan)
	}
}

func watchAndroidLogs() {
	now := time.Now().Format("01-02 15:04:05.000")
	execCmd := exec.Command("adb", "logcat", "-T", now)
	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating StdoutPipe", err)
		return
	}

	if err := execCmd.Start(); err != nil {
		fmt.Println("Error starting command", err)
		return
	}

	m := initialModel("android")

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
