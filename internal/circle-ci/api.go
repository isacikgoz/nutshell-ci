package circleci

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Step is a part of the Build
type Step struct {
	Name    string    `json:"name"`
	Actions []*Action `json:"actions"`
}

// GetBuild fetches the result of a circle-ci build
func GetBuild(ctx context.Context, token, path string) (*Build, error) {
	path = strings.Replace(path, "circleci.com/", "circleci.com/api/v1.1/project/", 1)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	req.SetBasicAuth(token, "")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not do the request: %w", err)
	}
	defer res.Body.Close()

	b := new(Build)
	err = json.NewDecoder(res.Body).Decode(b)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %w", err)
	}
	b.path = path
	return b, nil
}
