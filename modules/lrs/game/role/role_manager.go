package role

type Role interface {
	DoSomething()
}

type newRoleFunc func() interface{}

type RoleManager struct {
	roleList map[string]newRoleFunc
}

func NewRoleManager() *RoleManager {
	return &RoleManager{
		roleList: map[string]newRoleFunc{
			Werewolves: newWerewolves,
		},
	}
}

func (rm *RoleManager) GetRoleList() []string {
	var res []string
	for k := range rm.roleList {
		res = append(res, k)
	}
	return res
}
