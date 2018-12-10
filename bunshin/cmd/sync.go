package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func runSync() {
	if confirmSync() {
		fmt.Println("YEAH!")
	} else {
		fmt.Println("OH..")
	}
}

func confirmSync() bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Are you sure to sync files? (Y/n): ")
		text, _ := reader.ReadString('\n')
		switch strings.ToLower(strings.TrimSpace(text)) {
		case "", "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Println("Please input `y` or `n`.")
			continue
		}
	}
}
