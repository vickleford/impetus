package impetus

import (
	"os"
	"testing"
)

func TestGetRepositoriesOnOrgPaginates(t *testing.T) {
	if tok := os.Getenv("IMPETUS_GIT_TOKEN"); tok == "" {
		t.Skip("IMPETUS_GIT_TOKEN not set")
	}
	gheClient := buildGithubEnterpriseClient()
	ghe := githubApiClient{client: gheClient}
	repos, err := ghe.getRepositoriesOnOrg("hpid")
	if err != nil {
		t.Errorf("Error listing repos on org: %q", err)
	}
	if count := len(repos); count < 31 {
		t.Errorf("Did not seem to get enough repos. Got %d", count)
	}
}

func TestGetOpenPullsPaginates(t *testing.T) {
	if tok := os.Getenv("IMPETUS_GIT_TOKEN"); tok == "" {
		t.Skip("IMPETUS_GIT_TOKEN not set")
	}
	gheClient := buildGithubEnterpriseClient()
	ghe := githubApiClient{client: gheClient}
	pulls, err := ghe.getOpenPulls("hpid", "ping-topology")
	if err != nil {
		t.Errorf("Error getting pulls: %q", err)
	}
	if numpulls := len(pulls); numpulls < 5 {
		t.Errorf("Did not get seem to get enough pulls. Got %d", numpulls)
	}
}
