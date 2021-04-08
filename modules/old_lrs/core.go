package lrs1

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

const (
	noGame         = "no_game"
	waitingForJoin = "join"
	roleSending    = "role_sending"
	duringGame     = "during_game"
	night          = "night"
	daytime        = "daytime"
)

const (
	modelID          = "game.lrs"
	newGameKeyWord   = "狼人杀"
	joinGameKeyWord  = "+"
	quitGameKeyWord  = "-"
	startGameKeyWord = "开始"
)

const (
	wolf     = "wolf"
	villager = "villager"
	prophet  = "prophet"
	witch    = "witch"
	guard    = "guard"
	hunter   = "hunter"
)

var godList = [4]string{prophet, witch, guard, hunter}

func (s *server) checkOrder(state, action string) bool {
	if state == noGame && action == newGameKeyWord {
		return true
	}
	if state == waitingForJoin && (action == joinGameKeyWord || action == quitGameKeyWord) {
		return true
	}
	if state == waitingForJoin && action == startGameKeyWord {
		return true
	}
	id, err := strconv.Atoi(action)
	if err != nil {
		return false
	}
	if (state == night || state == daytime) && id >= 0 && id <= s.pNum {
		return true
	}
	return false
}

func (s *server) generateRole() {
	heroNum, villagerNum := 0, 0
	m, n := s.pNum%3, s.pNum/3
	switch m {
	case 0:
		heroNum = n
		villagerNum = n
	case 1:
		heroNum = n
		villagerNum = n + 1
	case 2:
		heroNum = n + 1
		villagerNum = n
	}
	rand.Seed(time.Now().UnixNano())
	for _, v := range s.players {
		if villagerNum <= 0 {
			break
		}
		if rand.Intn(100)%2 == 0 {
			v.role = villager
			villagerNum--
		}
	}
	i := 0
	for k, v := range s.players {
		if i >= heroNum {
			break
		}
		if v.role == "" {
			v.role = godList[i]
			s.godList[i] = k
			i++
		}
	}
	for _, v := range s.players {
		if villagerNum <= 0 {
			break
		}
		if v.role == "" {
			v.role = villager
			villagerNum--
		}
	}
	for _, v := range s.players {
		if heroNum <= 0 {
			break
		}
		if v.role == "" {
			v.role = wolf
			heroNum--
		}
	}
}

func (s *server) sendWolfCommand(qqClient *client.QQClient) {
	//大群
	qqClient.SendGroupMessage(qqGroup, &message.SendingMessage{
		Elements: []message.IMessageElement{
			message.NewText("天黑了，请留意小群/私信指令"),
		},
	})
	//狼群
	m := fmt.Sprintf("统一意见后，@我 发送目标号码(一人发送即可)发送0空刀")
	s.buildPlayerList(&m)
	qqClient.SendGroupMessage(wolfGroup, &message.SendingMessage{
		Elements: []message.IMessageElement{
			message.NewText(m),
		},
	})
}

func (s *server) sendGodCommand(deadID int, qqClient *client.QQClient) {
	for _, v := range s.godList {
		if v == 0 {
			break
		}
		switch s.players[v].role {
		case prophet:
			s.sendProphetCommand(deadID, v, qqClient)
		case witch:
			s.sendWitchCommand(deadID, v, qqClient)
		case guard:
			s.sendGuardCommand(deadID, v, qqClient)
		case hunter:
			s.sendHunterCommand(deadID, v, qqClient)
		default:
			return
		}
	}
}

func (s *server) sendProphetCommand(deadID int, uid int64, qqClient *client.QQClient) {}
func (s *server) sendWitchCommand(deadID int, uid int64, qqClient *client.QQClient)   {}
func (s *server) sendGuardCommand(deadID int, uid int64, qqClient *client.QQClient)   {}
func (s *server) sendHunterCommand(deadID int, uid int64, qqClient *client.QQClient) {
	if s.idToUin[deadID] != uid {
		return
	}
	m := fmt.Sprintf("你已出局，回复要带走的玩家编号，0放弃")
	s.buildPlayerList(&m)
	qqClient.SendPrivateMessage(uid, &message.SendingMessage{
		Elements: []message.IMessageElement{
			message.NewText(m),
		},
	})
}

func (s *server) buildPlayerList(m *string) {
	for i := 1; i <= s.pNum; i++ {
		*m += "\n" + strconv.Itoa(i) + ":" + s.players[s.idToUin[i]].name
		if s.players[s.idToUin[i]].dead {
			*m += "(out)"
		}
	}
}
