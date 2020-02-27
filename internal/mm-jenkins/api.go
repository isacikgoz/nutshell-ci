package mmjenkins

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/isacikgoz/nutshell-ci/internal/fails"
)

const jenkinsLink = "https://build.mattermost.com/blue/rest/organizations/jenkins/pipelines/mp/pipelines/mattermost-server-pr-new/branches/"

type Node struct {
	Result      string `json:"result"`
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	link        string
}

type Build struct {
	Steps []*Step
}

type Step struct {
	Link  string
	Fails []*Node
}

func GetBuild(ctx context.Context, branch, pr string) (*Build, error) {
	branch = jenkinsLink + branch + "/runs/1/nodes"
	pr = jenkinsLink + pr + "/runs/1/nodes/"

	b := &Build{
		Steps: make([]*Step, 0),
	}

	b.Steps = append(b.Steps, &Step{Link: branch})
	b.Steps = append(b.Steps, &Step{Link: pr})

	for _, step := range b.Steps {
		step.getFailingNodes(ctx)
	}

	return b, nil
}

func (n *Node) PrintFail() error {
	url := n.link + "/" + n.ID + "/log/"
	cmd := exec.Command("curl", url)
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("could not pipe: %w", err)
	}
	cmd.Start()
	go func() {
		err = fails.Print(os.Stdout, reader)
		if err != nil {
			fmt.Printf("could not print fails: %s\n", err)
		}
	}()
	cmd.Wait()
	reader.Close()

	return nil
}

func (s *Step) getFailingNodes(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", s.Link, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data := []*Node{}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return err
	}

	nodes := make([]*Node, 0)
	for _, node := range data {
		if node.Result == "FAILURE" {
			node.link = s.Link
			nodes = append(nodes, node)
		}
	}
	s.Fails = nodes
	return nil
}
