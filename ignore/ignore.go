package ignore

import (
	"k8s.io/helm/pkg/ignore"
	"os"
	"strings"
)

type Target interface {
	Path() string
	AbsolutePath() string
}

type Matcher struct {
	rules *ignore.Rules
}

func Parse(lines []string) (*Matcher, error) {
	rules, err := ignore.Parse(strings.NewReader(joinLines(lines)))
	if err != nil {
		return nil, err
	}
	return &Matcher{rules}, nil
}

func (r *Matcher) Match(t Target) bool {
	fi, err := os.Stat(t.AbsolutePath())
	if err != nil {
		return false
	}
	return r.rules.Ignore(t.Path(), fi)
}

func joinLines(lines []string) string {
	var b strings.Builder
	for idx, l := range lines {
		if strings.HasPrefix(l, "!") {
			continue // TODO: support accept pattern
		}
		if idx != 0 {
			b.WriteString("\n")
		}
		b.WriteString(l)
	}
	return b.String()
}
