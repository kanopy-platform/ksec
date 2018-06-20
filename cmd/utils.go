package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func outputTabular(lines []string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	w.Flush()
}
