package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	circleci "github.com/isacikgoz/wtf-ci/internal/circle-ci"
	"github.com/isacikgoz/wtf-ci/internal/github"
	mmjenkins "github.com/isacikgoz/wtf-ci/internal/mm-jenkins"
)

func main() {

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

	links, err := github.GetCircleCIFails(ctx, oid)
	if err != nil {
		return err
	}
	for _, link := range links {
		fmt.Printf("Found a failing check at %s ❌\n", link)
		if err := readCircleCILogs(ctx, link); err != nil {
			//return err
		}
	}
	if err := readJenkinsLogs(ctx, branch, pr); err != nil {
		//return err
	}
	fmt.Printf("Done. 🎉\n")
	return nil
}

func readCircleCILogs(ctx context.Context, path string) error {

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	build, err := circleci.GetBuild(context.Background(), os.Getenv("CIRCLECI_TOKEN"), path)
	if err != nil {
		return err
	}

	steps, _ := build.GetFailingSteps()
	for _, step := range steps {
		act := step.Actions[0]
		fmt.Printf("The step %q took %s to finish ⏱\n", act.Name, act.Duration())

		err = build.FindFails(act)
		if err != nil {
			return err
		}
	}
	return nil
}

func readJenkinsLogs(ctx context.Context, branch string, pr int) error {
	err := mmjenkins.FindFails(ctx, branch, fmt.Sprintf("PR-%d", pr))
	if err != nil {
		return err
	}
	return nil
}
