package impetus

import (
	"context"

	"github.com/google/go-github/github"
)

type githubApiServicer interface {
	getOpenPulls(string, string) ([]*github.PullRequest, error)
	getRepositoriesOnOrg(string) ([]*github.Repository, error)
}

type githubApiClient struct {
	client *github.Client
}

func (c githubApiClient) getOpenPulls(org, repo string) ([]*github.PullRequest, error) {
	var allpulls []*github.PullRequest
	options := &github.PullRequestListOptions{State: "open"}

	for {
		pulls, resp, err := c.client.PullRequests.List(
			context.Background(), org, repo, options)
		if err != nil {
			return allpulls, err
		}
		allpulls = append(allpulls, pulls...)
		if resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}

	return allpulls, nil
}

func (c githubApiClient) getRepositoriesOnOrg(org string) ([]*github.Repository, error) {
	var allrepos []*github.Repository
	opt := &github.RepositoryListByOrgOptions{}

	for {
		repos, resp, err := c.client.Repositories.ListByOrg(
			context.Background(), org, opt)
		if err != nil {
			return allrepos, err
		}
		allrepos = append(allrepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allrepos, nil
}
