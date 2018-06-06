package impetus

import (
	"fmt"
	"log"
	"strings"
)

type controller struct {
	repositories  []string
	organizations []string
	roomID        string
	prscanner     *pullScanner
	hipchatroom   hipchatRoomServicer
}

func (c *controller) ScanAndReport() {
	for _, repo := range c.repositories {
		o, p := splitOrgAndRepo(repo)
		idlePrs := c.prscanner.getIdleOpenPulls(o, p)
		for _, idlePr := range idlePrs {
			message := fmt.Sprintf("%s still needs a code review", *idlePr.URL)
			err := c.hipchatroom.sendMessageToRoom(message)
			if err != nil {
				log.Printf("Could not send message to room: %q", err)
			}
		}
	}

	for _, org := range c.organizations {
		idlePrs := c.prscanner.getIdleOpenPullsOnOrg(org)
		for _, idlePr := range idlePrs {
			message := fmt.Sprintf("%s still needs a code review", *idlePr.URL)
			err := c.hipchatroom.sendMessageToRoom(message)
			if err != nil {
				log.Printf("Could not send message to room: %q", err)
			}
		}
	}
}

func splitOrgAndRepo(orgrepo string) (string, string) {
	parts := strings.Split(orgrepo, "/")
	return parts[0], parts[1]
}

func NewController(hipchatRoomID string, githubRepositories, githubOrganizations []string, scanner *pullScanner) *controller {
	hipchatroom := newHipchatRoomClient(hipchatRoomID)
	c := controller{
		repositories:  githubRepositories,
		organizations: githubOrganizations,
		roomID:        hipchatRoomID,
		prscanner:     scanner,
		hipchatroom:   hipchatroom,
	}
	return &c
}
