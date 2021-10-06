# gogi (GOGIT)
This repository bulk adds a payload into all repositories owned by a user or organization.

# [How to set up GoLang](https://www.digitalocean.com/community/tutorials/how-to-build-and-install-go-programs)
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