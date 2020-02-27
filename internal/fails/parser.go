package fails

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	indent        = "    "
	failIndicator = "--- FAIL:"
)

// Print focuses to fails of a test output
func Print(w io.Writer, r io.Reader) error {
	var subMode, skip, found bool
	var line, prevLine string

	s := bufio.NewScanner(r)
	fmt.Fprintln(w, "Searching for the fail logs..")
	for s.Scan() {

		prevLine = line
		line = s.Text()
		if strings.HasPrefix(line, failIndicator) {
			fmt.Fprintln(w, line)
			found = true
			subMode = true
			continue
		}

		// no indentation means we returned to the top level
		if !strings.HasPrefix(line, indent) {
			subMode = false
		}

		if subMode {
			// if there is a fail remove skip indicator
			if strings.Contains(line, failIndicator) {
				skip = false
				continue
			}

			if skip {
				continue
			}

			// don't print the passed tests
			if strings.Contains(line, "--- PASS:") && !skip {
				count := strings.Count(line, indent)
				fmt.Fprintln(w, strings.Repeat(indent, count)+"...")
				skip = true
				continue
			}

			// this should print the actual test that have failed
			if strings.Contains(prevLine, failIndicator) {
				fmt.Fprintln(w, prevLine)
			}
			// this should print the error messages
			fmt.Fprintln(w, line)
		}
	}
	if !found {
		fmt.Fprintln(w, "Could not find any fail logs..")
	}
	return nil
}
