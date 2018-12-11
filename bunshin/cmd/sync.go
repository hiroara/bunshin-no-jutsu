package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hiroara/bunshin-no-jutsu/filesync"
)

func runSync(srcDir, destDir string) error {
	err := checkDestinationAvailability(destDir)
	if err != nil {
		return err
	}
	if confirmSync(srcDir, destDir) {
		filesync.NewDirectory(srcDir).Sync(filesync.NewDirectory(destDir))
	}
	return nil
}

func confirmSync(srcDir, destDir string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Sync: %s => %s\n", srcDir, destDir)
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

func checkDestinationAvailability(path string) error {
	if path == "" {
		return fmt.Errorf("Any destination directory is not configured. Please check your configuration file.")
	}
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if !stat.Mode().IsDir() {
		return fmt.Errorf("`%s` is not a directory. Please check it.", path)
	}
	return nil
}
