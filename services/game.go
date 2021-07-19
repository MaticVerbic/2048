package services

import (
	"bytes"
	"math/rand"
	"text/template"
	"time"

	"github.com/go-errors/errors"
	"github.com/sirupsen/logrus"
)

type Left int

const (
	moveCheckL Left = iota
	moveCheckR
	moveLeft
	addLeft
	invalidLeft
)

type Game struct {
	id   string
	Rows [4][4]int
	Log  *logrus.Entry
	Done bool
}

func NewGame(id string, log *logrus.Entry) *Game {
	// seed each game separately
	rand.Seed(time.Now().UTC().UnixNano())

	// generate initial 2 numbers, they can either be 2 or 4.
	firstY := rand.Intn(4)
	secondY := rand.Intn(4)

	firstX := rand.Intn(4)
	secondX := rand.Intn(4)
	for firstX != secondX && firstY != secondY {
		secondX = rand.Intn(4)
	}

	digits := []int{2, 4}
	firstDigit := digits[rand.Intn(len(digits))]
	secondDigit := digits[rand.Intn(len(digits))]

	rows := [4][4]int{}
	rows[firstY][firstX] = firstDigit
	rows[secondY][secondX] = secondDigit

	// return game
	return &Game{
		id:   id,
		Rows: rows,
		Log:  log,
	}
}

func (g *Game) Draw() (string, error) {
	tmp := `{{ range .Rows }}
| {{ range . }}{{ . }} | {{ end }}
{{ end }}`
	tmpl, err := template.New("test").Parse(tmp)
	if err != nil {
		return "", errors.WrapPrefix(err, "failed to parse template", 0)
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, g); err != nil {
		return "", errors.WrapPrefix(err, "failed to execute template", 0)
	}

	return buf.String(), nil
}

func (g *Game) Move() (bool, error) {
	g.addOne()

	if g.Done {
		return true, nil
	}

	return false, nil
}

func (g *Game) addOne() {
	empty := [][]int{}

	for k, row := range g.Rows {
		for j, val := range row {
			if val == 0 {
				empty = append(empty, []int{k, j})
			}
		}
	}

	if len(empty) == 0 {
		g.Done = true
		return
	}

	choice := empty[rand.Intn(len(empty))]
	four := rand.Intn(1)

	if four == 1 {
		g.Rows[choice[0]][choice[1]] = 4
	}

	g.Rows[choice[0]][choice[1]] = 2
}

func (g *Game) MoveLeft() {
	for row := 0; row < len(g.Rows); row++ {
		_, checkL, checkR := 0, 0, 1

		for checkL < len(g.Rows[row])-1 {
			switch evalLeft(g.Rows[row], checkL, checkR) {
			case moveCheckL:
				checkL++
				checkR++
			case moveCheckR:
				checkR++
			case moveLeft:
				g.Rows[row][checkL] = g.Rows[row][checkR]
				g.Rows[row][checkR] = 0
				checkR = checkL + 1
			case addLeft:
				g.Rows[row][checkL] += g.Rows[row][checkR]
				g.Rows[row][checkR] = 0
				checkR++
				checkL++

			}

			if checkR > len(g.Rows[row])-1 {
				checkR = len(g.Rows[row]) - 1
			}
		}
	}
}

func evalLeft(rows [4]int, left, right int) Left {
	if rows[left] != 0 && rows[right] != 0 && rows[left] != rows[right] {
		return moveCheckL
	}

	if rows[left] == 0 && rows[right] != 0 {
		return moveLeft
	}

	if rows[left] != 0 && rows[right] == rows[left] {
		return addLeft
	}

	if rows[right] == 0 && right == len(rows)-1 {
		return moveCheckL
	}

	if rows[right] == 0 {
		return moveCheckR
	}

	return invalidLeft
}

func (g *Game) MoveRight() {

}
