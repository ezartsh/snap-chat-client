package forms

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"
	"snap_client/models"
	"strings"
	"time"

	"github.com/charmbracelet/huh/spinner"
)

type FormContact struct {
	authSession models.Auth
	username    string
	formAction  string
}

type ContactRequestBody struct {
	Username string `json:"username"`
}

type ContactResponseBody struct {
	AccessToken string `json:"access_token"`
	Name        string `json:"name"`
	Username    string `json:"username"`
}

type ContactList struct {
	Username string `json:"username"`
}

func NewFormContact(authSession models.Auth, formAction string, username string) *FormContact {
	return &FormContact{
		authSession: authSession,
		username:    username,
		formAction:  formAction,
	}
}

func (r *FormContact) GetList() ([]ContactList, error) {
	requestUrl := r.formAction + "/contacts"
	contacts := []ContactList{}

	newContact := ContactRequestBody{
		Username: r.username,
	}

	msgByte, err := json.Marshal(newContact)
	if err != nil {
		return contacts, err
	}

	request, err := http.NewRequest("GET", requestUrl, bytes.NewBuffer(msgByte))
	request.Header.Set("Authorization", "Bearer "+r.authSession.Token())
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		return contacts, err
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
		return contacts, err
	}

	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	if !slices.Contains([]int{200, 201}, response.StatusCode) {
		return contacts, errors.New(strings.ReplaceAll(string(body), "\n", ""))
	}

	if err := json.Unmarshal(body, &contacts); err != nil {
		return contacts, err
	}

	return contacts, nil
}

func (r *FormContact) Add() error {
	requestUrl := r.formAction + "/contacts/create"

	newContact := ContactRequestBody{
		Username: r.username,
	}

	msgByte, err := json.Marshal(newContact)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(msgByte))
	request.Header.Set("Authorization", "Bearer "+r.authSession.Token())
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		return err
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
		return err
	}

	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	if !slices.Contains([]int{200, 201}, response.StatusCode) {
		return errors.New(strings.ReplaceAll(string(body), "\n", ""))
	}

	return nil
}

func (r *FormContact) Remove() error {
	requestUrl := r.formAction + "/contacts/remove"

	newContact := ContactRequestBody{
		Username: r.username,
	}

	msgByte, err := json.Marshal(newContact)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(msgByte))
	request.Header.Set("Authorization", "Bearer "+r.authSession.Token())
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		return err
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
		return err
	}

	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	if !slices.Contains([]int{200, 201}, response.StatusCode) {
		return errors.New(strings.ReplaceAll(string(body), "\n", ""))
	}

	return nil
}
