package github

import (
	"context"
	"io/ioutil"
	"net/http"
	"regexp"
)

// GetCircleCIFails looks for failed circle-ci builds
func GetCircleCIFails(ctx context.Context, oid string) ([]string, error) {
	url := "https://api.github.com/repos/mattermost/mattermost-server/commits/" + oid + "/check-runs"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.antiope-preview+json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`\[(\S+)\]\(((\S+)\?\S+)\) - Failed`)
	submatches := re.FindAllSubmatch(data, -1)
	links := make([]string, 0)
	for _, sub := range submatches {
		if len(sub) < 4 {
			continue
		}
		links = append(links, string(sub[3]))
	}

	return links, nil
}
