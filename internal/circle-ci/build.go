package circleci

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/isacikgoz/nutshell-ci/internal/fails"
)

// Build is the circle-ci job that is under investigation
type Build struct {
	path  string
	Steps []*Step `json:"steps"`
}

func (b *Build) GetFailingSteps() ([]*Step, error) {
	steps := make([]*Step, 0)
	for _, step := range b.Steps {
		if step.Actions[0].Failed {
			steps = append(steps, step)
		}
	}
	return steps, nil
}

// FindFails gets the output file of the an action
func (b *Build) FindFails(a *Action) error {
	url := b.path + "/output/" + strconv.Itoa(a.Step) + "/0?file=true&allocation-id=" + a.AllocationID
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("could not do request: %w", err)
	}
	defer res.Body.Close()

	err = fails.Print(os.Stdout, res.Body)
	if err != nil {
		return fmt.Errorf("could not print fails: %w", err)
	}

	return nil
}
