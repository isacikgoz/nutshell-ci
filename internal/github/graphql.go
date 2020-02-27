package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"golang.org/x/oauth2"
)

type query struct {
	Query string `json:"query"`
}

type response struct {
	Data json.RawMessage
}

type data struct {
	Repository struct {
		PullRequest struct {
			Oid  string `json:"headRefOid"`
			Head string `json:"headRefName"`
		} `json:"pullRequest"`
	} `json:"repository"`
}

// GetLatestCommitOfPR fetches the latest commit ID
func GetLatestCommitOfPR(ctx context.Context, token, user, repo string, num int) (string, string, error) {
	q := generateGraphQLQuery(user, repo, num)
	qr := query{Query: q}

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(qr)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", &buf)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/json")

	src := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)
	client := oauth2.NewClient(ctx, src)

	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return "", "", fmt.Errorf("non-OK status code: %v body: %q", res.Status, body)
	}
	defer res.Body.Close()

	var out response
	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		return "", "", err
	}

	var d data
	err = json.Unmarshal(out.Data, &d)
	if err != nil {
		return "", "", err
	}

	return d.Repository.PullRequest.Oid, d.Repository.PullRequest.Head, nil
}

func generateGraphQLQuery(user, repo string, id int) string {
	vars := struct {
		Owner, Repository string
		ID                int
	}{
		Owner:      user,
		Repository: repo,
		ID:         id,
	}
	tpl := template.Must(template.New("pull").Parse(queryTpl))

	w := bytes.NewBuffer(make([]byte, 0))
	tpl.Execute(w, vars)

	return fmt.Sprint(w)
}
