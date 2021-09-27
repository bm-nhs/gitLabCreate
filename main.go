package main

import (
	"context"
	"fmt"
	copy "gogi/copy"
	"gogi/github"
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
	commitBranch := string(os.Getenv("branchName"))
	prDescription := string(os.Getenv("prDescription"))

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := gh.NewClient(tc)
	repositories, _, _ := client.Repositories.ListByOrg(ctx, targetOrg, nil)

	for i := 0; i < len(repositories); i++ {
		repoName := *repositories[i].Name
		url := *repositories[i].CloneURL
		r, err := git.PlainCloneContext(ctx, repoName, false, &git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: "2",
				Password: token,
			},
			URL: url,
		})
		if err == nil {
			target := "./" + repoName + "/."
			w, _ := r.Worktree()
			headRef, _ := r.Head()

			ref := plumbing.NewHashReference(plumbing.ReferenceName(branchTarget(commitBranch)), headRef.Hash())
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
				println("err with checkout")
				println(w.Status())
				println(err)
				err = nil
			}
		}
		//make PR
		payload := github.CreatePullRequestPayload{
			Title: commitBranch,
			Head:  commitBranch,
			Base:  *repositories[i].DefaultBranch,
			Body: prDescription,
		}
		err = github.PullRequest(payload, targetOrg, repoName, token)
		if err != nil {
			println("error with PR")
			println(err)
			err = nil
		}
		//clean up
		err = os.RemoveAll("./" + repoName)
		if err != nil {
			println("err with dir cleanup")
			println(err)
			err = nil
		}
	}
}

func branchTarget(branchName string) string {
	return fmt.Sprintf("refs/heads/%s", branchName)
}
