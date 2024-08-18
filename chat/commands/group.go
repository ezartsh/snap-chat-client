package commands

import (
	"errors"
	"fmt"
	"io"
	"snap_client/forms"
	"snap_client/models"
	"snap_client/repl"
	"strings"

	"gopkg.in/readline.v1"
)

func Group(rl *readline.Instance, room *repl.Room, line string, authSession models.Auth, httpAddr string) error {
	splitContent := strings.Split(line, " ")
	var act string
	var name string

	if len(splitContent) < 2 {
		return errors.New("Wrong command format. Valid format must be : /group [action] <group_name | optional>")
	}

	_, act = splitContent[0], splitContent[1]

	if splitContent[1] != "list" {

		if len(splitContent) < 3 {
			return errors.New("Wrong command format. Valid format must be : /group [action] <group_name | optional>")
		}

		name = strings.Join(splitContent[2:], " ")
		name = strings.Trim(name, " ")

		if name == "" {
			return errors.New("Group name must be provided.")
		}

	}

	if act == "list" {

		form := forms.NewFormGroup(authSession, httpAddr)
		groups, err := form.GetList()
		if err == nil {
			if len(groups) > 0 {
				for key, group := range groups {
					io.WriteString(
						rl.Stdout(),
						repl.ClientInfoPrompt(authSession, room, fmt.Sprintf("%d. %s - %s", key+1, group.Name, group.Key)),
					)
				}
			} else {
				io.WriteString(
					rl.Stdout(),
					repl.ClientInfoPrompt(authSession, room, "No group available on the list."),
				)
			}
		}

	} else if act == "create" {

		form := forms.NewFormGroup(authSession, httpAddr)
		if err := form.Create(name); err != nil {
			return err
		}
		io.WriteString(
			rl.Stdout(),
			repl.ClientInfoPrompt(authSession, room, "New group successfully created."),
		)

	} else if act == "join" {

		form := forms.NewFormGroup(authSession, httpAddr)
		if err := form.Join(name); err != nil {
			return err
		}
		io.WriteString(
			rl.Stdout(),
			repl.ClientInfoPrompt(authSession, room, "Successfully join to the group."),
		)

	} else if act == "leave" {

		form := forms.NewFormGroup(authSession, httpAddr)
		if err := form.Leave(name); err != nil {
			return err
		}

		io.WriteString(
			rl.Stdout(),
			repl.ClientInfoPrompt(authSession, room, "Successfully leave the group."),
		)

	}
	return nil
}
