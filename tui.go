package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ethan-green3/pokedexcli/pokeapi"
)

var (
	appStyle = lipgloss.NewStyle().
			Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1)

	promptLineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))

	detailLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("69"))

	detailValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))
)

type viewMode int

const (
	modeShell viewMode = iota
	modePokedexList
	modePokemonDetail
)

type pokemonItem struct {
	name    string
	pokemon pokeapi.PokemonToCatch
}

func (p pokemonItem) Title() string       { return p.name }
func (p pokemonItem) Description() string { return "Caught Pokémon" }
func (p pokemonItem) FilterValue() string { return p.name }

type model struct {
	cfg config

	input    textinput.Model
	viewport viewport.Model

	history []string

	width  int
	height int
	ready  bool

	mode viewMode

	pokedexList list.Model
	selected    *pokeapi.PokemonToCatch
}

func newModel(cfg config) model {
	ti := textinput.New()
	ti.Prompt = "Pokedex > "
	ti.Placeholder = "..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 60

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false

	l := list.New([]list.Item{}, delegate, 20, 10)
	l.Title = "Your Pokedex"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(true)
	l.DisableQuitKeybindings()

	m := model{
		cfg:   cfg,
		input: ti,
		history: []string{
			titleStyle.Render("Welcome to Pokedex TUI"),
			subtitleStyle.Render("Type help to see available commands."),
			helpStyle.Render("Press Ctrl+C to quit."),
		},
		mode:        modeShell,
		pokedexList: l,
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case modePokedexList:
		return m.updatePokedexList(msg)
	case modePokemonDetail:
		return m.updatePokemonDetail(msg)
	default:
		return m.updateShell(msg)
	}
}

func (m model) updateShell(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		contentHeight := max(msg.Height-9, 8)

		contentWidth := max(msg.Width-6, 20)

		if !m.ready {
			m.viewport = viewport.New(contentWidth, contentHeight)
			m.ready = true
		} else {
			m.viewport.Width = contentWidth
			m.viewport.Height = contentHeight
		}

		m.input.Width = max(20, msg.Width-18)
		m.pokedexList.SetSize(contentWidth, contentHeight)
		m.syncViewport()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			raw := strings.TrimSpace(m.input.Value())
			if raw == "" {
				return m, nil
			}

			m.history = append(m.history, promptLineStyle.Render("Pokedex > ")+raw)

			args := cleanInput(raw)
			if len(args) > 0 {
				switch args[0] {
				case "exit":
					m.history = append(m.history, "Closing the Pokedex...")
					m.syncViewport()
					return m, tea.Quit

				case "pokedex":
					m.refreshPokedexList()
					if len(m.pokedexList.Items()) == 0 {
						m.history = append(m.history, "Your pokedex is empty. Catch a Pokémon first.")
						m.input.SetValue("")
						m.syncViewport()
						return m, nil
					}
					m.mode = modePokedexList
					m.input.SetValue("")
					return m, nil

				case "inspect":
					if len(args) < 2 {
						m.history = append(m.history, errorStyle.Render("Error: usage: inspect <pokemon>"))
						m.input.SetValue("")
						m.syncViewport()
						return m, nil
					}

					pokemon, ok := m.cfg.Pokedex[args[1]]
					if !ok {
						m.history = append(m.history, errorStyle.Render("Error: that pokemon is not in your Pokedex, you need to catch it first!"))
						m.input.SetValue("")
						m.syncViewport()
						return m, nil
					}

					m.selected = &pokemon
					m.mode = modePokemonDetail
					m.input.SetValue("")
					return m, nil
				}
			}

			output, err := runCommandCapture(&m.cfg, raw)
			if output != "" {
				m.history = append(m.history, output)
			}
			if err != nil {
				m.history = append(m.history, errorStyle.Render("Error: "+err.Error()))
			}

			m.input.SetValue("")
			m.syncViewport()
			return m, nil
		}
	}

	var cmds []tea.Cmd

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	if m.ready {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) updatePokedexList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		contentHeight := max(msg.Height-8, 8)

		contentWidth := max(msg.Width-6, 20)

		m.pokedexList.SetSize(contentWidth, contentHeight)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			m.mode = modeShell
			return m, nil
		case "enter":
			selectedItem, ok := m.pokedexList.SelectedItem().(pokemonItem)
			if !ok {
				return m, nil
			}
			p := selectedItem.pokemon
			m.selected = &p
			m.mode = modePokemonDetail
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.pokedexList, cmd = m.pokedexList.Update(msg)
	return m, cmd
}

