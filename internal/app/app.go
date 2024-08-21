package app

import (
	"crypton/internal/game"
	"crypton/internal/point"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	timeout      = time.Second * 1
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

type model struct {
	stage      *string
	idToS      map[int]string
	Stoid      map[string]int
	menu       menu
	settings   settings
	playground playground
	play       play
}
type menu struct {
	cursor   int
	choices  []string
	selected int
}

type settings struct {
	areaN     int
	cursor    int
	choices   []string
	selected  *int
	textInput textinput.Model
}
type playground struct {
	choices  [][]point.Point
	selected map[string]point.Point
	cursor   point.Point
	arena    *game.GOL
	back     []string
}
type play struct {
	ticker *time.Ticker
}

func InitialModel(n int) model {
	stage := "menu"
	playgroundstage := -1
	t := textinput.New()
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
			areaN:     n,
			choices:   []string{"size", "fill in", "play"},
			selected:  &playgroundstage,
			textInput: t,
		},
		playground: playground{
			choices:  make([][]point.Point, 0),
			selected: make(map[string]point.Point, 0),
			cursor:   point.Point{X: 1, Y: 1},
			arena:    game.NewGIL(10),
			back:     []string{"back"},
		},
		play: play{
			ticker: time.NewTicker(timeout),
		},
	}
}

func (m model) Init() tea.Cmd {
	m.play.ticker = time.NewTicker(timeout) // Создание таймера на 100 мс
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch *(m.settings.selected) {
	case 0:

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter, tea.KeyCtrlC:
				m.settings.areaN, _ = strconv.Atoi(m.settings.textInput.Value())
				m.settings.areaN += 2
				m.playground.arena = game.NewGIL(m.settings.areaN)
			case tea.KeyBackspace, tea.KeyEsc:
				*m.stage = "settings"
				*m.settings.selected = -1
			}

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
				if m.settings.choices[m.settings.cursor] == "play" {
					*m.stage = "play"
				}

			}
		}

	case "playground":

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case tea.KeyBackspace.String(), "esq":
				*m.stage = "settings"
				*m.settings.selected = -1
			case "ctrl+c", "q":
				return m, tea.Quit
			case "left", "a":
				if m.playground.cursor.Y > 1 {
					m.playground.cursor.Y--
				}
			case "right", "d":
				if m.playground.cursor.Y < m.settings.areaN-2 {
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
				if status := m.playground.arena.Area[m.playground.cursor.X][m.playground.cursor.Y].Status; status == "dead" || status == "active" {
					m.playground.arena.Area[m.playground.cursor.X][m.playground.cursor.Y].Status = "live"
					m.playground.selected[fmt.Sprintf("%s-%s", m.playground.cursor.X, m.playground.cursor.Y)] = m.playground.cursor

				} else {

					m.playground.arena.Area[m.playground.cursor.X][m.playground.cursor.Y].Status = "dead"
					delete(m.playground.selected, fmt.Sprintf("%s-%s", m.playground.cursor.X, m.playground.cursor.Y))

				}
			}

			return m, nil

		}

	case "play":

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit

			}
		case timer.TickMsg:
			m.playground.arena.Step()
			m.playground.arena.Print()

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
				*m.stage = m.settings.choices[i]
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

				if x == m.playground.cursor.X && y == m.playground.cursor.Y && m.playground.arena.Area[x][y].Status != "busy" {

					m.playground.arena.Area[x][y].Status = "active"
				} else {
					m.playground.arena.Area[x][y].Status = "dead"
				}
				for _, elem := range m.playground.selected {
					m.playground.arena.Area[elem.X][elem.Y].Status = "busy"

				}

			}

		}
		busyP := "\tPoint\n"
		for i := range m.playground.selected {
			busyP += fmt.Sprintf("x=%v;y=%v\n", m.playground.selected[i].X, m.playground.selected[i].Y)
		}
		m.playground.arena.Print()
		return *m.playground.arena.GameMap + "\n" + busyP
	case "play":

		return s + *m.playground.arena.GameMap
	}

	return s
}

func Run(n int) {
	m := InitialModel(n)
	p := tea.NewProgram(m)

	go func() {
		for range m.play.ticker.C { // Отправка сообщений каждые 100 мс
			p.Send(timer.TickMsg{}) // Отправка сообщения в основной цикл
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("Ошибка:", err)
		os.Exit(1)
	}
}
