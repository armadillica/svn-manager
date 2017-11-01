package svnman

// CreateRepo is contains the info required to create a repository.
type CreateRepo struct {
	RepoID    string `json:"repo_id"`
	ProjectID string `json:"project_id"`
	Creator   string `json:"creator"` // Full Name <email> notation.
}

// ModifyAccessGrantEntry contains info about one user to allow access.
type ModifyAccessGrantEntry struct {
	Username string `json:"username"`
	Password string `json:"password"` // always bcrypted and base64-encoded.
}

// ModifyAccess contains the changes in access rules for users of a specific repository.
type ModifyAccess struct {
	Grant  []ModifyAccessGrantEntry `json:"grant"`
	Revoke []string                 `json:"revoke"` // list of usernames
}
