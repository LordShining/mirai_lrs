package game

import (
	"sync"

	"github.com/LordShining/mirai_lrs/MiraiGo-Template/bot"

	"github.com/Mrs4s/MiraiGo/message"
)

type Server struct {
	state       string
	wolfNum     int
	godNum      int
	villagerNum int
	godList     [4]int64
	pNum        int
	idToUin     []int64
	vote        [13]int
	mu          sync.Mutex
}

func (s *Server) tryS(b *bot.Bot) {
	b.QQClient.SendPrivateMessage(1226286757, &message.SendingMessage{
		Elements: []message.IMessageElement{
			message.NewText("test"),
		},
	})
}
