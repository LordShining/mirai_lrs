package lrs

import (
	"fmt"
	"strings"
	"sync"

	"github.com/LordShining/mirai_lrs/bot"
	"github.com/LordShining/mirai_lrs/modules/lrs/game"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

const (
	modelID         = "game.lrs.adapter"
	mainGroup int64 = 883866459
	wolfGroup int64 = 765359710
	testQQ    int64 = 1226286757
)

type adapter struct {
	server    *game.Server
	mainGroup int64
	wolfGroup int64
}

var instance *adapter

func init() {
	instance = &adapter{
		server: game.NewServer(mainGroup, wolfGroup),
	}
	bot.RegisterModule(instance)
}

func (a *adapter) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       modelID,
		Instance: instance,
	}
}

func (a *adapter) Init() {
	fmt.Println("loading lrs adapter")
}

func (a *adapter) PostInit() {}

func (a *adapter) Serve(b *bot.Bot) {
	registerFunc(b)
}

func (a *adapter) Start(b *bot.Bot) {}

func (a *adapter) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func dealWithPrivateMessage(qqClient *client.QQClient, msg *message.PrivateMessage) {
	if msg.Sender.Uin != testQQ && !instance.server.IsPlayer(msg.Sender.Uin) {
		return
	}
	var keyword string
	for _, elem := range msg.Elements {
		switch e := elem.(type) {
		case *message.TextElement:
			keyword += " " + e.Content
		}
	}
	kList := strings.Fields(keyword)
	sender := game.Sender{
		QQUin: msg.Sender.Uin,
		Name:  msg.Sender.DisplayName(),
	}
	err := instance.server.InputMessage(sender, kList)
	if err != nil {
		qqClient.SendPrivateMessage(msg.Sender.Uin, &message.SendingMessage{
			Elements: []message.IMessageElement{
				message.NewText(err.Error()),
			},
		})
	}
}

func dealWithGroupMessage(qqClient *client.QQClient, msg *message.GroupMessage) {
	if msg.GroupCode != mainGroup && msg.GroupCode != wolfGroup {
		return
	}
	isAtMe := false
	var keyword string
	for _, elem := range msg.Elements {
		switch e := elem.(type) {
		case *message.AtElement:
			if e.Display == "@LordShining" {
				isAtMe = true
			}
			break
		case *message.TextElement:
			keyword += " " + e.Content
		}
	}
	if !isAtMe {
		return
	}
	kList := strings.Fields(keyword)
	sender := game.Sender{
		QQUin:   msg.Sender.Uin,
		Name:    msg.Sender.DisplayName(),
		QQGroup: msg.GroupCode,
	}
	err := instance.server.InputMessage(sender, kList)
	if err != nil {
		qqClient.SendGroupMessage(msg.GroupCode, &message.SendingMessage{
			Elements: []message.IMessageElement{
				message.NewText("@" + msg.Sender.DisplayName() +
					" " + err.Error()),
			},
		})
	}
}

func dealWithJoinGroup(qqClient *client.QQClient, req *client.UserJoinGroupRequest) {
	if req.GroupCode != wolfGroup {
		return
	}
	if !instance.server.IsPlayer(req.RequesterUin) {
		return
	}
	sender := game.Sender{
		QQUin: req.RequesterUin,
	}
	err := instance.server.InputMessage(sender, []string{"加群"})
	if err != nil {

	}
}

func registerFunc(b *bot.Bot) {
	b.OnPrivateMessage(dealWithPrivateMessage)
	b.OnGroupMessage(dealWithGroupMessage)
	b.OnUserWantJoinGroup(dealWithJoinGroup)
}
