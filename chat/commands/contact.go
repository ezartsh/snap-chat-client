package commands

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"snap_client/forms"
	"snap_client/models"
	"snap_client/repl"
	"strings"

	"gopkg.in/readline.v1"
)

func Contact(rl *readline.Instance, room *repl.Room, line string, authSession models.Auth, httpAddr string) error {
	splitContent := strings.Split(line, " ")
	var act string
	var username string

	if len(splitContent) < 2 {
		return errors.New("Wrong command format. Valid format must be : /contact [action] <username | optional>")
	}

	_, act = splitContent[0], splitContent[1]

	if !slices.Contains([]string{"list", "add", "remove"}, act) {
		return errors.New("Wrong command format for contact action. Valid format for action must be either : list, add or remove")
	}

	if act != "list" {

		if len(splitContent) != 3 {
			return errors.New("Wrong command format. Valid format must be : /contact [action] <username | optional>")
		}

		username = splitContent[2]

	}

	if act == "list" {

		form := forms.NewFormContact(authSession, httpAddr, username)
		contacts, err := form.GetList()
		if err == nil {
			if len(contacts) > 0 {
				for key, contact := range contacts {
					io.WriteString(
						rl.Stdout(), repl.ClientInfoPrompt(authSession, room, fmt.Sprintf("%d. %s", key+1, contact.Username)))
				}
			} else {
				io.WriteString(
					rl.Stdout(), repl.ClientInfoPrompt(authSession, room, "No contact available on the list."))
			}
		}

	} else if act == "add" {

		form := forms.NewFormContact(authSession, httpAddr, username)
		if err := form.Add(); err != nil {
			return err
		}

		io.WriteString(
			rl.Stdout(), repl.ClientInfoPrompt(authSession, room, "New Contact successfully added."))

	} else if act == "remove" {

		form := forms.NewFormContact(authSession, httpAddr, username)
		if err := form.Remove(); err != nil {
			return err
		}
		io.WriteString(
			rl.Stdout(), repl.ClientInfoPrompt(authSession, room, "Contact successfully removed."))

	}
	return nil
}
