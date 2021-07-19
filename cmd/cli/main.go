package main

import (
	"fmt"
	"game/config"
	"game/services"
	"os"

	"github.com/eiannone/keyboard"
)

func main() {
	cfg := config.New()
	game := services.NewGame("1", cfg.Log)

	out, err := game.Draw()
	if err != nil {
		game.Log.Error(err)
	}

	fmt.Println(out)

	for !game.Done {
		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			game.Log.Fatal(err)
		}

		fmt.Printf("input: %c\n", char)

		switch char {
		case 'a':
			game.MoveLeft()
			game.Move()

		case 'd':
			game.MoveLeft()
			game.Move()
		case 'q':
			os.Exit(0)
		}

		out, err := game.Draw()
		if err != nil {
			game.Log.Error(err)
		}

		fmt.Print(out)
	}

}
