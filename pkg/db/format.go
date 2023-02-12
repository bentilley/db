package db

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

// ColumnFormat right pads (with spaces) the value of each row so that the
// first letters of each column align.
func ColumnFormat(searchDetails [][]string) []string {
	maxWidths := lo.Reduce(
		searchDetails,
		func(agg []int, details []string, _ int) []int {
			for i, s := range details {
				if len(s) > agg[i] {
					agg[i] = len(s)
				}
			}
			return agg
		},
		make([]int, len(searchDetails[0])),
	)
	return lo.Map(
		searchDetails,
		func(details []string, _ int) string {
			N := len(details)
			formattedDetails := []string{}
			for i, s := range details {
				if i == N-1 { // last element doesn't need padding
					formattedDetails = append(formattedDetails, s)
					continue
				}
				formattedDetails = append(formattedDetails, pad(s, maxWidths[i]))
			}
			return strings.Join(formattedDetails, " ")
		},
	)
}

func pad(s string, width int) string {
	return fmt.Sprintf(fmt.Sprintf("%%-%ds", width), s)
}
