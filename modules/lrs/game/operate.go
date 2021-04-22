package game

import (
	"fmt"
	"strconv"

	"github.com/LordShining/mirai_lrs/bot"
	"github.com/LordShining/mirai_lrs/modules/lrs/game/state"
	"github.com/Mrs4s/MiraiGo/message"
)

type keyOperate func(s *Server, sender Sender, keyList ...string)

func joinGroupOperate(s *Server, sender Sender, keyList ...string) {}

//测试操作
func testOperate(s *Server, sender Sender, keyList ...string) {
	bot.Instance.QQClient.SendGroupMessage(s.mainGroup, &message.SendingMessage{
		Elements: []message.IMessageElement{
			message.NewText(fmt.Sprintf("test\n sender: %v keyword(s): %v", sender.QQUin, keyList)),
		},
	})
}

//创建房间操作
func createRoomOperate(s *Server, sender Sender, keyList ...string) {
	if sender.QQGroup != s.mainGroup {
		return
	}
	var m string
	var err error
	if s.state != state.NoRoom {
		m = "@" + sender.Name + " 正在游戏中，请勿重复建房"
		goto send
	}
	if len(keyList) < 1 {
		m = "创建房间说明:\n" +
			"指令（狼人杀 板子名）以指定的板子创建房间"
		goto send
	}
	err = s.CreateRoom(keyList[1])
	if err != nil {
		m = err.Error()
		goto send
	}
	m = "@全体成员 来玩狼人杀呀~"
	for i, v := range s.GetPlayerList() {
		m += "\n" + strconv.Itoa(i+1) + ": " + v
	}
send:
	bot.Instance.QQClient.SendGroupMessage(s.mainGroup, &message.SendingMessage{
		Elements: []message.IMessageElement{
			message.NewText(m),
		},
	})
}

//新增板子操作
func addBoardOperate(s *Server, sender Sender, keyList ...string) {
	if sender.QQGroup != s.mainGroup {
		return
	}
	var m string
	if s.state != state.NoRoom {
		m = "@" + sender.Name + " 正在游戏中，请勿添加板子"
		goto send
	}
	if len(keyList) < 1 {
		m = "添加板子说明:\n" +
			"指令（加板子 现有角色）获得当前已有角色列表\n" +
			"指令（加板子 板子名 3 狼 3 村民 1 预言家 1 女巫 1 守卫）新增板子，如果板子名重复则覆盖"
		goto send
	}
	if keyList[0] == "现有角色" {
		rl := s.rm.GetRoleList()
		m = "现有角色（如需增加新角色请联系开发人员）"
		for _, v := range rl {
			m += "\n" + v
		}
		goto send
	}
send:
	bot.Instance.QQClient.SendGroupMessage(s.mainGroup, &message.SendingMessage{
		Elements: []message.IMessageElement{
			message.NewText(m),
		},
	})
}
