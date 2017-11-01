package svnman

// CreateRepo is contains the info required to create a repository.
type CreateRepo struct {
	RepoID    string `json:"repo_id"`
	ProjectID string `json:"project_id"`
	Creator   string `json:"creator"` // Full Name <email> notation.
}

// ModifyAccess contains the changes in access rules for users of a specific repository.
type ModifyAccess struct {
	Grant []struct {
		Username string `json:"username"`
		Password string `json:"password"` // always bcrypted and base64-encoded.
	} `json:"grant"`
	Revoke []string `json:"revoke"` // list of usernames
}
