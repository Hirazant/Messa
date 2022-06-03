package telega

import (
	"Diplom/app/model"
	"log"
	"path/filepath"
	"strings"

	"github.com/zelenin/go-tdlib/client"
)

type Telega struct {
	Client *client.Client
}

func Auth() *Telega {
	authorizer := client.ClientAuthorizer()
	go client.CliInteractor(authorizer)

	// or bot authorizer
	// botToken := "000000000:gsVCGG5YbikxYHC7bP5vRvmBqJ7Xz6vG6td"
	// authorizer := client.BotAuthorizer(botToken)

	const (
		apiId   = 12138087
		apiHash = "7f71a056fc4f70176516af8122d24685"
	)

	authorizer.TdlibParameters <- &client.TdlibParameters{
		UseTestDc:              false,
		DatabaseDirectory:      filepath.Join(".tdlib", "database"),
		FilesDirectory:         filepath.Join(".tdlib", "files"),
		UseFileDatabase:        true,
		UseChatInfoDatabase:    true,
		UseMessageDatabase:     true,
		UseSecretChats:         false,
		ApiId:                  apiId,
		ApiHash:                apiHash,
		SystemLanguageCode:     "en",
		DeviceModel:            "Server",
		SystemVersion:          "1.0.0",
		ApplicationVersion:     "1.0.0",
		EnableStorageOptimizer: true,
		IgnoreFileNames:        false,
	}

	_, err := client.SetLogVerbosityLevel(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: 1,
	})
	if err != nil {
		log.Fatalf("SetLogVerbosityLevel error: %s", err)
	}

	tdlibClient, err := client.NewClient(authorizer)
	if err != nil {
		log.Fatalf("NewClient error: %s", err)
	}

	optionValue, err := tdlibClient.GetOption(&client.GetOptionRequest{
		Name: "version",
	})
	if err != nil {
		log.Fatalf("GetOption error: %s", err)
	}

	log.Printf("TDLib version: %s", optionValue.(*client.OptionValueString).Value)

	return &Telega{Client: tdlibClient}
}

func Short(s string, i int) string {
	strS := strings.Split(s, " ")
	var res []string
	if len(strS)-1 > i {
		res = strS[0:i]
	} else {
		res = strS
	}
	return strings.Join(res, " ")
}

func (tg *Telega) GetDialogs() ([]model.Dialog, error) {
	reqChats := client.GetChatsRequest{Limit: 15}
	chats, err := tg.Client.GetChats(&reqChats)
	if err != nil {
		return nil, err
	}
	res := make([]model.Dialog, 0, len(chats.ChatIds))
	for _, val := range chats.ChatIds {
		reqChat := client.GetChatRequest{ChatId: val}
		chat, err := tg.Client.GetChat(&reqChat)
		if err != nil {
			return nil, err
		}
		if chat.LastMessage.Content.MessageContentType() == client.TypeMessageText {
			res = append(res, model.Dialog{
				Id:            int(chat.Id),
				Title:         chat.Title,
				LastMessageId: chat.LastMessage.Id,
				LastMessage:   Short(chat.LastMessage.Content.(*client.MessageText).Text.Text, 10),
				Platform:      "telega",
				Date:          int(chat.LastMessage.Date),
			})
		}
	}

	return res, nil
}

func (tg *Telega) GetMessagesFromId(id int) ([]model.Message, error) {
	var res []model.Message
	reqChatH := client.GetChatHistoryRequest{
		ChatId:        int64(id),
		FromMessageId: 0,
		Offset:        0,
		Limit:         5,
		OnlyLocal:     false,
	}
	messages, err := tg.Client.GetChatHistory(&reqChatH)
	if err != nil {
		return nil, err
	}
	reqUser := client.GetUserRequest{
		UserId: messages.Messages[0].SenderId.(*client.MessageSenderUser).UserId,
	}
	author, err := tg.Client.GetUser(&reqUser)
	if err != nil {
		return nil, err
	}
	res = append(res, model.Message{
		MessageText:     messages.Messages[0].Content.(*client.MessageText).Text.Text,
		MessageId:       int(messages.Messages[0].Id),
		MessageIdAuthor: int(author.Id),
		MessageAuthor:   author.FirstName + " " + author.LastName,
		MessageDate:     int(messages.Messages[0].Date),
	})

	reqChatH = client.GetChatHistoryRequest{
		ChatId:        int64(id),
		FromMessageId: int64(res[0].MessageId),
		Offset:        0,
		Limit:         15,
		OnlyLocal:     false,
	}

	messag, err := tg.Client.GetChatHistory(&reqChatH)
	if err != nil {
		return nil, err
	}
	for _, val := range messag.Messages {

		reqUser = client.GetUserRequest{
			UserId: val.SenderId.(*client.MessageSenderUser).UserId,
		}
		author, err = tg.Client.GetUser(&reqUser)
		if err != nil {
			return nil, err
		}

		if val.Content.MessageContentType() == client.TypeMessageText {
			res = append(res, model.Message{
				MessageText:     val.Content.(*client.MessageText).Text.Text,
				MessageId:       int(val.Id),
				MessageAuthor:   author.FirstName + " " + author.LastName,
				MessageIdAuthor: int(author.Id),
				MessageDate:     int(val.Date),
			})
		}
	}

	return res, err
}

func (tg *Telega) SendMessage(id int, text string) {
	messageText := client.InputMessageText{Text: &client.FormattedText{Text: text}}

	req := client.SendMessageRequest{
		ChatId:              int64(id),
		InputMessageContent: &messageText,
	}
	_, err := tg.Client.SendMessage(&req)
	if err != nil {
		log.Fatalln(err)
	}
}
