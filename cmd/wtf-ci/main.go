package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	circle "github.com/isacikgoz/wtf-ci/internal/circle-ci"
	"github.com/isacikgoz/wtf-ci/internal/fails"
	"github.com/isacikgoz/wtf-ci/internal/github"
	jenkins "github.com/isacikgoz/wtf-ci/internal/mm-jenkins"
)

func main() {

	if len(os.Args) <= 1 {
		if err := fails.Print(os.Stdout, os.Stdin); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	pr := os.Args[1]
	s := strings.Split(pr, "/")
	id, err := strconv.Atoi(s[len(s)-1])
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	if err := run(ctx, s[len(s)-4], s[len(s)-3], id); err != nil {
		fmt.Fprintf(os.Stderr, "program exited with error: %s", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, owner, repo string, pr int) error {
	fmt.Printf("Fetching the latest commit of PR#%d..\n", pr)
	oid, branch, err := github.GetLatestCommitOfPR(
		context.Background(),
		os.Getenv("GITHUB_TOKEN"),
		owner,
		repo,
		pr,
	)
	if err != nil {
		return err
	}
	fmt.Printf("Finding the failing check of commit %s..\n", oid[:7])

	readCircleCILogs(ctx, oid)
	readJenkinsLogs(ctx, branch, pr)

	fmt.Printf("Done. 🎉\n")
	return nil
}

func readCircleCILogs(ctx context.Context, oid string) error {

	links, err := github.GetCircleCIFails(ctx, oid)
	if err != nil {
		return err
	}
	for _, link := range links {
		fmt.Printf("Found a failing check at %s ❌\n", link)
		build, err := circle.GetBuild(context.Background(), os.Getenv("CIRCLECI_TOKEN"), link)
		if err != nil {
			return err
		}

		steps, _ := build.GetFailingSteps()
		for _, step := range steps {
			act := step.Actions[0]

			err = build.FindFails(act)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func readJenkinsLogs(ctx context.Context, branch string, pr int) error {
	build, err := jenkins.GetBuild(ctx, branch, fmt.Sprintf("PR-%d", pr))
	if err != nil {
		return err
	}

	for _, step := range build.Steps {
		for _, node := range step.Fails {
			fmt.Printf("Found a failing check at %s ❌\n", strings.TrimSuffix(step.Link, "/runs/1/nodes"))
			err = node.PrintFail()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
