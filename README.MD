# gogi (GOGIT)
This repository bulk adds a payload into all repositories owned by a user or organization. 


# [How to set up GoLang](https://www.digitalocean.com/community/tutorials/how-to-build-and-install-go-programs)
## Query GitHub API for all repositories within an organization:
```
curl "https://api.github.com/users/$GHUSER/repos?access_token=$GITHUB_API_TOKEN" | grep -w clone_url
```

## Create repository in GitLab with GitLab API V4
```
curl -H "Content-Type:application/json" https://gitlab.com/api/v4/projects?private_token=TOKEN -d "{ \"name\": \"newRepo\", \"namespace_id\": \"ID\" }"            
```

namespace_id = user or org name for GitLab
name = name of repository
private_token= gitLab private token

TOKEN = gitLab personal access token
###### (User Settings > Access Tokens)
ID    = gitLab group namespace_id
###### request all groups associated with account to get namespace id namespace.id
```
curl --header "PRIVATE-TOKEN: <token>" "https://gitlab.com/api/v4/namespaces"    
```


##### setting up .env vars for testing
``` # github personal access token set to read private repositories and write
githubPAT = <ghp>
targetOrg = <Org Name Example: nhsllc>
```
### Enviroment Vars
# export githubPAT = gitHub-auth-token
# targetOrg = org-alias
# payload = ./path/to/payload
# branchName = name-of-new-branch
