package game

import (
	"sync"

	"github.com/LordShining/mirai_lrs/bot"

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

func (s *Server) TryS() {
	b := bot.Instance
	b.QQClient.SendPrivateMessage(1226286757, &message.SendingMessage{
		Elements: []message.IMessageElement{
			message.NewText("test"),
		},
	})
}
