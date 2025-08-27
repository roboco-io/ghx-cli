package graphql

// Repository queries and mutations

// RepositoryQuery represents the repository query structure
type RepositoryQuery struct {
	Repository *RepositoryDetails `graphql:"repository(owner: $owner, name: $name)"`
}

// RepositoryDetails represents a GitHub repository
type RepositoryDetails struct {
	ID          string          `graphql:"id"`
	Name        string          `graphql:"name"`
	Owner       RepositoryOwner `graphql:"owner"`
	IsPrivate   bool            `graphql:"isPrivate"`
	Visibility  string          `graphql:"visibility"`
	Description *string         `graphql:"description"`
	URL         string          `graphql:"url"`
}

// RepositoryOwner represents a repository owner (User or Organization)
type RepositoryOwner struct {
	Login string `graphql:"login"`
	Type  string `graphql:"__typename"`
}

// RepositoryInfo represents simplified repository information
type RepositoryInfo struct {
	ID          string
	Name        string
	Owner       string
	OwnerType   string
	IsPrivate   bool
	Visibility  string
	Description *string
	URL         string
}

// BuildRepositoryVariables builds variables for repository queries
func BuildRepositoryVariables(owner, name string) map[string]interface{} {
	return map[string]interface{}{
		"owner": owner,
		"name":  name,
	}
}

// ParseRepositoryResponse parses repository query response
func ParseRepositoryResponse(resp *RepositoryQuery) (*RepositoryInfo, error) {
	if resp.Repository == nil {
		return nil, nil // Repository not found or inaccessible
	}

	repo := resp.Repository
	return &RepositoryInfo{
		ID:          repo.ID,
		Name:        repo.Name,
		Owner:       repo.Owner.Login,
		OwnerType:   repo.Owner.Type,
		IsPrivate:   repo.IsPrivate,
		Visibility:  repo.Visibility,
		Description: repo.Description,
		URL:         repo.URL,
	}, nil
}
