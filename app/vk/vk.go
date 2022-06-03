package vk

import (
	"Diplom/app/model"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

type Vk struct {
	vkApi *api.VK
	token string
}

func NewVk(token string) *Vk {
	return &Vk{
		token: token,
		vkApi: api.NewVK(token),
	}
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

func (vk *Vk) GetConversation() ([]model.Dialog, error) {
	p := params.MessagesGetConversationsBuilder{}

	conversations, err := vk.vkApi.MessagesGetConversations(p.Params)
	if err != nil {
		return nil, err
	}
	dialogs := make([]model.Dialog, 0, len(conversations.Items))
	user := make([]string, 0)
	for _, conver := range conversations.Items {
		if conver.LastMessage.Text != "" {
			if conver.Conversation.Peer.Type == "chat" {

				dialog := model.Dialog{
					Id:          conver.Conversation.Peer.ID,
					Type:        conver.Conversation.Peer.Type,
					LastMessage: Short(conver.LastMessage.Text, 10),
					Title:       conver.Conversation.ChatSettings.Title,
					Date:        conver.LastMessage.Date,
					Platform:    "vk",
				}
				dialogs = append(dialogs, dialog)
			}

			if conver.Conversation.Peer.Type == "user" {
				dialog := model.Dialog{
					Id:          conver.Conversation.Peer.ID,
					Type:        conver.Conversation.Peer.Type,
					LastMessage: Short(conver.LastMessage.Text, 10),
					Date:        conver.LastMessage.Date,
					Platform:    "vk",
				}
				dialogs = append(dialogs, dialog)
				user = append(user, strconv.Itoa(conver.Conversation.Peer.ID))
			}
		}
	}

	if len(user) > 0 {
		usersName, err := vk.getUsersProfile(user)
		if err != nil {
			return nil, err
		}
		for i, dialog := range dialogs {
			if dialog.Type == "user" {
				dialogs[i].Title = usersName[strconv.Itoa(dialog.Id)]
			}
		}
	}

	return dialogs, nil
}

func (vk *Vk) getUsersProfile(Ids []string) (map[string]string, error) {
	paramsUser := params.NewUsersGetBuilder()
	paramsUser.UserIDs(Ids)
	paramsUser.Fields([]string{"photo_50"})
	get, err := vk.vkApi.UsersGet(paramsUser.Params)
	if err != nil {
		return nil, err
	}

	resMap := make(map[string]string)
	for _, user := range get {
		resMap[strconv.Itoa(user.ID)] = user.FirstName + " " + user.LastName
	}

	return resMap, nil
}

func (vk *Vk) GetMessagesFromId(id int) ([]model.Message, error) {
	p := params.NewMessagesGetHistoryBuilder()
	p.UserID(id)
	messageID, err := vk.vkApi.MessagesGetHistory(p.Params)
	if err != nil {
		return nil, err
	}
	res := make([]model.Message, 0, len(messageID.Items))
	authorsId := make(map[string]interface{})
	for _, message := range messageID.Items {
		if message.Text != "" {
			authorsId[strconv.Itoa(message.FromID)] = nil
			res = append(res, model.Message{
				MessageText:     message.Text,
				MessageId:       message.ID,
				MessageIdAuthor: message.FromID,
				MessageDate:     message.Date,
			})
		}
	}
	authorIsStr := make([]string, 0, len(authorsId))
	for author := range authorsId {
		authorIsStr = append(authorIsStr, author)
	}
	authors, err := vk.getUsersProfile(authorIsStr)

	for i, mess := range res {
		res[i].MessageAuthor = authors[strconv.Itoa(mess.MessageIdAuthor)]
	}
	return res, err
}

func (vk *Vk) SendVkMessage(id int, text string) {
	p := params.NewMessagesSendBuilder()
	p.PeerID(id)
	p.Message(text)
	randId := rand.Intn(1890000-160000) + 160000
	p.RandomID(randId)
	_, err := vk.vkApi.MessagesSend(p.Params)
	if err != nil {
		log.Fatalln(err)
	}
}
