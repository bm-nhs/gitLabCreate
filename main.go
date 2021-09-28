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
	// Load env vars
	godotenv.Load(".env")
	token := string(os.Getenv("githubPAT"))
	targetOrg := string(os.Getenv("targetOrg"))
	payload := string(os.Getenv("payload"))
	branchName := string(os.Getenv("commitBranch"))
	prDescription := string(os.Getenv("prDescription"))

	// Initialize oauth connection so we can grab a list of all repos in target org
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := gh.NewClient(tc)
	repositories, _, err := client.Repositories.ListByOrg(ctx, targetOrg, nil)

	if err != nil {
		println("err with getting client")
		err = nil
	}

	// For each repo within a target organization or user targeted in GitHub
	// 1) Clone repo
	// 2) Create a new Branch
	// 3) Copy payload folder into project
	// 4) Commit payload to branch
	// 5) Push new branch to origin
	// 6) Create Pull Request to default branch using new branch
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
			bt := branchTarget(branchName)
			ref := plumbing.NewHashReference(plumbing.ReferenceName(bt), headRef.Hash())
			r.Storer.SetReference(ref)
			err = w.Checkout(&git.CheckoutOptions{
				Branch: ref.Name(),
			})
			if err != nil {
				println(branchTarget(*repositories[i].DefaultBranch))
				println(*repositories[i].DefaultBranch)
				println(ref.Name())
				println("err with checkout")
				err = nil
			}
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
		} else { println("error with clone")}
		//make PR
		payload := github.CreatePullRequestPayload{
			Title: branchName,
			Head:  branchName,
			Base:  *repositories[i].DefaultBranch,
			Body: prDescription,
		}
		github.PullRequest(payload, targetOrg, repoName, token)
		//clean up
		err = os.RemoveAll("./" + repoName)
		if err != nil {
			println("error with pull")
			println(err)
		}
	}
}

func branchTarget(bn string) string {
	return fmt.Sprintf("refs/heads/%s", bn)
}
