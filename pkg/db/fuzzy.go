package db

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

// FuzzyFind returns the index of the selected item.
func FuzzyFind(input []string) (int, error) {
	maxLineNum := int(math.Floor(math.Log10(float64(len(input)))) + 1)
	indexedInput := lo.Map(
		input,
		func(line string, i int) string {
			return fmt.Sprintf("%s %s", pad(strconv.Itoa(i+1), maxLineNum), line)
		},
	)
	fzfCmd := exec.Command("fzf", "--height", "100%")
	fzfCmd.Stdin = strings.NewReader(strings.Join(indexedInput, "\n"))
	fzfCmd.Stderr = os.Stderr
	out, err := fzfCmd.Output()
	if err != nil {
		return 0, fmt.Errorf("fun fuzzy command: %w", err)
	}
	indexStr := strings.SplitN(string(out), " ", 2)[0]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return 0, fmt.Errorf("parse index: %w", err)
	}
	return index, nil
}
