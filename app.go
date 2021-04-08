package main

import (
	"os"
	"os/signal"

	"github.com/LordShining/mirai_lrs/bot"
	"github.com/LordShining/mirai_lrs/config"
	"github.com/LordShining/mirai_lrs/utils"

	_ "github.com/LordShining/mirai_lrs/modules/logging"
	_ "github.com/LordShining/mirai_lrs/modules/lrs"
)

func init() {
	utils.WriteLogToFS()
	config.Init()
}

func main() {
	// 快速初始化
	bot.Init()

	// 初始化 Modules
	bot.StartService()

	// 使用协议
	// 不同协议可能会有部分功能无法使用
	// 在登陆前切换协议
	bot.UseProtocol(bot.IPad)

	// 登录
	bot.Login()

	// 刷新好友列表，群列表
	bot.RefreshList()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	bot.Stop()
}
