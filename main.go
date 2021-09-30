package main

import (
	"context"
	"fmt"

	git "github.com/go-git/go-git/v5"
	plumbing "github.com/go-git/go-git/v5/plumbing"
	http "github.com/go-git/go-git/v5/plumbing/transport/http"
	gh "github.com/google/go-github/v38/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"os"
)

func main() {
	// Load env vars
	godotenv.Load(".env")
	token := string(os.Getenv("githubPAT"))
	targetOrg := string(os.Getenv("targetOrg"))

	branchName := string(os.Getenv("commitBranch"))


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

			w, _ := r.Worktree()
			headRef, _ := r.Head()
			bt := branchTarget(branchName)
			ref := plumbing.NewHashReference(plumbing.ReferenceName(bt), headRef.Hash())
			r.Storer.SetReference(ref)
			err = w.Checkout(&git.CheckoutOptions{
				Branch: ref.Name(),
			})
			r.Push(&git.PushOptions{
				RemoteName: "origin",
				Auth: &http.BasicAuth{
					Username: "2",
					Password: token,
				},
			})

		} else { println("error with clone")}
		// Clean up
		err = os.RemoveAll("./" + repoName)
	}
}

func branchTarget(bn string) string {
	return fmt.Sprintf("refs/heads/%s", bn)
}
