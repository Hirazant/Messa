package handlers

import (
	"Diplom/app/model"
	telega "Diplom/app/telegram"
	"Diplom/app/vk"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type App struct {
	Dialogs  []model.Dialog
	TGClient *telega.Telega
	vk       *vk.Vk
}

func (a *App) Start() {
	vk := vk.NewVk("e8c2df0f8fd0570233bbab9372eafa4d13e6b4a92e1260e7616ef56a6c99139924713388af0ab02f60cf6")
	dialogs, err := vk.GetConversation()
	if err != nil {
		log.Fatalln(err)
	}

	tgClient := telega.Auth()

	dialogsT, err := tgClient.GetDialogs()
	if err != nil {
		log.Fatalln(err)
	}

	dialogs = append(dialogs, dialogsT...)

	sort.Slice(dialogs, func(i, j int) (less bool) {
		return dialogs[i].Date > dialogs[j].Date
	})
	dialogs = dialogs[0:15]

	a.Dialogs = dialogs
	a.TGClient = tgClient
	a.vk = vk
}

func (a *App) UpdateDialogs() {
	dialogs, err := a.vk.GetConversation()
	if err != nil {
		log.Fatalln(err)
	}

	dialogsT, err := a.TGClient.GetDialogs()
	if err != nil {
		log.Fatalln(err)
	}

	dialogs = append(dialogs, dialogsT...)

	sort.Slice(dialogs, func(i, j int) (less bool) {
		return dialogs[i].Date > dialogs[j].Date
	})
	dialogs = dialogs[0:15]

	a.Dialogs = dialogs
}

func (a *App) ShowIndexPage(c *gin.Context) {

	a.UpdateDialogs()

	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"index.html",
		// Pass the data that the page uses
		gin.H{
			"title":   "Home Page",
			"dialogs": a.Dialogs,
		},
	)

}

func (a *App) GetMessages(c *gin.Context) {
	// Проверим валидность ID
	if dialogId, err := strconv.Atoi(c.Param("dialog_id")); err == nil {
		// Проверим существование топика
		if messages, err := a.getMessagesByID(dialogId); err == nil {
			title, err := a.findDialogById(dialogId)
			if err != nil {
				log.Fatalln(err)
			}
			// Вызовем метод HTML из Контекста Gin для обработки шаблона
			c.HTML(
				// Зададим HTTP статус 200 (OK)
				http.StatusOK,
				// Используем шаблон index.html
				"messages.html",
				// Передадим данные в шаблон
				gin.H{
					"title":    title.Title,
					"data":     messages.MessagesList,
					"id":       messages.DialogId,
					"platform": title.Platform,
				},
			)

		} else {
			// Если топика нет, прервём с ошибкой
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		// При некорректном ID в URL, прервём с ошибкой
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func (a *App) findDialogById(id int) (*model.Dialog, error) {
	for _, d := range a.Dialogs {
		if d.Id == id {
			return &d, nil
		}
	}
	return nil, errors.New("Article not found")
}

func (a *App) getMessagesByID(id int) (*model.Messages, error) {
	dialog, err := a.findDialogById(id)
	if err != nil {
		return nil, err
	}
	if dialog.Platform == "vk" {
		return a.getMessagesFromVK(id)
	} else if dialog.Platform == "telega" {
		return a.getMessagesFromTelega(id)
	}
	return nil, nil
}

func (a *App) getMessagesFromVK(id int) (*model.Messages, error) {
	messList, err := a.vk.GetMessagesFromId(id)
	return &model.Messages{
		DialogId:     id,
		MessagesList: messList,
	}, err
}

func (a *App) getMessagesFromTelega(id int) (*model.Messages, error) {
	messList, err := a.TGClient.GetMessagesFromId(id)
	return &model.Messages{
		DialogId:     id,
		MessagesList: messList,
	}, err
}

type jsonPost struct {
	Message  string `json:"message"`
	IdString string `json:"idString"`
	Platform string `json:"platform"`
}

func (a *App) SendMessage(c *gin.Context) {
	var json jsonPost
	err := c.BindJSON(&json)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(json.Message, json.IdString, json.Platform)
	id, err := strconv.Atoi(json.IdString)
	if err != nil {
		log.Fatalln(err)
	}
	if json.Platform == "vk" {
		a.vk.SendVkMessage(id, json.Message)
	}
	if json.Platform == "telega" {
		a.TGClient.SendMessage(id, json.Message)
	}
}
