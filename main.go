package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"snap_client/chat/commands"
	"snap_client/forms"
	"snap_client/models"
	"snap_client/repl"
	"snap_client/utils"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
	"gopkg.in/readline.v1"
)

func main() {
	flag.Parse()

	addr := os.Args[len(os.Args)-1]

	var wsAddr string
	var wsAddrPath string = "/ws"
	var httpAddr string

	var websocketConn *websocket.Conn
	var err error
	var authSession models.Auth
	var room *repl.Room

	if !strings.HasPrefix(addr, "ws://") && !strings.HasPrefix(addr, "wss://") {
		if strings.HasPrefix(addr, "http://") {
			wsAddr = "ws://" + strings.TrimLeft(addr, "http://") + wsAddrPath
		} else if strings.HasPrefix(addr, "https://") {
			wsAddr = "wss://" + strings.TrimLeft(addr, "https://") + wsAddrPath
		} else {
			wsAddr = "ws://" + addr + wsAddrPath
		}
	} else {
		wsAddr = addr + wsAddrPath
	}

	if !strings.HasPrefix(addr, "http://") && !strings.HasPrefix(addr, "https://") {
		if strings.HasPrefix(addr, "ws://") {
			httpAddr = "http://" + strings.TrimLeft(addr, "ws://")
		} else if strings.HasPrefix(addr, "wss://") {
			httpAddr = "https://" + strings.TrimLeft(addr, "wss://")
		} else {
			httpAddr = "http://" + addr
		}
	} else {
		httpAddr = addr
	}

	request, err := http.NewRequest("POST", httpAddr+"/health", nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#f53911")).MarginTop(1).MarginBottom(1).Render(err.Error()))
		fmt.Println("Chat app successfully exited.")
		os.Exit(1)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#f53911")).MarginTop(1).MarginBottom(1).Render(err.Error()))
		fmt.Println("Chat app successfully exited.")
		os.Exit(1)
	}

	if !slices.Contains([]int{200, 201}, response.StatusCode) {
		fmt.Printf(lipgloss.NewStyle().Foreground(lipgloss.Color("#f53911")).MarginTop(1).MarginBottom(1).Render("X Chat server offline."))
		fmt.Println("Chat app successfully exited.")
		os.Exit(1)
	}

	var answer string

	prompt := &survey.Select{
		Message: "Login if you already have an account, or register for new one.",
		Options: []string{"Login", "Register", "Exit"},
		VimMode: true,
	}

	if err := survey.AskOne(prompt, &answer); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if strings.ToLower(answer) == "exit" {
		fmt.Println("Chat app successfully exited.")
		os.Exit(0)
	} else if strings.ToLower(answer) == "register" {

		var formRegister = forms.NewFormRegister(httpAddr)
		if auth := formRegister.Run(); auth != nil {
			authSession = *auth
		}

	} else if strings.ToLower(answer) == "login" {

		var formLogin = forms.NewFormLogin(httpAddr)
		if auth := formLogin.Run(); auth != nil {
			authSession = *auth
		}

	}

	fmt.Println(lipgloss.NewStyle().Margin(1, 0).Render("You're Loggedin !"))

	requestHeader := make(http.Header)
	requestHeader.Set("Authorization", "Bearer "+authSession.Token())

	spinner.New().
		Title("Please wait while connecting to chat server...").
		Action(func() {
			utils.Sleep(2)
			websocketConn, _, err = websocket.DefaultDialer.Dial(wsAddr, requestHeader)
		}).
		Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	defer websocketConn.Close()

	fmt.Println("Chat successfully connected.")
	fmt.Println(lipgloss.NewStyle().MarginTop(1).MarginBottom(1).Render("Welcome to Snap chat."))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	rl, err := readline.NewEx(&readline.Config{
		Prompt: "\033[32m»\033[0m ",
	})
	if err != nil {
		panic(err)
	}

	defer rl.Close()

	rl.SetPrompt(repl.Prompt(authSession, time.Now(), room))
	go func() {
		c := time.Tick(30 * time.Second)
		for now := range c {
			rl.SetPrompt(repl.Prompt(authSession, now, room))
			rl.Refresh()
		}
	}()

	go func() {
		for {
			_, message, err := websocketConn.ReadMessage()
			if err != nil {
				fmt.Printf("<<server: %s>>\n", err)
				close(interrupt)
				return
			}

			var msg models.ClientMessage

			if err := json.Unmarshal(message, &msg); err != nil {
				fmt.Printf(err.Error())
				close(interrupt)
				return
			}

			if msg.IsError {
				if slices.Contains(msg.Target, authSession.Username()) {
					io.WriteString(rl.Stdout(), repl.ClientErrorPrompt(authSession, room, string(msg.Message)))
				}

			} else {
				if msg.Sender.Username != authSession.Username() && slices.Contains(msg.Target, authSession.Username()) {
					io.WriteString(rl.Stdout(), repl.IncomingPromptWithSenderAndMessage(msg))
				}
			}
		}
	}()

	go func() {
		for {
			line, err := rl.Readline()
			if err == readline.ErrInterrupt {
				interrupt <- os.Interrupt
				return
			} else if err != nil {
				close(interrupt)
				return
			}
			if len(line) == 0 {
				continue
			}

			if strings.HasPrefix(line, "/group") {
				if err := commands.Group(rl, room, line, authSession, httpAddr); err != nil {
					io.WriteString(rl.Stdout(), repl.ClientErrorPrompt(authSession, room, err.Error()))
				}
				continue
			}

			if strings.HasPrefix(line, "/contact") {
				if err := commands.Contact(rl, room, line, authSession, httpAddr); err != nil {
					io.WriteString(rl.Stdout(), repl.ClientErrorPrompt(authSession, room, err.Error()))
				}
				continue
			}

			if strings.HasPrefix(line, "/dm") {
				room, err = commands.PrivateChat(websocketConn, rl, line, authSession, httpAddr, room)
				if err != nil {
					io.WriteString(rl.Stdout(), repl.ClientErrorPrompt(authSession, room, err.Error()))
					continue
				}
				if room != nil {
					rl.SetPrompt(repl.Prompt(authSession, time.Now(), room))
					rl.Refresh()
				}
				continue
			}

			if strings.HasPrefix(line, "/gm") {
				room, err = commands.GroupChat(websocketConn, rl, line, authSession, httpAddr, room)
				if err != nil {
					io.WriteString(rl.Stdout(), repl.ClientErrorPrompt(authSession, room, err.Error()))
					continue
				}
				if room != nil {
					rl.SetPrompt(repl.Prompt(authSession, time.Now(), room))
					rl.Refresh()
				}
				continue
			}

			if strings.HasPrefix(line, "/exit") {
				room = nil
				rl.SetPrompt(repl.Prompt(authSession, time.Now(), room))
				rl.Refresh()
				continue
			}

			if room != nil {

				room.Message = []byte(line)

				msgByte, err := json.Marshal(*room)
				if err != nil {
					fmt.Println("err:", err)
					return
				}
				err = websocketConn.WriteMessage(websocket.TextMessage, msgByte)
				if err != nil {
					fmt.Println("err:", err)
					return
				}
			}
		}
	}()

	select {
	case <-interrupt:
		websocketConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		fmt.Println("<<client: sent websocket close frame>>")
		websocketConn.Close()
		os.Exit(0)
	}
}
