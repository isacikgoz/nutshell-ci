package mmjenkins

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/isacikgoz/wtf-ci/internal/output"
)

// kv-set-improvements-5.21/runs/1/nodes/
// PR-13942/runs/1/nodes/

const jenkinsLink = "https://build.mattermost.com/blue/rest/organizations/jenkins/pipelines/mp/pipelines/mattermost-server-pr-new/branches/"

type nodes struct {
	Values []node
}

type node struct {
	Result      string `json:"result"`
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

func FindFails(ctx context.Context, branch, pr string) error {

	branch = jenkinsLink + branch + "/runs/1/nodes"
	pr = jenkinsLink + pr + "/runs/1/nodes/"

	logURLs, err := getLog(ctx, branch)
	if err != nil {
		return err
	}
	for _, logURL := range logURLs {
		fmt.Printf("Found a failing check at %s ❌\n", logURL)
		if err := printFail(logURL); err != nil {
			return err
		}
	}

	logURLs, err = getLog(ctx, pr)
	for _, logURL := range logURLs {
		fmt.Printf("Found a failing check at %s ❌\n", logURL)
		if err := printFail(logURL); err != nil {
			return err
		}
	}
	return nil
}

func printFail(url string) error {
	cmd := exec.Command("curl", url)
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("could not pipe: %w", err)
	}
	cmd.Start()
	go func() {
		err = output.PrintFails(os.Stdout, reader)
		if err != nil {
			fmt.Printf("could not print fails: %s\n", err)
		}
	}()
	cmd.Wait()
	reader.Close()
	return nil
}

func getLog(ctx context.Context, url string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data := []*node{}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	links := make([]string, 0)
	for _, node := range data {
		if node.Result == "FAILURE" {
			link := url + "/" + node.ID + "/log/"
			links = append(links, link)
		}
	}
	return links, nil
}
