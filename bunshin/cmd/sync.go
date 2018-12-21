package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/hiroara/bunshin-no-jutsu/filesync"
)

func runSync(srcDir, destDir string, dryrun bool, del bool) error {
	targets, idx, err := listFilesWithIndex(destDir, del)
	if err != nil {
		return err
	}

	if !dryrun {
		err = os.MkdirAll(destDir, 0777)
		if err != nil {
			return err
		}
	}
	return withCheck(srcDir, destDir, dryrun, del, func() error {
		err = run(srcDir, destDir, dryrun, del, func(target filesync.Target) error {
			d, newFile, err := filesync.Copy(target.Prefix(), destDir, target.Path(), dryrun)
			if err != nil {
				return err
			}
			delete(idx, target.Path())
			if d == nil {
				return nil
			}
			if newFile {
				fmt.Printf("%s => %s (new)\n", target.AbsolutePath(), d.AbsolutePath())
			} else {
				fmt.Printf("%s => %s\n", target.AbsolutePath(), d.AbsolutePath())
			}
			return nil
		})
		if err != nil {
			return err
		}
		return deleteFilesWithIndex(targets, idx, dryrun)
	})
}

func deleteFilesWithIndex(targets []filesync.Target, index map[string]int, dryrun bool) error {
	indices := make([]int, 0, len(index))
	for _, i := range index {
		indices = append(indices, i)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(indices)))
	for _, i := range indices {
		t := targets[i]
		if dryrun {
			fmt.Printf("%s will be deleted.\n", t.AbsolutePath())
		} else {
			err := t.Delete()
			if err != nil {
				return err
			}
			fmt.Printf("%s is deleted.\n", t.AbsolutePath())
		}
	}
	return nil
}

func listFilesWithIndex(dir string, del bool) ([]filesync.Target, map[string]int, error) {
	if del {
		d := filesync.NewDirectory(dir, ".")
		ts, err := d.ListTargets()
		if err != nil {
			return nil, nil, err
		}
		indices := make(map[string]int, len(ts))
		for i, t := range ts {
			indices[t.Path()] = i
		}
		return ts, indices, nil
	} else {
		return []filesync.Target{}, make(map[string]int, 0), nil
	}
}

func run(srcDir, destDir string, dryrun, del bool, f func(filesync.Target) error) error {
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
}

func withCheck(srcDir, destDir string, dryrun, del bool, f func() error) error {
	err := checkDestinationAvailability(destDir)
	if err != nil {
		return err
	}
	if !confirmSync(srcDir, destDir, dryrun, del) {
		return nil
	}
	return f()
}

func confirmSync(srcDir, destDir string, dryrun, del bool) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s: %s => %s\n", operationName(dryrun, del), srcDir, destDir)
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

func operationName(dryrun, del bool) string {
	name := "Sync"
	if del {
		name += " with delete"
	}
	if dryrun {
		name += " (dry-run)"
	}
	return name
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
