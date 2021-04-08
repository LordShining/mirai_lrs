package lrs1

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"

	"MiraiGo-Template/bot"
)

const (
	qqGroup int64 = 883866459
	// qqGroup   int64 = 721023753
	wolfGroup int64 = 765359710
)

type player struct {
	id                int
	role              string
	dead              bool
	name              string
	canReceivePrivate bool
}

type server struct {
	state       string
	wolfNum     int
	godNum      int
	villagerNum int
	godList     [4]int64
	pNum        int
	players     map[int64]*player
	idToUin     []int64
	vote        [13]int
	mu          sync.Mutex
}

var roleMap = map[string]string{
	wolf:     "狼人",
	villager: "村民",
	prophet:  "预言家",
	witch:    "女巫",
	guard:    "守卫",
	hunter:   "猎人",
}

var instance *server

func init() {
	instance = &server{
		state:   noGame,
		idToUin: make([]int64, 15),
	}
	instance.players = make(map[int64]*player)
	bot.RegisterModule(instance)
}

func (s *server) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       modelID,
		Instance: instance,
	}
}

func (s *server) Init() {
	fmt.Println("lrs")
}

func (s *server) PostInit() {

}

func (s *server) Serve(b *bot.Bot) {
	registerFunc(b)
}

func (s *server) Start(b *bot.Bot) {

}

func (s *server) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func dealWithGroupMessage(qqClient *client.QQClient, msg *message.GroupMessage) {
	if msg.GroupCode != qqGroup && msg.GroupCode != wolfGroup {
		return
	}
	isAtMe := false
	var keyWord string
	for _, elem := range msg.Elements {
		switch e := elem.(type) {
		case *message.AtElement:
			if e.Display == "@LordShining" {
				isAtMe = true
			}
			break
		case *message.TextElement:
			keyWord = e.Content[1:]
		}
	}
	if !isAtMe {
		return
	}
	model, err := bot.GetModule(modelID)
	if err != nil {
		fmt.Println(err)
		return
	}
	client, ok := model.Instance.(*server)
	if !ok {
		fmt.Println("model type error")
	}

	//指令检查
	fmt.Println(keyWord)
	ok = client.checkOrder(client.state, keyWord)
	if !ok {
		qqClient.SendGroupMessage(msg.GroupCode, &message.SendingMessage{
			Elements: []message.IMessageElement{
				message.NewText("@" + msg.Sender.DisplayName() +
					fmt.Sprintf("order error: unknown order or wrong state, now state:%s", client.state)),
			},
		})
		return
	}

	switch msg.GroupCode {
	case qqGroup:
		client.dealWithQQGroupMessage(keyWord, qqClient, msg)
	case wolfGroup:
		client.dealWithWolfGroupMessage(keyWord, qqClient)
	default:
	}
}

func (s *server) dealWithQQGroupMessage(keyWord string, qqClient *client.QQClient, msg *message.GroupMessage) {
	//解析指令
	switch keyWord {
	case newGameKeyWord:
		//新游戏
		qqClient.SendGroupMessage(qqGroup, &message.SendingMessage{
			Elements: []message.IMessageElement{
				message.AtAll(),
				message.NewText("来玩狼人杀啦"),
			},
		})
		s.state = waitingForJoin
	case joinGameKeyWord:
		//加入
		if _, ok := s.players[msg.Sender.Uin]; !ok {
			s.mu.Lock()
			s.pNum++
			s.players[msg.Sender.Uin] = &player{name: msg.Sender.DisplayName()}
			m := fmt.Sprintf("狼人杀组队中")
			i := 1
			for _, v := range s.players {
				m += "\n" + strconv.Itoa(i) + ":" + v.name
				i++
			}
			qqClient.SendGroupMessage(qqGroup, &message.SendingMessage{
				Elements: []message.IMessageElement{
					message.NewText(m),
				},
			})
			s.mu.Unlock()

		}
	case quitGameKeyWord:
		//离开
		if _, ok := s.players[msg.Sender.Uin]; !ok {
			qqClient.SendGroupMessage(qqGroup, &message.SendingMessage{
				Elements: []message.IMessageElement{
					message.NewText("@" + msg.Sender.DisplayName() +
						" 你还没有加入~"),
				},
			})
		} else {
			s.mu.Lock()
			s.pNum--
			delete(s.players, msg.Sender.Uin)
			m := fmt.Sprintf("狼人杀组队中")
			i := 1
			for _, v := range s.players {
				m += "\n" + strconv.Itoa(i) + ":" + v.name
				i++
			}
			qqClient.SendGroupMessage(qqGroup, &message.SendingMessage{
				Elements: []message.IMessageElement{
					message.NewText(m),
				},
			})
			s.mu.Unlock()
		}
	case startGameKeyWord:
		//开局，分配角色
		s.mu.Lock()
		if s.pNum < 6 {
			qqClient.SendGroupMessage(qqGroup, &message.SendingMessage{
				Elements: []message.IMessageElement{
					message.NewText("玩家少于6人，无法开始"),
				},
			})
			s.mu.Unlock()
			return
		}
		s.state = roleSending
		s.mu.Unlock()
		qqClient.SendGroupMessage(qqGroup, &message.SendingMessage{
			Elements: []message.IMessageElement{
				message.NewText("车门已焊死，正在分发角色"),
			},
		})
		//分配角色
		s.generateRole()
		i := 0
		//通知
		for k, v := range s.players {
			i++
			s.idToUin[i] = k
			v.id = i
			// if
			qqClient.SendPrivateMessage(k, &message.SendingMessage{
				Elements: []message.IMessageElement{
					message.NewText(fmt.Sprintf("你的角色是: %s", roleMap[v.role])),
				},
			})
		}
		m := fmt.Sprintf("游戏即将开始")
		s.buildPlayerList(&m)
		qqClient.SendGroupMessage(qqGroup, &message.SendingMessage{
			Elements: []message.IMessageElement{
				message.NewText(m),
			},
		})
		qqClient.SendGroupMessage(qqGroup, &message.SendingMessage{
			Elements: []message.IMessageElement{
				message.NewText(fmt.Sprintf("狼人请加群: %v", wolfGroup)),
			},
		})
		s.state = duringGame
	default:
	}
}

