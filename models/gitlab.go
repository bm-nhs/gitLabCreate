package models

type gitLab struct {
	name              string
	namespace_id      string
	privateToken      string
}

type gitLabRepo struct {
	url   			  string
}

type gitOptions struct {
	gitTargetToBackup string
	gitBackupLocation string
	sshPrivateKey     string
	actionToInsert    string
}

type gitHubRepo struct {
	gitURL            string
	gitDisc           string

}

type gitOrg struct {
	gitOrgURL         string
	gitHubRepos       []gitHubRepo
}

type