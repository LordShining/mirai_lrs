package lrs

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"

	"github.com/Logiase/MiraiGo-Template/bot"
	// "github.com/Logiase/MiraiGo-Template/utils"
)

const qqGroup int64 = 721023753

type player struct {
	id   int
	role string
	dead bool
	name string
}

type server struct {
	state       string
	wolfNum     int
	godNum      int
	villagerNum int
	godList     [4]int64
	wolfGroup   int64
	pNum        int
	players     map[int64]*player
	idToUin     []int64
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

func dealWithGroupMessage(msg *message.GroupMessage) {
	if msg.GroupCode != qqGroup {
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
	model, err := bot.GetModule(modelID)
	if err != nil {
		fmt.Println(err)
		return
	}
	client, ok := model.Instance.(*server)
	if !ok {
		fmt.Println("model type error")
	}
	if isAtMe {
		//解析指令
		fmt.Println(keyWord)
		ok := checkOrder(client.state, keyWord)
		if !ok {
			bot.Instance.QQClient.SendGroupMessage(qqGroup, &message.SendingMessage{
				Elements: []message.IMessageElement{
					message.NewText("@" + msg.Sender.DisplayName() +
						fmt.Sprintf("order error: unknown order or wrong state, now state:%s", client.state)),
				},
			})
			return
		}
		switch keyWord {
		case newGameKeyWord:
			//新游戏
			bot.Instance.QQClient.SendGroupMessage(qqGroup, &message.SendingMessage{
				Elements: []message.IMessageElement{
					message.AtAll(),
					message.NewText("来玩狼人杀啦"),
				},
			})
			client.state = waitingForJoin
		case joinGameKeyWord:
			//加入
			if _, ok := client.players[msg.Sender.Uin]; !ok {
				client.mu.Lock()
				client.pNum++
				client.players[msg.Sender.Uin] = &player{name: msg.Sender.DisplayName()}
				m := fmt.Sprintf("狼人杀组队中")
				i := 1
				for _, v := range client.players {
					m += "\n" + strconv.Itoa(i) + ":" + v.name
					i++
				}
				bot.Instance.QQClient.SendGroupMessage(qqGroup, &message.SendingMessage{
					Elements: []message.IMessageElement{
						message.NewText(m),
					},
				})
				client.mu.Unlock()

			}
		case quitGameKeyWord:
			//离开
			if _, ok := client.players[msg.Sender.Uin]; !ok {
				bot.Instance.QQClient.SendGroupMessage(qqGroup, &message.SendingMessage{
					Elements: []message.IMessageElement{
						message.NewText("@" + msg.Sender.DisplayName() +
							" 你还没有加入哦"),
					},
				})
			} else {
				client.mu.Lock()
				client.pNum--
				delete(client.players, msg.Sender.Uin)
				m := fmt.Sprintf("狼人杀组队中")
				i := 1
				for _, v := range client.players {
					m += "\n" + strconv.Itoa(i) + ":" + v.name
					i++
				}
				bot.Instance.QQClient.SendGroupMessage(qqGroup, &message.SendingMessage{
					Elements: []message.IMessageElement{
						message.NewText(m),
					},
				})
				client.mu.Unlock()
			}
		case startGameKeyWord:
			//开局，分配角色
			client.mu.Lock()
			if client.pNum < 6 {
				bot.Instance.QQClient.SendGroupMessage(qqGroup, &message.SendingMessage{
					Elements: []message.IMessageElement{
						message.NewText("玩家少于6人，无法开始"),
					},
				})
				client.mu.Unlock()
				return
			}
			client.state = roleSending
			client.mu.Unlock()
			bot.Instance.QQClient.SendGroupMessage(qqGroup, &message.SendingMessage{
				Elements: []message.IMessageElement{
					message.NewText("车门已焊死，正在分发角色"),
				},
			})
			//分配角色
			generateRole(client.pNum, &client.players, &client.godList)
			i := 0
			//通知
			for k, v := range client.players {
				i++
				client.idToUin[i] = k
				v.id = i
				// if
				bot.Instance.QQClient.SendPrivateMessage(k, &message.SendingMessage{
					Elements: []message.IMessageElement{
						message.NewText(fmt.Sprintf("你的角色是: %s", roleMap[v.role])),
					},
				})
			}
			//狼人群
			// bot.Instance.QQClient.
			//开始
			m := fmt.Sprintf("游戏即将开始")
			for i = 1; i <= client.pNum; i++ {
				m += "\n" + strconv.Itoa(i) + ":" + client.players[client.idToUin[i]].name
			}
			bot.Instance.QQClient.SendGroupMessage(qqGroup, &message.SendingMessage{
				Elements: []message.IMessageElement{
					message.NewText(m),
				},
			})
			bot.Instance.QQClient.SendGroupMessage(qqGroup, &message.SendingMessage{
				Elements: []message.IMessageElement{
					message.NewText(fmt.Sprintf("debug message:\n server:\n %v", client)),
				},
			})
			client.state = duringGame
		}
	}
}

func registerFunc(b *bot.Bot) {
	b.OnGroupMessage(func(qqClient *client.QQClient, groupMessage *message.GroupMessage) {
		dealWithGroupMessage(groupMessage)
	})
}
