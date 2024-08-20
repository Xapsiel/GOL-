package app

import (
	"crypton/internal/game"
	"crypton/internal/point"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type model struct {
	stage *string

	menu       menu
	settings   settings
	playground playground
	cursor     point.Point
}
type menu struct {
	cursor   int
	choices  []string
	selected int
}

type settings struct {
	areaN     int
	arena     [][]point.Point
	cursor    int
	choices   []string
	selected  *int
	textInput textinput.Model
}
type playground struct {
	choices  [][]point.Point
	selected map[string]point.Point
	cursor   point.Point
	arena    *game.GIL
}

func InitialModel(n int) model {
	stage := "menu"
	playgroundstage := -1
	t := textinput.New()
	t.Cursor.Style = cursorStyle
	t.Placeholder = "Размер арены"
	t.Focus()
	t.PromptStyle = focusedStyle
	t.TextStyle = focusedStyle
	return model{
		stage: &stage,
		menu: menu{
			choices:  []string{"start", "end"},
			selected: -1,
		},
		settings: settings{
			areaN:     10,
			choices:   []string{"Размер арены: ", "Расставить клетки", "start", "Назад"},
			selected:  &playgroundstage,
			textInput: t,
		},
		playground: playground{
			choices:  make([][]point.Point, 0),
			selected: make(map[string]point.Point, 0),
			cursor:   point.Point{X: 1, Y: 1},
			arena:    game.NewGIL(10),
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch *(m.settings.selected) {
	case 0:

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
				m.settings.areaN, _ = strconv.Atoi(m.settings.textInput.Value())
				m.playground.arena = game.NewGIL(m.settings.areaN)

			}
			// We handle errors just like any other message
		}
		m.settings.textInput, _ = m.settings.textInput.Update(msg)
		break
	case 1:
		*m.stage = "playground"
	}
	switch *(m.stage) {
	case "menu":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "w":
				if m.menu.cursor > 0 {
					m.menu.cursor--
				}
			case "down", "j":
				if m.menu.cursor < len(m.menu.choices)-1 {
					m.menu.cursor++
				}
			case "enter", " ":
				m.menu.selected = m.menu.cursor
			}

		}
	case "settings":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "w":
				if m.settings.cursor > 0 {
					m.settings.cursor--
				}
			case "down", "j":
				if m.settings.cursor < len(m.settings.choices)-1 {
					m.settings.cursor++
				}
			case "enter", " ":
				*m.settings.selected = m.settings.cursor
			}
		}
	case "playground":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q", "esq":
				return m, tea.Quit
			case "left", "a":
				if m.playground.cursor.Y > 1 {
					m.playground.cursor.Y--
				}
			case "rigth", "d":
				if m.playground.cursor.Y < len(m.settings.arena)-1 {
					*m.stage = strconv.Itoa(m.settings.areaN)
					m.playground.cursor.Y++
				}
			case "up", "w":
				if m.playground.cursor.X > 1 {
					m.playground.cursor.X--
				}
			case "down", "s":
				if m.playground.cursor.X < m.settings.areaN-2 {
					m.playground.cursor.X++
				}
			case "enter", " ":
				if ok := m.playground.arena.Area[m.playground.cursor.X][m.playground.cursor.Y].Status; ok == "dead" {
					m.playground.arena.Area[m.playground.cursor.X][m.playground.cursor.Y].Status = "live"
					m.playground.selected[fmt.Sprintf("%s-%s", m.playground.cursor.X, m.playground.cursor.Y)] = m.playground.cursor
				} else {
					m.playground.arena.Area[m.playground.cursor.X][m.playground.cursor.Y].Status = "dead"
					delete(m.playground.selected, fmt.Sprintf("%s-%s", m.playground.cursor.X, m.playground.cursor.Y))
				}
			}

			return m, nil

		}
	}
	return m, nil
}

func (m model) View() string {

	s := fmt.Sprintf("%s\n\n", *m.stage)
	cursor := " "
	switch *(m.settings.selected) {
	case 0:
		s += m.settings.textInput.View() + "\n"
	}
	switch *(m.stage) {
	case "menu":
		for i := range m.menu.choices {
			if i == m.menu.cursor && i != m.menu.selected {
				cursor = "->"
			} else if i == m.menu.selected {
				*m.stage = "settings"
				cursor = "x"
			} else {
				cursor = " "
			}
			s += fmt.Sprintf("%s %s\n", cursor, m.menu.choices[i])

		}
		return s
	case "settings":
		for i := range m.settings.choices {
			if i == m.settings.cursor && i != *m.settings.selected {
				cursor = "->"
			} else if i == *m.settings.selected {
				cursor = "x"
				*m.settings.selected = i
			} else {
				cursor = " "
			}
			s += fmt.Sprintf("%s %s\n", cursor, m.settings.choices[i])

		}
		return s
	case "playground":
		for x := range m.playground.arena.Area {
			if x == 0 || x == m.settings.areaN-1 {
				continue
			}
			for y := range m.playground.arena.Area {
				if y == 0 || y == m.settings.areaN-1 {
					continue
				}
				m.playground.arena.Area[x][y].Status = "dead"
				for _, elem := range m.playground.selected {
					m.playground.arena.Area[elem.X][elem.Y].Status = "busy"

				}
				if x == m.playground.cursor.X && y == m.playground.cursor.Y {

					m.playground.arena.Area[x][y].Status = "active"
				} else {

				}

			}

		}
		return m.playground.arena.Print()
	}

	return s
}
