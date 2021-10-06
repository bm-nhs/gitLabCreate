package main

import (
	"context"
	"fmt"

	copy "goGitBack/copy"
	github "goGitBack/github"
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
	prDescription := string(os.Getenv("prDescription"))
	branchName := string(os.Getenv("commitBranch"))
	target := string(os.Getenv("target"))
	payloadDir := string(os.Getenv("payloadDir"))
	// Initialize oauth connection so we can grab a list of all repos in target org
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := gh.NewClient(tc)

	// Initialize options for pagination
	var listOptions = gh.ListOptions{
		Page: 1,
		PerPage: 100,
	}
	repositoryListByOrgOptions := gh.RepositoryListByOrgOptions{
		Type: "all",
		ListOptions: listOptions,
	}
	repositories, githHubResponse, err := client.Repositories.ListByOrg(ctx, targetOrg, &repositoryListByOrgOptions)

	if err != nil {
		println("err with getting client")
		err = nil
	}

	//Paginate through repositories.
	for i := 1; i < githHubResponse.LastPage; i++ {
		listOptions.Page = i
		repositoryListByOrgOptions.ListOptions = listOptions
		pagination, _, err := client.Repositories.ListByOrg(ctx, targetOrg, &repositoryListByOrgOptions)
		if err != nil {
			println(err)
			return
		}
		repositories = append(repositories,pagination...)
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
			copy.Copy(payloadDir, target)
			w.Add(payloadDir)
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
		//err = os.RemoveAll("./" + repoName)
		if err != nil {
			println("error with pull")
			println(err)
		}

	}
}

func branchTarget(bn string) string {
	return fmt.Sprintf("refs/heads/%s", bn)
}
