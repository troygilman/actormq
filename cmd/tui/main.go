package main

import "github.com/troygilman/actormq/tui"

func main() {
	if err := tui.Run(); err != nil {
		panic(err)
	}
}
