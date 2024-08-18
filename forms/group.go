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

type FormGroup struct {
	authSession models.Auth
	formAction  string
}

type GroupCreateRequestBody struct {
	Name string `json:"name"`
}

type GroupKeyRequestBody struct {
	Key string `json:"key"`
}

type GroupResponseBody struct {
	AccessToken string `json:"access_token"`
	Name        string `json:"name"`
}

type GroupList struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

func NewFormGroup(authSession models.Auth, formAction string) *FormGroup {
	return &FormGroup{
		authSession: authSession,
		formAction:  formAction,
	}
}

func (r *FormGroup) GetList() ([]GroupList, error) {
	requestUrl := r.formAction + "/groups"
	groups := []GroupList{}

	request, err := http.NewRequest("GET", requestUrl, nil)
	request.Header.Set("Authorization", "Bearer "+r.authSession.Token())
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		return groups, err
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
		return groups, err
	}

	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	if !slices.Contains([]int{200, 201}, response.StatusCode) {
		return groups, errors.New(strings.ReplaceAll(string(body), "\n", ""))
	}

	if err := json.Unmarshal(body, &groups); err != nil {
		return groups, err
	}

	return groups, nil
}

func (r *FormGroup) Create(name string) error {
	requestUrl := r.formAction + "/groups/create"

	newGroup := GroupCreateRequestBody{
		Name: name,
	}

	msgByte, err := json.Marshal(newGroup)
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

func (r *FormGroup) Join(key string) error {
	requestUrl := r.formAction + "/groups/join"

	newGroup := GroupKeyRequestBody{
		Key: key,
	}

	msgByte, err := json.Marshal(newGroup)
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

func (r *FormGroup) Leave(key string) error {
	requestUrl := r.formAction + "/groups/leave"

	newGroup := GroupKeyRequestBody{
		Key: key,
	}

	msgByte, err := json.Marshal(newGroup)
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
