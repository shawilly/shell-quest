package tui

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shanewilliams/shell-quest/internal/db"
	"github.com/shanewilliams/shell-quest/internal/shell"
	"github.com/shanewilliams/shell-quest/internal/shell/commands"
	"github.com/shanewilliams/shell-quest/internal/world"
)

type AppState int

const (
	StateGame AppState = iota
	StateProfileSelect
	StateTierSelect
	StateNameInput
	StateWelcome
	StateAdventureLog
	StateParentMode
)

// Model is the root Bubble Tea model.
type Model struct {
	state      AppState
	db         *db.DB
	player     *db.Player
	fs         *shell.FS
	executor   *shell.Executor
	runner     *world.MissionRunner
	gameWorld  *world.World
	missionIdx int
	cwd        string
	width      int
	height     int

	// shell pane
	shellInput textinput.Model
	shellVP    viewport.Model
	shellLines []string // raw lines for rebuilding viewport on resize

	// story pane
	clueText  string
	storyText string
	failCount int

	// profile selection
	profiles    []*db.Player
	selectedIdx int
	nameInput   textinput.Model

	// parent mode
	mathInput      textinput.Model
	mathA, mathB   int
	parentUnlocked bool
}

// NewGameModel creates a model ready to play a mission.
func NewGameModel(d *db.DB, player *db.Player, fs *shell.FS, ex *shell.Executor, runner *world.MissionRunner, startCWD, clue string) Model {
	m := Model{
		state:     StateGame,
		db:        d,
		player:    player,
		fs:        fs,
		executor:  ex,
		runner:    runner,
		cwd:       startCWD,
		clueText:  clue,
		storyText: "Welcome to Skull Island, young pirate! Arr!",
	}
	m.shellInput = newShellInput()
	m.shellInput.Focus()
	m.shellVP = newShellViewport(80, 20)
	m.mathInput = newMathInput()
	m.nameInput = newNameInput()
	return m
}

// NewStartupModel creates a model starting at the welcome screen.
func NewStartupModel(d *db.DB) Model {
	players, _ := d.ListPlayers()
	m := Model{
		state:    StateWelcome,
		db:       d,
		profiles: players,
	}
	m.nameInput = newNameInput()
	m.mathInput = newMathInput()
	m.shellInput = newShellInput()
	m.shellInput.Focus()
	m.shellVP = newShellViewport(80, 20)
	return m
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		rightWidth := msg.Width - msg.Width/2
		shellInner := rightWidth - 8
		if shellInner < 10 {
			shellInner = 10
		}
		vpHeight := msg.Height - 8
		if vpHeight < 4 {
			vpHeight = 4
		}
		m.shellVP = viewport.New(shellInner, vpHeight)
		m.shellVP.SetContent(strings.Join(m.shellLines, "\n"))
	case tea.KeyMsg:
		switch m.state {
		case StateWelcome:
			return m.handleWelcomeKey(msg)
		case StateGame:
			return m.handleKey(msg)
		case StateProfileSelect:
			return m.handleProfileKey(msg)
		case StateTierSelect:
			return m.handleTierKey(msg)
		case StateNameInput:
			return m.handleNameKey(msg)
		case StateAdventureLog:
			return m.handleAdventureLogKey(msg)
		case StateParentMode:
			return m.handleParentModeKey(msg)
		}
	}
	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case StateWelcome:
		return m.welcomeView()
	case StateGame:
		return m.gameView()
	case StateProfileSelect:
		return m.profileSelectView()
	case StateTierSelect:
		return m.tierSelectView()
	case StateNameInput:
		return m.nameInputView()
	case StateAdventureLog:
		return m.adventureLogView()
	case StateParentMode:
		return m.parentModeView()
	default:
		return "Loading..."
	}
}

func (m Model) handleWelcomeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEnter:
		m.state = StateProfileSelect
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyCtrlL:
		m.state = StateAdventureLog
	case tea.KeyCtrlP:
		m.state = StateParentMode
		m.mathA = rand.Intn(10) + 1
		m.mathB = rand.Intn(10) + 1
		m.mathInput.SetValue("")
		m.mathInput.Focus()
		m.parentUnlocked = false
	case tea.KeyEnter:
		return m.submitCommand()
	default:
		var cmd tea.Cmd
		m.shellInput, cmd = m.shellInput.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) handleAdventureLogKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc, tea.KeyEnter, tea.KeyCtrlL:
		m.state = StateGame
	}
	return m, nil
}

