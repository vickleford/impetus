package impetus

import (
	"os"

	"github.com/tbruyelle/hipchat-go/hipchat"
)

type hipchatRoomServicer interface {
	sendMessageToRoom(string) error
}

type hipchatRoomClient struct {
	roomID string
	client *hipchat.Client
}

func (c *hipchatRoomClient) sendMessageToRoom(message string) error {
	notifReq := &hipchat.NotificationRequest{Message: message}
	_, err := c.client.Room.Notification(c.roomID, notifReq)
	if err != nil {
		return err
	}
	return nil
}

func newHipchatRoomClient(roomID string) *hipchatRoomClient {
	token := os.Getenv("IMPETUS_HIPCHAT_ROOM_TOKEN")
	theirClient := hipchat.NewClient(token)
	wrappedClient := hipchatRoomClient{roomID: roomID, client: theirClient}
	return &wrappedClient
}
