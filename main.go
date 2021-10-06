package main

import (
	"context"
	"fmt"

	"goGitBack/copy"
	"goGitBack/github"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	gh "github.com/google/go-github/v38/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"

)

func main() {
	// Load env vars
	err := godotenv.Load(".env")
	if err != nil {
		println("failed to load .env file ... review README.MD and configure")
		err = nil
	}
	token := os.Getenv("githubPAT")
	targetOrg := os.Getenv("targetOrg")
	prDescription := os.Getenv("prDescription")
	branchName := os.Getenv("commitBranch")
	target := os.Getenv("target")
	payloadDir := os.Getenv("payloadDir")
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
	repositories, gitHubResponse, err := client.Repositories.ListByOrg(ctx, targetOrg, &repositoryListByOrgOptions)

	if err != nil {
		println("err with getting client")
		err = nil
	}

	//Paginate through repositories.
	for i := 1; i < gitHubResponse.LastPage; i++ {
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
			err = r.Storer.SetReference(ref)

			err = w.Checkout(&git.CheckoutOptions{
				Branch: ref.Name(),
			})
			err = copy.Copy(payloadDir, target)
			if err != nil {
				println("failed to copy payload to target repository make sure you are targeting the correct directory")
				err = nil
			}
			_, err = w.Add(payloadDir)
			if err != nil {
				println("failed to copy payload to target repository make sure you are targeting the correct directory")
				err = nil
			}
			_, err = w.Commit("Added Payload", &git.CommitOptions{
				All: true,
			})
			if err != nil {
				println("Failed to commit targeted payload")
				err = nil
			}
			err = r.Push(&git.PushOptions{
				RemoteName: "origin",
				Auth: &http.BasicAuth{
					Username: "2",
					Password: token,
				},
			})
			if err != nil {
				println("Failed to push to target repository")
				err = nil
			}

		} else { println("error with clone")}

		//make PR
		payload := github.CreatePullRequestPayload{
			Title: branchName,
			Head:  branchName,
			Base:  *repositories[i].DefaultBranch,
			Body: prDescription,
		}
		err = github.PullRequest(payload, targetOrg, repoName, token)
		if err != nil {
			println("error with pull")
			println(err)
		}
		//clean up
		//err = os.RemoveAll("./" + repoName)
		if err != nil {
			println("error cleaning up directories")
			println(err)
		}

	}
}

func branchTarget(bn string) string {
	return fmt.Sprintf("refs/heads/%s", bn)
}