func (m Model) handleParentModeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.parentUnlocked {
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			m.state = StateGame
			m.parentUnlocked = false
		case tea.KeyRunes:
			if msg.String() == "q" {
				return m, tea.Quit
			}
		}
		return m, nil
	}

	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		m.state = StateGame
	case tea.KeyEnter:
		expected := fmt.Sprintf("%d", m.mathA+m.mathB)
		if m.mathInput.Value() == expected {
			m.parentUnlocked = true
			m.mathInput.Blur()
		} else {
			m.mathInput.SetValue("")
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.mathInput, cmd = m.mathInput.Update(msg)
	return m, cmd
}

func (m Model) handleProfileKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	count := len(m.profiles) + 1 // +1 for "New Profile"
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyUp:
		if m.selectedIdx > 0 {
			m.selectedIdx--
		}
	case tea.KeyDown:
		if m.selectedIdx < count-1 {
			m.selectedIdx++
		}
	case tea.KeyEnter:
		if m.selectedIdx == len(m.profiles) {
			// "New Profile" selected
			m.state = StateNameInput
			m.nameInput.SetValue("")
			m.nameInput.Focus()
		} else {
			// Existing profile selected
			m.player = m.profiles[m.selectedIdx]
			return m.startGame()
		}
	}
	return m, nil
}

func (m Model) handleTierKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	tiers := []string{"beginner", "explorer", "master"}
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyUp:
		if m.selectedIdx > 0 {
			m.selectedIdx--
		}
	case tea.KeyDown:
		if m.selectedIdx < len(tiers)-1 {
			m.selectedIdx++
		}
	case tea.KeyEnter:
		if m.selectedIdx != 0 {
			// Only Beginner is available; other tiers are coming soon
			return m, nil
		}
		tier := tiers[m.selectedIdx]
		player, err := m.db.CreatePlayer(strings.TrimSpace(m.nameInput.Value()), tier)
		if err != nil {
			return m, nil
		}
		m.player = player
		return m.startGame()
	}
	return m, nil
}

func (m Model) handleNameKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEnter:
		if len(strings.TrimSpace(m.nameInput.Value())) > 0 {
			m.state = StateTierSelect
			m.selectedIdx = 0
			m.nameInput.Blur()
		}
		return m, nil
	case tea.KeyEsc:
		return m, nil
	}
	var cmd tea.Cmd
	m.nameInput, cmd = m.nameInput.Update(msg)
	return m, cmd
}

// startGame initializes the game with the selected player.
func (m Model) startGame() (Model, tea.Cmd) {
	w, err := world.LoadWorld("skull_island")
	if err != nil {
		return m, tea.Quit
	}
	mission := w.Missions[0]

	fs := shell.NewFS()
	paths := make([]string, 0, len(mission.Filesystem))
	for p := range mission.Filesystem {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	for _, p := range paths {
		entry := mission.Filesystem[p]
		if entry.Type == "dir" {
			if err := fs.Mkdir(p, entry.Hidden); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := fs.WriteFile(p, entry.Content, entry.Hidden); err != nil {
				log.Fatal(err)
			}
		}
	}

	ex := shell.NewExecutor(fs)
	RegisterCommands(ex, m.player.Tier)
	runner := world.NewMissionRunner(mission)

	cwd := mission.StartingCWD
	if cwd == "" {
		cwd = "/"
	}

	m.fs = fs
	m.executor = ex
	m.runner = runner
	m.gameWorld = w
	m.missionIdx = 0
	m.cwd = cwd
	m.clueText = mission.StartingClue
	m.storyText = missionIntro(mission)
	m.state = StateGame
	m.shellLines = nil
	m.shellVP.SetContent("")
	m.failCount = 0
	return m, nil
}

// missionIntro returns the story pane text for the start of a mission.
func missionIntro(mission world.Mission) string {
	if mission.IntroDialogue != nil && mission.IntroDialogue.Text != "" {
		return mission.IntroDialogue.NPC + ": " + mission.IntroDialogue.Text
	}
	return "Welcome to Skull Island, young pirate! Arr!"
}

func (m Model) submitCommand() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(m.shellInput.Value())
	m.shellInput.SetValue("")
	if input == "" {
		return m, nil
	}

	m.appendOutput(PromptStyle.Render("$ ") + input)
	result := m.executor.Execute(input, m.cwd)

	switch {
	case result.ExitLevel:
		m.profiles, _ = m.db.ListPlayers()
		m.state = StateWelcome
		return m, nil
	case result.Clear:
		m.shellLines = nil
		m.shellVP.SetContent("")
		return m, nil
	case result.Error != "":
		m.appendOutput(ErrorStyle.Render(result.Error))
		m.failCount++
	default:
		for _, line := range strings.Split(result.Output, "\n") {
			m.appendOutput(OutputStyle.Render(line))
		}
		if result.NewCWD != "" {
			m.cwd = result.NewCWD
		}
		m.failCount = 0
	}

	m = m.applyMissionEvent(result)

	if m.failCount >= 3 && m.runner != nil {
		if obj := m.runner.CurrentObjective(); obj != nil {
			m.storyText = "Psst! Try using '" + obj.Command + "' to progress..."
			m.failCount = 0
		}
	}

	return m, nil
}

