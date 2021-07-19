package services

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

var game *Game

func TestMain(m *testing.M) {
	// Discard logger output
	logrus.New().SetOutput(ioutil.Discard)

	// Create a new game
	game = &Game{
		id:   "test",
		Rows: [4][4]int{},
		Log:  logrus.NewEntry(logrus.New()),
	}

	os.Exit(m.Run())
}

func TestLeft(t *testing.T) {
	tests := []struct {
		Name     string
		Rows     [4][4]int
		Expected [4][4]int
	}{
		{
			Name: "simple move left test",
			Rows: [4][4]int{
				{2, 0, 0, 0},
				{0, 2, 0, 0},
				{0, 0, 2, 0},
				{0, 0, 0, 2},
			},
			Expected: [4][4]int{
				{2, 0, 0, 0},
				{2, 0, 0, 0},
				{2, 0, 0, 0},
				{2, 0, 0, 0},
			},
		},
		{
			Name: "complex move left test",
			Rows: [4][4]int{
				{2, 4, 0, 8},
				{2, 4, 0, 0},
				{0, 2, 0, 4},
				{8, 4, 2, 0},
			},
			Expected: [4][4]int{
				{2, 4, 8, 0},
				{2, 4, 0, 0},
				{2, 4, 0, 0},
				{8, 4, 2, 0},
			},
		},
		{
			Name: "simple add left test",
			Rows: [4][4]int{
				{2, 2, 0, 0},
				{2, 0, 2, 0},
				{2, 0, 0, 2},
				{0, 0, 2, 2},
			},
			Expected: [4][4]int{
				{4, 0, 0, 0},
				{4, 0, 0, 0},
				{4, 0, 0, 0},
				{4, 0, 0, 0},
			},
		},
		{
			Name: "complex add left test",
			Rows: [4][4]int{
				{2, 0, 2, 4},
				{4, 2, 2, 8},
				{2, 4, 0, 4},
				{4, 8, 2, 2},
			},
			Expected: [4][4]int{
				{4, 4, 0, 0},
				{4, 4, 8, 0},
				{2, 8, 0, 0},
				{4, 8, 4, 0},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			game.Rows = test.Rows
			game.MoveLeft()

			if ok := reflect.DeepEqual(game.Rows, test.Expected); !ok {
				out, err := game.Draw()
				if err != nil {
					t.Logf("drawing failed with err %v (%s)", err, test.Name)
					t.FailNow()
				}

				t.Logf("Test failed: %s\n%s", test.Name, out)
				t.FailNow()
			}
		})
	}
}
