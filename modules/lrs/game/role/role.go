package role

type werewolves struct {
}

func newWerewolves() interface{} {
	return &werewolves{}
}

func (w *werewolves) DoSomething() {}
