package main

import (
	"context"
	"fmt"
	copy "goGitBack/copy"
	"goGitBack/github"
	"os"

	git "github.com/go-git/go-git/v5"
	plumbing "github.com/go-git/go-git/v5/plumbing"
	http "github.com/go-git/go-git/v5/plumbing/transport/http"
	gh "github.com/google/go-github/v38/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	godotenv.Load(".env")
	token := string(os.Getenv("githubPAT"))
	targetOrg := string(os.Getenv("targetOrg"))
	payload := string(os.Getenv("payload"))
	branchName := string(os.Getenv("branchName"))

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := gh.NewClient(tc)
	repositories, _, _ := client.Repositories.ListByOrg(ctx, targetOrg, nil)

	for i := 0; i < len(repositories); i++ {
		repo := *repositories[i].Name
		url := *repositories[i].CloneURL
		r, err := git.PlainCloneContext(ctx, repo, false, &git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: "2",
				Password: token,
			},
			URL: url,
		})
		if err == nil {
			target := "./" + repo + "/."
			w, _ := r.Worktree()
			headRef, _ := r.Head()

			ref := plumbing.NewHashReference(plumbing.ReferenceName(branchTarget(branchName)), headRef.Hash())
			r.Storer.SetReference(ref)
			err = w.Checkout(&git.CheckoutOptions{
				Branch: ref.Name(),
			})
			copy.Copy(payload, target)
			w.Add(payload)
			w.Commit("Added Payload", &git.CommitOptions{
				All: true,
			})
			r.Push(&git.PushOptions{
				RemoteName: "origin",
				Auth: &http.BasicAuth{
					Username: "2",
					Password: token,
				},
			})
			if err != nil {
				println(err)
			}
		}
		//make PR
		payload := github.CreatePullRequestPayload{
			Title: branchName,
			Head:  branchName,
			Base:  "master",
		}
		github.PullRequest(payload, targetOrg, repo)
		//clean up
		err = os.RemoveAll("./" + repo)
		if err != nil {
			println(err)
		}
	}
}

func branchTarget(branchName string) string {
	return fmt.Sprintf("refs/heads/%s", branchName)
}
