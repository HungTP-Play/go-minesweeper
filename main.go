package main

import (
	"fmt"
	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
)

type MineSweeperModel struct {
	Width         int
	Height        int
	NumberOfMines int
	Board         [][]int
	Flagged       [][]bool
	Opened        [][]bool
	IsGameOver    bool
	IsWin         bool
	CursorX       int
	CursorY       int
}

func InitModel(width, height, numberOfMines int) tea.Model {
	board := make([][]int, height)
	flagged := make([][]bool, height)
	opened := make([][]bool, height)
	for i := 0; i < height; i++ {
		board[i] = make([]int, width)
		flagged[i] = make([]bool, width)
		opened[i] = make([]bool, width)
	}
	m := MineSweeperModel{
		Width:         width,
		Height:        height,
		NumberOfMines: numberOfMines,
		Board:         board,
		Flagged:       flagged,
		IsGameOver:    false,
		CursorX:       0,
		CursorY:       0,
		Opened:        opened,
		IsWin:         false,
	}
	m.placeMines()
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if m.Board[i][j] == -1 {
				continue
			}
			m.Board[i][j] = m.countMines(j, i)
		}
	}
	return m
}

func (m MineSweeperModel) Init() tea.Cmd {
	// Set the initial window size to 100 x 100
	return nil
}

func (m *MineSweeperModel) countMines(x, y int) int {
	count := 0
	for i := y - 1; i <= y+1; i++ {
		if i < 0 || i >= m.Height {
			continue
		}
		for j := x - 1; j <= x+1; j++ {
			if j < 0 || j >= m.Width {
				continue
			}
			if m.Board[i][j] == -1 {
				count++
			}
		}
	}
	return count
}

func (m *MineSweeperModel) placeMines() {
	for i := 0; i < m.NumberOfMines; i++ {
		x := rand.Intn(m.Width)
		y := rand.Intn(m.Height)
		if m.Board[y][x] == -1 {
			i--
			continue
		}
		m.Board[y][x] = -1
	}
}

func (m *MineSweeperModel) Open(x, y int) {
	if m.Board[y][x] == -1 {
		m.IsGameOver = true
		return
	}
	m.open(x, y)
}

func (m *MineSweeperModel) open(x, y int) {
	if m.Board[y][x] != 0 {
		m.Opened[y][x] = true
		return
	}
	m.Opened[y][x] = true

	for i := y - 1; i <= y+1; i++ {
		if i < 0 || i >= m.Height {
			continue
		}
		for j := x - 1; j <= x+1; j++ {
			if j < 0 || j >= m.Width {
				continue
			}
			if m.Opened[i][j] {
				continue
			}
			m.open(j, i)
		}
	}
}

func (m *MineSweeperModel) ToggleFlag(x, y int) {
	if m.Opened[y][x] {
		return
	}

	m.Flagged[y][x] = !m.Flagged[y][x]

	// Check if win
	rightFlags := 0
	for i := 0; i < m.Height; i++ {
		for j := 0; j < m.Width; j++ {
			if m.Flagged[i][j] && m.Board[i][j] == -1 {
				rightFlags++
			}
		}
	}

	if rightFlags == m.NumberOfMines {
		m.IsWin = true
	}
}

// Update function
func (m MineSweeperModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.CursorY > 0 {
				m.CursorY--
			}
			return m, nil
		case "down", "j":
			if m.CursorY < m.Height-1 {
				m.CursorY++
			}
			return m, nil
		case "left", "h":
			if m.CursorX > 0 {
				m.CursorX--
			}
			return m, nil
		case "right", "l":
			if m.CursorX < m.Width-1 {
				m.CursorX++
			}
			return m, nil
		case tea.KeySpace.String():
			// Open
			m.Open(m.CursorX, m.CursorY)
			if m.IsGameOver {
				return m, tea.Quit
			}
			return m, nil
		case "f":
			// Toggle flag
			m.ToggleFlag(m.CursorX, m.CursorY)
			if m.IsWin {
				// Open all
				for i := 0; i < m.Height; i++ {
					for j := 0; j < m.Width; j++ {
						m.Opened[i][j] = true
					}
				}
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		return m, nil
	}
	return m, nil
}

func calculateRemainingFlags(m MineSweeperModel) int {
	remainingFlags := m.NumberOfMines
	for i := 0; i < m.Height; i++ {
		for j := 0; j < m.Width; j++ {
			if m.Flagged[i][j] {
				remainingFlags--
			}
		}
	}
	return remainingFlags
}

// View the board with the given model
// Show the board
// - 🔺 for flagged
// - number
// - blank for empty (opened)
// - ? for closed
func (m MineSweeperModel) View() string {
	uiString := ""
	uiString += fmt.Sprintf("Mines: %d\n", m.NumberOfMines)
	uiString += fmt.Sprintf("Remaining Flags: %d\n", calculateRemainingFlags(m))
	uiString += "\n"
	if m.IsGameOver {
		uiString += "Game Over 😭\n"
	} else if m.IsWin {
		uiString += "You Win 🎉🎉🎉\n"
	}

	for i := 0; i < m.Height; i++ {
		for j := 0; j < m.Width; j++ {
			// Cursor
			if i == m.CursorY && j == m.CursorX {
				uiString += "🤔"
				continue
			}

			if m.Flagged[i][j] {
				uiString += "🧡"
				continue
			}

			if m.Opened[i][j] {
				if m.Board[i][j] == -1 {
					uiString += "💣"
				} else if m.Board[i][j] == 0 {
					uiString += "  "
				} else {
					uiString += fmt.Sprintf(" %d", m.Board[i][j])
				}
			} else {
				uiString += "🟦"
			}
		}
		uiString += "\n"
	}
	return uiString
}

// Wrap the model in a program
func main() {
	p := tea.NewProgram(InitModel(30, 23, 99))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
