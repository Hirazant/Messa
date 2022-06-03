package main

func initializeRoutes() {

	// определение роута главной страницы
	router.GET("/", app.ShowIndexPage)

	// Обработчик GET-запросов на /article/view/некоторый_dialog_id
	router.GET("/messages/view/:dialog_id", app.GetMessages)

	router.POST("/messages/send", app.SendMessage)
}
