package repl

import (
	"fmt"
	"snap_client/models"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type Room struct {
	RoomID      uuid.UUID
	RoomName    string
	RoomType    string // private or group
	MessageType string
	Target      []string
	Message     []byte
}

func Prompt(auth models.Auth, time time.Time, room *Room, msg ...string) string {
	tsStr := time.Format("[15:04] ")
	return fmt.Sprintf("%s\033[32m»\033[0m [%s] [%s] %s", tsStr, getRoomName(room), auth.Name(), strings.Join(msg, " "))
}

func ClientInfoPrompt(auth models.Auth, room *Room, message string) string {
	tsStr := time.Now().Format("[15:04] ")
	return fmt.Sprintf(
		"%s\033[32m»\033[0m [%s] [%s] %s\n",
		tsStr,
		getRoomName(room),
		auth.Name(),
		message,
	)
}

func ClientErrorPrompt(auth models.Auth, room *Room, message string) string {
	tsStr := time.Now().Format("[15:04] ")
	return fmt.Sprintf(
		"%s\033[31m»\033[0m [%s] [%s] %s\n",
		tsStr,
		getRoomName(room),
		auth.Name(),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#f53911")).Render(message),
	)
}

func IncomingPromptWithSenderAndMessage(message models.ClientMessage) string {
	tsStr := time.Now().Format("[15:04] ")
	return fmt.Sprintf("%s\033[31m«\033[0m [%s - %s] %s\n", tsStr, message.RoomName, message.Sender.Name, message.Message)
}

func getRoomName(room *Room) string {
	if room == nil {
		return "Lobby"
	}

	return room.RoomName
}
