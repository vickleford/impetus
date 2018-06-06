package impetus

import (
	"testing"
)

type hipchatRoomMock struct {
	sentMessage      string
	messageSendError error
}

func (r *hipchatRoomMock) sendMessageToRoom(message string) error {
	r.sentMessage = message
	return r.messageSendError
}

func TestHipchatGetsNotifiedWhenRepositoryScanned(t *testing.T) {
	zeejarmo := newMockRepo("zee", "jarmo")
	zeejarmo.addPullRequest(5, "Jan 2, 2018 at 3:00pm (MST)")
	mockedghclient := newGithubClientMock(zeejarmo)
	mockedscanner := pullScanner{gh: mockedghclient, clock: mockTime{}}

	roomMock := hipchatRoomMock{}

	controller := controller{
		repositories: []string{"zee/jarmo"},
		roomID:       "293238238fff",
		prscanner:    &mockedscanner,
		hipchatroom:  &roomMock,
	}

	controller.ScanAndReport()

	if roomMock.sentMessage != "https://github.example.com/zee/jarmo/pulls/5 still needs a code review" {
		t.Errorf("Sent the wrong message to hipchat. Sent %q\n", roomMock.sentMessage)
	}
}

func TestHipchatGetsNotifiedWhenOrgScanned(t *testing.T) {
	thingerflopper := newMockRepo("thinger", "flopper")
	thingerflopper.addPullRequest(17, "Jan 2, 2018 at 3:00pm (MST)")
	thingerflopper.addPullRequest(18, "Jan 17, 2018 at 1:59pm (MST)")
	mockedghclient := newGithubClientMock(thingerflopper)
	mockedscanner := pullScanner{gh: mockedghclient, clock: mockTime{}, IdleToleranceHours: 24}

	roomMock := hipchatRoomMock{}

	controller := controller{
		organizations: []string{"thinger"},
		roomID:        "1111111111",
		prscanner:     &mockedscanner,
		hipchatroom:   &roomMock,
	}

	controller.ScanAndReport()

	if roomMock.sentMessage != "https://github.example.com/thinger/flopper/pulls/17 still needs a code review" {
		t.Errorf("Sent the wrong message to hipchat. Sent %q\n", roomMock.sentMessage)
	}
}