func (s *server) dealWithWolfGroupMessage(keyWord string, qqClient *client.QQClient) {
	id, err := strconv.Atoi(keyWord)
	if err != nil {
		fmt.Println(err)
		return
	}
	if s.players[s.idToUin[id]].dead {
		qqClient.SendGroupMessage(wolfGroup, &message.SendingMessage{
			Elements: []message.IMessageElement{
				message.NewText("该玩家已出局，请重新选择"),
			},
		})
		return
	}
	//标记，神职通知
	s.players[s.idToUin[id]].dead = true
	s.sendGodCommand(id, qqClient)
}

func dealWithPrivateMessage(qqClient *client.QQClient, msg *message.PrivateMessage) {

	model, err := bot.GetModule(modelID)
	if err != nil {
		fmt.Println(err)
		return
	}
	client, ok := model.Instance.(*server)
	if !ok {
		fmt.Println("model type error")
	}

	if _, ok := client.players[msg.Sender.Uin]; !ok {
		return
	}
	if !client.players[msg.Sender.Uin].canReceivePrivate {
		return
	}
	var keyWord string
	for _, elem := range msg.Elements {
		switch e := elem.(type) {
		case *message.TextElement:
			keyWord = e.Content
		}
	}
	fmt.Println(keyWord)
	ok = client.checkOrder(client.state, keyWord)
	if !ok {
		qqClient.SendPrivateMessage(msg.Sender.Uin, &message.SendingMessage{
			Elements: []message.IMessageElement{
				message.NewText(fmt.Sprintf("order error: unknown order or wrong state, now state:%s", client.state)),
			},
		})
		return
	}

	switch client.state {
	case night:
		client.dealWithNightPrivate(keyWord, qqClient) //神职
	case daytime:
		client.dealWithDaytimePrivate(msg.Sender.Uin, qqClient) //放逐
	default:
		return
	}

	qqClient.SendGroupPoke(qqGroup, 1226286757) //戳一戳
	qqClient.GroupMuteAll(wolfGroup, false)
}

func (s *server) dealWithNightPrivate(keyWord string, qqClient *client.QQClient) {

}
func (s *server) dealWithDaytimePrivate(uid int64, qqClient *client.QQClient) {}

func dealWithJoinGroup(qqClient *client.QQClient, req *client.UserJoinGroupRequest) {

	if req.GroupCode != wolfGroup {
		return
	}
	model, err := bot.GetModule(modelID)
	if err != nil {
		fmt.Println(err)
		return
	}
	client, ok := model.Instance.(*server)
	if !ok {
		fmt.Println("model type error")
	}
	if _, ok := client.players[req.RequesterUin]; !ok || client.players[req.RequesterUin].role != wolf {
		qqClient.SolveGroupJoinRequest(req, false, false, "你不是狼人~")
		return
	}
	qqClient.SolveGroupJoinRequest(req, true, false, "")
	//记录人数
	wolfGroupInfo, err := qqClient.GetGroupInfo(wolfGroup)
	if err != nil {
		fmt.Println(err)
		return
	}
	if wolfGroupInfo.MemberCount-1 == uint16(client.wolfNum) {
		//狼齐，进入夜晚
		client.state = night
		client.sendWolfCommand(qqClient)
	}
}

func registerFunc(b *bot.Bot) {
	b.OnGroupMessage(dealWithGroupMessage)
	b.OnPrivateMessage(dealWithPrivateMessage)
	b.OnUserWantJoinGroup(dealWithJoinGroup)
}
