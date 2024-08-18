package forms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"snap_client/models"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

type FormLogin struct {
	name       string
	username   string
	password   string
	formAction string
}

type LoginResponseBody struct {
	AccessToken string `json:"access_token"`
	Name        string `json:"name"`
	Username    string `json:"username"`
}

func NewFormLogin(formAction string) *FormLogin {
	return &FormLogin{
		formAction: formAction + "/login",
	}
}

func (r *FormLogin) Run() *models.Auth {

	promptUsername := &survey.Input{
		Message: "Enter your username here :",
	}
	if err := survey.AskOne(promptUsername, &r.username, survey.WithValidator(survey.Required)); err != nil {
		fmt.Println("Chat app successfully exited.")
		os.Exit(0)
	}

	promptPassword := &survey.Password{
		Message: "Please type your password :",
	}

	if err := survey.AskOne(promptPassword, &r.password, survey.WithValidator(survey.Required)); err != nil {
		fmt.Println("Chat app successfully exited.")
		os.Exit(0)
	}

	loginUser := models.User{
		Username: r.username,
		Password: r.password,
	}

	msgByte, err := json.Marshal(loginUser)
	if err != nil {
		fmt.Println("err:", err)
		os.Exit(0)
	}

	request, err := http.NewRequest("POST", r.formAction, bytes.NewBuffer(msgByte))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		fmt.Println("err:", err)
		os.Exit(0)
	}

	var response *http.Response
	client := &http.Client{}

	spinner.New().
		Title("Please wait ...").
		Action(func() {
			time.Sleep(1 * time.Second)
			response, err = client.Do(request)
		}).
		Run()

	if err != nil {
		fmt.Println("err:", err)
		os.Exit(0)
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	if !slices.Contains([]int{200, 201}, response.StatusCode) {

		fmt.Printf(lipgloss.NewStyle().Foreground(lipgloss.Color("#f53911")).MarginTop(1).MarginBottom(1).Render("X " + string(body)))

		return r.Run()

	}

	var loggedinResponse LoginResponseBody

	if err := json.Unmarshal(body, &loggedinResponse); err != nil {
		fmt.Println("err:", err)
		os.Exit(0)
	}

	return models.NewAuth(loggedinResponse.Name, loggedinResponse.Username, loggedinResponse.AccessToken)
}
