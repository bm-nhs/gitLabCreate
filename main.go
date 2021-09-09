package main

import (
	"io/ioutil"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	//token := string(os.Getenv("githubPAT"))
	//targetOrg := string(os.Getenv("targetOrg"))
	//payload := string(os.Getenv("payload"))
	//branchName := string(os.Getenv("branchName"))
	rmFile, _ := ioutil.ReadFile("./.payUnload")
	rmFiles := strings.Split(string(rmFile), "\n")


	for i := 0; i < len(rmFiles); i++ {
		println(rmFiles[i])
	}
	/*
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	repos, _, _ := client.Repositories.ListByOrg(ctx, targetOrg, nil)

	for i := 0; i < len(repos); i++ {
		repo := *repos[i].Name
		url := *repos[i].CloneURL
		r, err := git.PlainCloneContext(ctx, repo, false, &git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: "2",
				Password: token,
			},
			URL: url,
		})
		if err == nil {
			target := "./" + repo + "/."

			copy.Copy(payload, target)
			w, _ := r.Worktree()
			w.Checkout(&git.CheckoutOptions{Branch: plumbing.ReferenceName(branchName)})
			w.Add(payload)
			w.Commit("Added Payload", &git.CommitOptions{})
			r.Push(&git.PushOptions{
				RemoteName: "origin",
				Auth: &http.BasicAuth{
					Username: "2",
					Password: token,
				},
			})
		}
	}
	*/
}
