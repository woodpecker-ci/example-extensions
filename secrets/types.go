package secrets

type Secret struct {
	ID     int64    `json:"id"`
	OrgID  int64    `json:"org_id"`
	RepoID int64    `json:"repo_id"`
	Name   string   `json:"name"`
	Value  string   `json:"value,omitempty"`
	Images []string `json:"images"`
	Events []string `json:"events"`
}
