package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gen2brain/beeep"
)

type phase int

const (
	focusPhase phase = iota
	restPhase
	finishedPhase
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#EF4444")).
			Padding(0, 1).
			MarginBottom(1)

	focusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")).
			Bold(true)

	restStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#737373")).
			MarginTop(1)
)

type model struct {
	phase       phase
	repetition  int
	totalReps   int
	focusDur    time.Duration
	restDur     time.Duration
	timer       timer.Model
	progress    progress.Model
	quitting    bool
	interrupted bool
}

func (m model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case " ":
			return m, m.timer.Toggle()
		case "r":
			return m, m.timer.Init()
		}

	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		return m.nextPhase()
	}

	return m, nil
}

func (m model) nextPhase() (tea.Model, tea.Cmd) {
	switch m.phase {
	case focusPhase:
		if m.restDur > 0 {
			m.phase = restPhase
			m.timer = timer.New(m.restDur)
			return m, tea.Batch(alertCmd("You can rest now"), m.timer.Init())
		}
		return m.handleFocusEnd(alertCmd("You can rest now"))
	case restPhase:
		return m.handleFocusEnd()
	}
	return m, nil
}

func (m model) handleFocusEnd(extraCmds ...tea.Cmd) (tea.Model, tea.Cmd) {
	m.repetition++
	if m.repetition >= m.totalReps {
		m.phase = finishedPhase
		cmds := append(extraCmds, alertCmd("Finished!"), tea.Quit)
		return m, tea.Batch(cmds...)
	}
	m.phase = focusPhase
	m.timer = timer.New(m.focusDur)
	cmds := append(extraCmds, alertCmd("Start to focus!"), m.timer.Init())
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.quitting {
		return "See you later! 🍅\n"
	}

	s := titleStyle.Render("TOMATE") + "\n"

	var phaseName string
	var currentStyle lipgloss.Style
	var totalDur time.Duration

	switch m.phase {
	case focusPhase:
		phaseName = "FOCUS"
		currentStyle = focusStyle
		totalDur = m.focusDur
	case restPhase:
		phaseName = "REST"
		currentStyle = restStyle
		totalDur = m.restDur
	case finishedPhase:
		return s + "All done! 🍅\n"
	}

	s += fmt.Sprintf("Repetition: %d/%d\n", m.repetition+1, m.totalReps)
	s += currentStyle.Render(phaseName) + " " + m.timer.View() + "\n\n"

	remaining := m.timer.Timeout
	percent := 1.0 - float64(remaining)/float64(totalDur)
	s += m.progress.ViewAs(percent) + "\n"

	s += helpStyle.Render("space: pause/resume • r: reset • q: quit")

	return s
}

func alertCmd(msg string) tea.Cmd {
	return func() tea.Msg {
		_ = beeep.Alert("TOMATE", msg, "🍅")
		return nil
	}
}

func main() {
	beeep.AppName = "tomate"

	focusMin := flag.Int("fm", 0, "Focus minutes")
	focusSec := flag.Int("fs", 0, "Focus seconds")
	restMin := flag.Int("dm", 0, "Rest minutes")
	restSec := flag.Int("ds", 0, "Rest seconds")
	reps := flag.Int("r", 1, "Repetitions")
	flag.Parse()

	focusDur := time.Duration(*focusMin)*time.Minute + time.Duration(*focusSec)*time.Second
	restDur := time.Duration(*restMin)*time.Minute + time.Duration(*restSec)*time.Second

	if focusDur <= 0 {
		fmt.Println("Please provide a valid focus duration.")
		os.Exit(1)
	}

	m := model{
		phase:     focusPhase,
		totalReps: *reps,
		focusDur:  focusDur,
		restDur:   restDur,
		timer:     timer.New(focusDur),
		progress:  progress.New(progress.WithGradient("#FF9B9B", "#EF4444")),
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
