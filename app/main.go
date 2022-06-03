package main

import (
	"Diplom/app/handlers"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine
var app handlers.App

func main() {
	router = gin.Default()
	router.LoadHTMLGlob("./templates/*")

	initializeRoutes()

	app.Start()

	router.Run()

	//vk := vk.NewVk("97a1b539b3ff778dff3f9f04f0b01999732c2e17545fb124bca28d675f2b7cb9a86c70371f786f3869126")
	//dialogs, err := vk.GetConversation()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//for _, dialog := range dialogs {
	//	fmt.Printf("Peple - %s, LastMessage - %s, Platform - %s, date - %d", dialog.Title, dialog.LastMessage, dialog.Platform, dialog.Date)
	//	fmt.Println("----------------------------------------------------------------")
	//}
	//
	//tgClient := telega.Auth()
	//
	//dialogsT, err := tgClient.GetDialogs()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Printf("\n\n\ntelega:\n")
	//for _, chat := range dialogsT {
	//	fmt.Printf("Peple - %s, LastMessage - %s, Platform - %s, date - %d", chat.Title, chat.LastMessage, chat.Platform, chat.Date)
	//	fmt.Println("----------------------------------------------------------------")
	//}
}
