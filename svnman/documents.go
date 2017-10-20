package svnman

// CreateRepo is contains the info required to create a repository.
type CreateRepo struct {
	RepoID              string `json:"repo_id"`
	AuthenticationRealm string `json:"auth_realm"`
	ProjectID           string `json:"project_id"`
	Creator             struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	} `json:"creator"`
}
