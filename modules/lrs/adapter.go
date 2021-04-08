package lrs

import (
	"fmt"
	"sync"

	"github.com/LordShining/mirai_lrs/MiraiGo-Template/bot"
	"github.com/LordShining/mirai_lrs/MiraiGo-Template/modules/lrs/game"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

const (
	modelID = "game.lrs.adapter"
)

type adapter struct {
	server *game.Server
}

var instance *adapter

func init() {
	instance = &adapter{
		server: new(game.Server), //要换构造函数
	}
	// instance.players = make(map[int64]*player)
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

}

func registerFunc(b *bot.Bot) {

}
