package services

import (
	"bytes"
	"math/rand"
	"text/template"
	"time"

	"github.com/go-errors/errors"
	"github.com/sirupsen/logrus"
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

func (g *Game) reduceRow(row [4]int) [4]int {
	out := [4]int{}

	i := 0
	for _, elem := range row {
		if elem != 0 {
			out[i] = elem
			i++
		}
	}

	return out
}

func (g *Game) sumRow(row [4]int) [4]int {
	for i := len(row) - 1; i >= 1; i-- {
		left, right := row[i], row[i-1]
		if left == right && left != 0 {
			row[i] = left + right
			row[i-1] = 0
		}
	}
	return g.reduceRow(row)
}

func (g *Game) reverseArr(rows [4]int) [4]int {
	out := [4]int{}

	for i, elem := range rows {
		out[len(rows)-1-i] = elem
	}

	return out
}

func (g *Game) getCol(col int) ([4]int, error) {
	out := [4]int{}
	if col > 4 {
		return out, errors.New("index out of range")
	}

	for i, row := range g.Rows {
		out[i] = row[col]
	}

	return out, nil
}

func (g *Game) MoveLeft() {
	for r, row := range g.Rows {
		reduced := g.reduceRow(row)
		g.Rows[r] = g.sumRow(reduced)
	}
}

func (g *Game) MoveRight() {
	for r, row := range g.Rows {
		reduced := g.reduceRow(g.reverseArr(row))
		g.Rows[r] = g.reverseArr(g.sumRow(reduced))
	}
}

func (g *Game) MoveUp() error {
	for j := 0; j < len(g.Rows); j++ {
		col, err := g.getCol(j)
		if err != nil {
			return errors.WrapPrefix(err, "failed to get col", 0)
		}
		reduced := g.reduceRow(col)
		row := g.sumRow(reduced)

		for i := 0; i < len(row); i++ {
			g.Rows[i][j] = row[i]
		}
	}

	return nil
}

func (g *Game) MoveDown() error {
	for j := 0; j < len(g.Rows); j++ {
		col, err := g.getCol(j)
		if err != nil {
			return errors.WrapPrefix(err, "failed to get col", 0)
		}
		reduced := g.reduceRow(g.reverseArr(col))
		row := g.reverseArr(g.sumRow(reduced))

		for i := 0; i < len(row); i++ {
			g.Rows[i][j] = row[i]
		}
	}

	return nil
}
