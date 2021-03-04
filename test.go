package main

import "fmt"

func main() {
	keyWidth := float64(640) / float64(52)

	switchAtFloat := 0.0
	switchAt := 0

	keyIndex := -1
	for i := 0; i < 640; i++ {
		if float64(i) >= switchAtFloat {
			keyIndex++
			switchAtFloat += keyWidth

			fmt.Printf("%d %d - %d\n", keyIndex, switchAt, int(switchAtFloat)+1)
			switchAt = int(switchAtFloat) + 1

		}
	}
}
