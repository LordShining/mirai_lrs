package lrs

import (
	"math/rand"
	"time"
)

const (
	noGame         = "no_game"
	waitingForJoin = "join"
	roleSending    = "role_sending"
	duringGame     = "during_game"
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

func checkOrder(state, action string) bool {
	if state == noGame && action == newGameKeyWord {
		return true
	}
	if state == waitingForJoin && (action == joinGameKeyWord || action == quitGameKeyWord) {
		return true
	}
	if state == waitingForJoin && action == startGameKeyWord {
		return true
	}
	return false
}

func generateRole(pNum int, players *map[int64]*player, gods *[4]int64) {
	heroNum, villagerNum := 0, 0
	m, n := pNum%3, pNum/3
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
	for _, v := range *players {
		if villagerNum <= 0 {
			break
		}
		if rand.Intn(100)%2 == 0 {
			v.role = villager
			villagerNum--
		}
	}
	i := 0
	for k, v := range *players {
		if i >= heroNum {
			break
		}
		if v.role == "" {
			v.role = godList[i]
			gods[i] = k
			i++
		}
	}
	for _, v := range *players {
		if villagerNum <= 0 {
			break
		}
		if v.role == "" {
			v.role = villager
			villagerNum--
		}
	}
	for _, v := range *players {
		if heroNum <= 0 {
			break
		}
		if v.role == "" {
			v.role = wolf
			heroNum--
		}
	}
}
