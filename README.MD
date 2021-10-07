# gogi (GOGIT)
This repository bulk adds a payload into all repositories owned by a user or organization.

# [How to set up GoLang](https://www.digitalocean.com/community/tutorials/how-to-build-and-install-go-programs)
#### go version: go.1.17.1
#To install golang with WSL  
```wget https://golang.org/dl/go1.17.1.linux-amd64.tar.gz```
```sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.17.1.linux-amd64.tar.gz```
```export PATH=$PATH:/usr/local/go/bin```
## Query GitHub API for all repositories within an organization:
```
curl "https://api.github.com/users/$GHUSER/repos?access_token=$GITHUB_API_TOKEN" | grep -w clone_url
```

### Enviroment Vars
# export githubPAT = gitHub-auth-token
# targetOrg = org-alias
# payloadDir = ./path/to/payload
# commitBranch = name-of-new-branch
# prDescription = "Body of the PR request. Put something descriptive here"