func (m model) updatePokemonDetail(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			m.mode = modePokedexList
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "\n  Loading..."
	}

	switch m.mode {
	case modePokedexList:
		return appStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				titleStyle.Render("Pokedex CLI"),
				subtitleStyle.Render("Select a Pokémon and press Enter to inspect"),
				"",
				panelStyle.Width(m.pokedexList.Width()+2).Render(m.pokedexList.View()),
				"",
				helpStyle.Render("↑/↓ move • Enter inspect • Esc back • q quit"),
			),
		)

	case modePokemonDetail:
		return appStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				titleStyle.Render("Pokedex CLI"),
				subtitleStyle.Render("Pokémon Detail"),
				"",
				panelStyle.Width(max(50, m.width-10)).Render(m.renderPokemonDetail()),
				"",
				helpStyle.Render("Esc back • q quit"),
			),
		)

	default:
		header := lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("Pokedex CLI"),
			subtitleStyle.Render("Bubble Tea shell with interactive Pokedex"),
		)

		outputPanel := panelStyle.
			Width(m.viewport.Width + 2).
			Height(m.viewport.Height + 2).
			Render(m.viewport.View())

		footer := lipgloss.JoinVertical(
			lipgloss.Left,
			panelStyle.Width(m.viewport.Width+2).Render(m.input.View()),
			helpStyle.Render("Commands: help, map, mapb, explore <area>, catch <pokemon>, inspect <pokemon>, pokedex, exit"),
		)

		ui := lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			"",
			outputPanel,
			"",
			footer,
		)

		return appStyle.Render(ui)
	}
}

func (m *model) renderPokemonDetail() string {
	if m.selected == nil {
		return "No Pokémon selected."
	}

	p := *m.selected
	var b strings.Builder

	b.WriteString(detailLabelStyle.Render("Name: ") + detailValueStyle.Render(p.Name) + "\n")
	b.WriteString(detailLabelStyle.Render("Height: ") + detailValueStyle.Render(fmt.Sprintf("%d", p.Height)) + "\n")
	b.WriteString(detailLabelStyle.Render("Weight: ") + detailValueStyle.Render(fmt.Sprintf("%d", p.Weight)) + "\n")
	b.WriteString("\n")
	b.WriteString(detailLabelStyle.Render("Stats:") + "\n")
	for _, stat := range p.Stats {
		fmt.Fprintf(&b, "  - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	b.WriteString("\n")
	b.WriteString(detailLabelStyle.Render("Types:") + "\n")
	for _, t := range p.Types {
		fmt.Fprintf(&b, "  - %s\n", t.Type.Name)
	}

	return strings.TrimRight(b.String(), "\n")
}

func (m *model) syncViewport() {
	if !m.ready {
		return
	}
	m.viewport.SetContent(strings.Join(m.history, "\n\n"))
	m.viewport.GotoBottom()
}

func (m *model) refreshPokedexList() {
	names := make([]string, 0, len(m.cfg.Pokedex))
	for name := range m.cfg.Pokedex {
		names = append(names, name)
	}
	sort.Strings(names)

	items := make([]list.Item, 0, len(names))
	for _, name := range names {
		p := m.cfg.Pokedex[name]
		items = append(items, pokemonItem{
			name:    p.Name,
			pokemon: p,
		})
	}

	m.pokedexList.SetItems(items)
	if len(items) > 0 {
		m.pokedexList.Select(0)
	}
}

func runCommandCapture(cfg *config, raw string) (string, error) {
	args := cleanInput(raw)
	if len(args) == 0 {
		return "", nil
	}

	cmd, ok := commands[args[0]]
	if !ok {
		return "", fmt.Errorf("unknown command: %s", args[0])
	}

	return captureStdout(func() error {
		return cmd.callback(cfg, args...)
	})
}

func captureStdout(fn func() error) (string, error) {
	oldStdout := os.Stdout

	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}

	os.Stdout = w

	outCh := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outCh <- buf.String()
	}()

	runErr := fn()

	_ = w.Close()
	os.Stdout = oldStdout

	output := <-outCh
	_ = r.Close()

	return strings.TrimRight(output, "\n"), runErr
}
