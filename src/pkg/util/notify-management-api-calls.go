package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/sean-cunniffe/notify-common/src/pkg/component"
)

const(
	notificationUrlEnv = "notification_management_url"
	registerSubUrl = "/register" 
	messageSubUrl = "/message"
)

var Notifier *notifier

type notifier struct{
	Component *component.Component
}

type NotificationBody struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

func (n notifier) SendNotification(message string){
	url, err := getNotificationMngtUrl()
	if err != nil{
		log.Println(err.Error())
	}
	url += messageSubUrl

	if n.Component == nil {
		log.Println("cannot send notification, component not registered")
		return
	}

	body := NotificationBody{Message: message, From: n.Component.Name}
	data, _ := json.Marshal(body)
	// TODO create body with from message
	resp, err := http.Post(url, "application/json", bytes.NewReader(data))

	if err != nil{
		log.Println(err.Error())
	}else{
		log.Printf("%.v\n", resp)
	}
	

}

/*
	Register a component with notify management
*/
func SetupNotificationHandling(name, href,command, description string) *notifier{
	// if href or command are "" then dont set but still allow component to send notification
	comp := &component.Component{Name:name, Href:href, Command:command, Description:description}
	if href != "" && command != "" && name != "" {
		url, err := getNotificationMngtUrl()
		if err != nil{
			log.Println(err.Error())
		}
		url += registerSubUrl
		tempBody, _ := json.Marshal(*comp)
		resp, err := http.Post(url, "application/json", bytes.NewReader(tempBody))
		if err != nil {
			log.Println(err)
		}else{
			log.Printf("%.v\n", resp)
		}
	}else{
		log.Printf("cannot setup component notification with following component %.v\n", *comp)
	}
	Notifier = &notifier{Component: comp}
	return Notifier
}

func getNotificationMngtUrl() (string, error){
	url, exists := os.LookupEnv(notificationUrlEnv)

	if !exists {
		return "", errors.New(notificationUrlEnv+" does not exist, cannot register component")
		
	}

	return url, nil

}