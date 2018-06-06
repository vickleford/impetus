package impetus

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-github/github"
)

const longTimeFormat = "Jan 2, 2006 at 3:04pm (MST)"

type mockRepo struct {
	organization string
	repository   *github.Repository
	pullrequests []*github.PullRequest
}

func (p *mockRepo) addPullRequest(prnumber int, updatedAt string) *github.PullRequest {
	t, err := time.Parse(longTimeFormat, updatedAt)
	if err != nil {
		fmt.Printf("t is a %v and the problem was: %q\n", t, err)
	}
	prurl := fmt.Sprintf("https://github.example.com/%s/%s/pulls/%d", p.organization, *p.repository.Name, prnumber)
	pr := github.PullRequest{
		Number:    &prnumber,
		UpdatedAt: &t,
		URL:       &prurl,
	}
	p.pullrequests = append(p.pullrequests, &pr)
	return &pr
}

func newMockRepo(org, name string) *mockRepo {
	repository := github.Repository{
		Name: &name,
	}
	mrepository := mockRepo{
		organization: org,
		repository:   &repository,
	}
	return &mrepository
}

type githubClientMock struct {
	scannedOrganization string
	scannedRepository   string
	repositories        []*mockRepo
	err                 error
}

func newGithubClientMock(repositories ...*mockRepo) *githubClientMock {
	var mock githubClientMock
	for _, repo := range repositories {
		mock.repositories = append(mock.repositories, repo)
	}
	return &mock
}

func (c *githubClientMock) getOpenPulls(org, repo string) ([]*github.PullRequest, error) {
	c.scannedOrganization = org
	c.scannedRepository = repo
	var returnable []*github.PullRequest
	for _, proj := range c.repositories {
		if proj.organization != org {
			continue
		}
		if *proj.repository.Name != repo {
			continue
		}
		returnable = proj.pullrequests
	}
	return returnable, c.err
}

func (c *githubClientMock) getRepositoriesOnOrg(org string) ([]*github.Repository, error) {
	c.scannedOrganization = org
	var returnable []*github.Repository
	for _, proj := range c.repositories {
		if proj.organization != org {
			continue
		}
		returnable = append(returnable, proj.repository)
	}
	return returnable, c.err
}

type mockTime struct {
	time.Time
}

func (m mockTime) Now() time.Time {
	t, _ := time.Parse(longTimeFormat, "Jan 17, 2018 at 2:00pm (MST)")
	return t
}

func (m mockTime) Since(t time.Time) time.Duration {
	return m.Now().Sub(t)
}

func TestIdlePrsOnSingleRepositoryAreReturned(t *testing.T) {
	barRepository := newMockRepo("foo", "bar")
	barRepository.addPullRequest(32, "Jan 2, 2018 at 3:04pm (MST)")
	barRepository.addPullRequest(33, "Jan 17, 2018 at 1:00pm (MST)")
	mockedghclient := newGithubClientMock(barRepository)

	mockTime := mockTime{}

	finder := pullScanner{gh: mockedghclient, clock: mockTime, IdleToleranceHours: 24}
	idlePrs := finder.getIdleOpenPulls("foo", "bar")
	if *idlePrs[0].Number != 32 {
		t.Errorf("man i got the wrong prs back: %v\n", idlePrs)
	}
	if l := len(idlePrs); l != 1 {
		t.Errorf("man i got too many prs back: %v\n", idlePrs)
	}
	if mockedghclient.scannedOrganization != "foo" {
		t.Errorf("scanned the wrong org: scanned %q\n", mockedghclient.scannedOrganization)
	}
	if mockedghclient.scannedRepository != "bar" {
		t.Errorf("scanned the wrong repository: scanned %q\n", mockedghclient.scannedRepository)
	}
}

func TestIdlePrsOnOrgAreReturned(t *testing.T) {
	notMyRepository := newMockRepo("theirorganization", "theirrepo")
	testrepo := newMockRepo("testorganization", "testrepo")
	testrepo.addPullRequest(32, "Jan 2, 2018 at 3:04pm (MST)")
	testrepo.addPullRequest(33, "Jan 17, 2018 at 1:08pm (MST)")
	mockedghclient := newGithubClientMock(testrepo, notMyRepository)

	mockTime := mockTime{}

	finder := pullScanner{gh: mockedghclient, clock: mockTime, IdleToleranceHours: 24}
	idlePrs := finder.getIdleOpenPullsOnOrg("testorganization")
	if *idlePrs[0].Number != 32 {
		t.Errorf("man i got the wrong prs back: %v\n", idlePrs)
	}
	if mockedghclient.scannedOrganization != "testorganization" {
		t.Errorf("scanned the wrong org: scanned %q\n", mockedghclient.scannedOrganization)
	}
}