// applyMissionEvent advances mission state when a command triggers an objective.
func (m Model) applyMissionEvent(result shell.Result) Model {
	if m.runner == nil || result.Event == nil {
		return m
	}
	if !m.runner.HandleEvent(result.Event) {
		return m
	}

	if !m.runner.IsComplete() {
		if hint := m.runner.CurrentHint(); hint != "" {
			m.storyText = hint
		}
		if obj := m.runner.CurrentObjective(); obj != nil {
			m.clueText = "Objective " + string(rune('0'+m.runner.CurrentObjectiveIndex())) + " complete! Keep going..."
		}
		return m
	}

	completed := m.runner.Mission()
	m.appendOutput(SuccessStyle.Render(completed.SuccessMessage))
	if completed.BugTaunt != "" {
		m.appendOutput(SuccessStyle.Render(completed.BugTaunt))
	}

	nextIdx := m.missionIdx + 1
	if m.gameWorld == nil || nextIdx >= len(m.gameWorld.Missions) {
		m.storyText = completed.Treasure
		m.clueText = "You've conquered Skull Island! You are a true Pirate Master!"
		return m
	}

	return m.loadMission(nextIdx)
}

// loadMission sets up the filesystem and runner for the mission at idx.
func (m Model) loadMission(idx int) Model {
	mission := m.gameWorld.Missions[idx]

	paths := make([]string, 0, len(mission.Filesystem))
	for p := range mission.Filesystem {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	for _, p := range paths {
		entry := mission.Filesystem[p]
		if entry.Type == "dir" {
			if err := m.fs.Mkdir(p, entry.Hidden); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := m.fs.WriteFile(p, entry.Content, entry.Hidden); err != nil {
				log.Fatal(err)
			}
		}
	}

	m.runner = world.NewMissionRunner(mission)
	m.missionIdx = idx
	m.clueText = mission.StartingClue
	m.storyText = missionIntro(mission)
	if mission.StartingCWD != "" {
		m.cwd = mission.StartingCWD
	}
	m.appendOutput(SuccessStyle.Render("=== New Mission: " + mission.Title + " ==="))
	return m
}

func (m *Model) appendOutput(line string) {
	m.shellLines = append(m.shellLines, line)
	m.shellVP.SetContent(strings.Join(m.shellLines, "\n"))
	m.shellVP.GotoBottom()
}

// RegisterCommands registers all commands with the executor based on player tier.
func RegisterCommands(ex *shell.Executor, tier string) {
	// Beginner always
	ex.Register(commands.NewPwd())
	ex.Register(commands.NewLs())
	ex.Register(commands.NewCd())
	ex.Register(commands.NewCat())
	ex.Register(commands.NewEcho())
	ex.Register(commands.NewClear())
	ex.Register(commands.NewExit())
	ex.Register(commands.NewHelp(tier))

	// Explorer commands always registered (all 10 missions need them; only beginner tier is selectable)
	ex.Register(commands.NewMkdir())
	ex.Register(commands.NewTouch())
	ex.Register(commands.NewCp())
	ex.Register(commands.NewMv())
	ex.Register(commands.NewRm())
	ex.Register(commands.NewFind())

	if tier == "master" {
		ex.Register(commands.NewGrep())
		ex.Register(commands.NewChmod())
		ex.Register(commands.NewMan())
		ex.Register(commands.NewHistory(func() []string { return ex.History }))
	}
}
