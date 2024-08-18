package commands

import (
	"encoding/json"
	"errors"
	"snap_client/models"
	"snap_client/repl"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gopkg.in/readline.v1"
)

func GroupChat(websocketConn *websocket.Conn, rl *readline.Instance, line string, authSession models.Auth, httpAddr string, room *repl.Room) (*repl.Room, error) {
	splitContent := strings.Split(line, " ")
	var groupKey string

	if len(splitContent) < 2 {
		return nil, errors.New("Wrong command format. Valid format must be : /gm <group_key>")
	}

	_, groupKey = splitContent[0], splitContent[1]

	if len(splitContent) > 2 {

		var message = strings.Join(splitContent[2:], " ")

		msgByte, err := json.Marshal(repl.Room{
			RoomName:    "Group",
			RoomType:    models.MT_GroupChat, // private or group
			MessageType: models.MT_GroupChat, // private or group
			Target:      []string{groupKey},
			Message:     []byte(message),
		})

		if err != nil {
			return nil, err
		}
		err = websocketConn.WriteMessage(websocket.TextMessage, msgByte)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	return &repl.Room{
		RoomID:      uuid.New(),
		RoomName:    "Group",
		RoomType:    models.MT_GroupChat, // private or group
		MessageType: models.MT_GroupChat, // private or group
		Target:      []string{groupKey},
	}, nil
}
