package game

import (
	"github.com/LordShining/mirai_lrs/modules/lrs/game/role"
)

type empty struct{}

type seat struct {
	player  int64
	role    role.Role
	isAlive bool
	buff    map[string]empty
}
