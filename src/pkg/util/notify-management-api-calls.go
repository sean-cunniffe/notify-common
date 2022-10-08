package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/sean-cunniffe/notify-common/src/pkg/component"
)

const (
	registerSubUrl = "/register"
	messageSubUrl  = "/message"
)

var Notifier *notifier

type notifier struct {
	Component *component.Component
	mngtUrl   string
}

type NotificationBody struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

func (n notifier) SendNotification(message string) error {
	url := n.mngtUrl
	url += messageSubUrl

	if n.Component == nil {
		return errors.New("cannot send notification, component not registered")
	}

	body := NotificationBody{Message: message, From: n.Component.Name}
	data, _ := json.Marshal(body)
	// TODO create body with from message
	resp, err := http.Post(url, "application/json", bytes.NewReader(data))

	if err != nil {
		return err
	} else {
		log.Printf("%.v\n", resp)
	}
	return nil
}

/*
	Register a component with notify management
*/
func SetupNotificationHandling(name, href, command, description, notificationMngtUrl string) (*notifier, error) {
	// if href or command are "" then dont set but still allow component to send notification
	comp := &component.Component{Name: name, Href: href, Command: command, Description: description}
	if href == "" || command == "" || name == "" {
		errMsg := fmt.Sprintf("cannot setup component notification with following component %.v\n", *comp)
		return nil, errors.New(errMsg)
	}
	notificationMngtUrl += registerSubUrl
	tempBody, _ := json.Marshal(*comp)
	_, err := http.Post(notificationMngtUrl, "application/json", bytes.NewReader(tempBody))
	if err != nil {
		return nil, err
	}
	Notifier = &notifier{Component: comp, mngtUrl: notificationMngtUrl}
	return Notifier, nil
}
