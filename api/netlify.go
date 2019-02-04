package api

import (
	"bytes"
	"fmt"
	"time"
)

type BuildState string
type BuildContext string

const (
	StateReady    BuildState = "ready"
	StateBuilding BuildState = "building"
	StateError    BuildState = "error"
)

const (
	ContextBranchDeploy BuildContext = "branch-deploy"
	ContextProduction   BuildContext = "production"
)

type WebhookMsg struct {
	ID           string        `json:"id"`
	SiteID       string        `json:"site_id"`
	BuildID      string        `json:"build_id"`
	State        BuildState    `json:"state"`
	SiteName     string        `json:"name"`
	URL          string        `json:"url"`
	SSLURL       string        `json:"ssl_url"`
	DeployURL    string        `json:"deploy_url"`
	DeploySSLURL string        `json:"deploy_ssl_url"`
	DeployTime   *int          `json:"deploy_time"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	ErrorMessage string        `json:"error_message"`
	CommitRef    string        `json:"commit_ref"`
	Branch       string        `json:"branch"`
	CommitURL    string        `json:"commit_url"`
	Title        string        `json:"title"`
	Context      BuildContext  `json:"context"`
	Summary      DeploySummary `json:"summary"`
	Committer    string        `json:"committer"`
}

type DeploySummary struct {
	Status   string `json:"status"`
	Messages []struct {
		Type        string `json:"type"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Details     string `json:"details"`
	} `json:"messages"`
}

func (c DeploySummary) String() string {
	b := bytes.NewBuffer([]byte{})
	for _, m := range c.Messages {
		b.WriteString(fmt.Sprintf("%s (%s)\n", m.Title, m.Description))
		b.WriteString(m.Details)
		b.WriteRune('\n')
	}
	return b.String()
}
