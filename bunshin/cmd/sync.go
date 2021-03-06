package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/hiroara/bunshin-no-jutsu/filesync"
	"github.com/hiroara/bunshin-no-jutsu/ignore"
)

func runSync(srcDir, destDir string, dryrun bool, del bool, ignLines []string) (bool, error) {
	return withCheck(srcDir, destDir, dryrun, del, func() error {
		im, err := ignore.Parse(ignLines)
		if err != nil {
			return err
		}

		targets, idx, err := listFilesWithIndex(destDir, del, im)
		if err != nil {
			return err
		}

		if !dryrun {
			err = os.MkdirAll(destDir, 0777)
			if err != nil {
				return err
			}
		}
		err = run(srcDir, destDir, dryrun, del, im, func(target filesync.Target) error {
			newFile, d, err := filesync.Copy(target, destDir, dryrun)
			if err != nil {
				return err
			}
			delete(idx, target.Path())
			if d == nil {
				return nil
			}
			if newFile {
				fmt.Printf("%s => %s (new)\n", target.String(), d.String())
			} else {
				fmt.Printf("%s => %s\n", target.String(), d.String())
			}
			return nil
		})
		if err != nil {
			return err
		}
		return deleteFilesWithIndex(srcDir, targets, idx, dryrun)
	})
}

func deleteFilesWithIndex(srcPrefix string, targets []filesync.Target, index map[string]int, dryrun bool) error {
	indices := make([]int, 0, len(index))
	for _, i := range index {
		indices = append(indices, i)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(indices)))
	for _, i := range indices {
		t := targets[i]
		if filesync.ContainsSymlink(srcPrefix, t.Path()) {
			continue
		}
		err := t.Delete(dryrun)
		if err != nil {
			return err
		}
		if dryrun {
			fmt.Printf("%s will be deleted.\n", t.String())
		} else {
			fmt.Printf("%s is deleted.\n", t.String())
		}
	}
	return nil
}

func listFilesWithIndex(dir string, del bool, im *ignore.Matcher) ([]filesync.Target, map[string]int, error) {
	if del {
		d := filesync.NewDirectory(dir, ".")
		ts, err := d.ListTargets(im)
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

func run(srcDir, destDir string, dryrun, del bool, im *ignore.Matcher, f func(filesync.Target) error) error {
	d := filesync.NewDirectory(srcDir, ".")
	ts, err := d.ListTargets(im)
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

func withCheck(srcDir, destDir string, dryrun, del bool, f func() error) (bool, error) {
	err := checkDestinationAvailability(destDir)
	if err != nil {
		return false, err
	}
	if !confirmSync(srcDir, destDir, dryrun, del) {
		return false, nil
	}
	return true, f()
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
