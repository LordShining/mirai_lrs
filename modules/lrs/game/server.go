package game

import (
	"fmt"
	"sync"

	"github.com/LordShining/mirai_lrs/modules/lrs/game/board"
	"github.com/LordShining/mirai_lrs/modules/lrs/game/role"
	"github.com/LordShining/mirai_lrs/modules/lrs/game/state"
)

type Sender struct {
	QQUin   int64
	Name    string
	QQGroup int64
}

type player struct {
	name   string
	seatID int
}

type Server struct {
	state string
	mu    sync.Mutex

	seats   map[int]*seat
	players map[int64]player
	keyMap  map[string]keyOperate

	wolfNum     int
	godNum      int
	villagerNum int

	mainGroup int64
	wolfGroup int64

	rm *role.RoleManager
	bm *board.BoardManager
}

func NewServer(mainGroup, wolfGroup int64) *Server {
	s := &Server{
		mainGroup: mainGroup,
		wolfGroup: wolfGroup,
		rm:        role.NewRoleManager(),
		bm:        board.NewBoardManager(),
		state:     state.NoRoom,
		keyMap: map[string]keyOperate{
			"测试":  testOperate,
			"加群":  joinGroupOperate,
			"狼人杀": createRoomOperate,
			"加板子": addBoardOperate,
		},
	}
	return s
}

func (s *Server) IsPlayer(uin int64) bool {
	_, ok := s.players[uin]
	return ok
}

func (s *Server) InputMessage(sender Sender, keyList []string) (err error) {
	if len(keyList) < 1 {
		return
	}
	if operate, ok := s.keyMap[keyList[0]]; ok {
		if len(keyList) == 1 {
			operate(s, sender)
			return
		}
		operate(s, sender, keyList[1:]...)
		return
	}
	return fmt.Errorf("invalid keyword(s): %v", keyList)
}

func (s *Server) CreateRoom(boardName string) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	//TODO
	//加载板子

	//添加房主

	s.state = state.WaitingForJoin
	return
}

func (s *Server) GetPlayerList() []string {
	var res []string
	switch s.state {
	case state.WaitingForJoin:
		for _, v := range s.players {
			res = append(res, v.name)
		}
	}
	return res
}
