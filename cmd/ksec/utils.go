package main

import (
	"bufio"
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

func askConfirmation(message string) bool {
	fmt.Printf("%s [y/N]: ", message)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response := scanner.Text()
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if response == "y" {
		return true
	}

	return false
}
