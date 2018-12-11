package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hiroara/bunshin-no-jutsu/filesync"
)

func runSync(srcDir, destDir string) error {
	return run(srcDir, destDir, func(target filesync.Target) error {
		d, err := target.Copy(destDir)
		if err != nil {
			return err
		}
		if d != nil {
			fmt.Printf("%s => %s\n", target.AbsolutePath(), d.AbsolutePath())
		}
		return nil
	})
}

func run(srcDir, destDir string, f func(filesync.Target) error) error {
	return withCheck(srcDir, destDir, func() error {
		d := filesync.NewDirectory(srcDir, ".")
		ts, err := d.ListTargets()
		if err != nil {
			return err
		}
		for _, t := range ts {
			err := f(t)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func withCheck(srcDir, destDir string, f func() error) error {
	err := checkDestinationAvailability(destDir)
	if err != nil {
		return err
	}
	if !confirmSync(srcDir, destDir) {
		return nil
	}
	return f()
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
