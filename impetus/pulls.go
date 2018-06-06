package impetus

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type clock interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

type realtime struct{}

func (realtime) Now() time.Time {
	return time.Now()
}

func (realtime) Since(t time.Time) time.Duration {
	return time.Now().Sub(t)
}

type pullScanner struct {
	gh                 githubApiServicer
	IdleToleranceHours time.Duration
	clock              clock
}

func (s *pullScanner) getIdleOpenPullsOnOrg(org string) []*github.PullRequest {
	var collectiveIdlePulls []*github.PullRequest

	repositories, err := s.gh.getRepositoriesOnOrg(org)
	if err != nil {
		log.Printf("Uhoh, couldn't get repositories for Organization %q because %q", org, err)
	}

	for _, repo := range repositories {
		repoIdlePulls := s.getIdleOpenPulls(org, *repo.Name)
		if len(repoIdlePulls) > 0 {
			collectiveIdlePulls = append(repoIdlePulls)
		}
	}

	return collectiveIdlePulls
}

func (s *pullScanner) getIdleOpenPulls(org, repository string) []*github.PullRequest {
	pulls, err := s.gh.getOpenPulls(org, repository)
	if err != nil {
		log.Printf("Uhoh, couldn't get pull requests because %q", err)
	}
	var idlePulls []*github.PullRequest

	for _, pull := range pulls {
		durationSinceLastUpdate := s.clock.Since(*pull.UpdatedAt)
		if durationSinceLastUpdate > s.IdleToleranceHours*time.Hour {
			log.Printf("%s was found to be idle", *pull.URL)
			idlePulls = append(idlePulls, pull)
		}
	}

	return idlePulls
}

func NewGhePullScanner() *pullScanner {
	gheUrl := buildGithubEnterpriseClient()
	ghClientWrapper := githubApiClient{client: gheUrl}
	scanner := pullScanner{gh: ghClientWrapper, clock: realtime{}}
	return &scanner
}

func buildGithubEnterpriseClient() *github.Client {
	gheURL := fmt.Sprintf("%s/api/v3/", os.Getenv("GHE_BASE_URL"))
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("IMPETUS_GIT_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client, err := github.NewEnterpriseClient(gheURL, gheURL, tc)
	if err != nil {
		log.Printf("Something went wrong setting up the github client: %q", err)
	}
	return client
}
