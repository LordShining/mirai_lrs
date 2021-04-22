package board

type BoardManager struct {
}

func NewBoardManager() *BoardManager {
	return &BoardManager{}
}

func (bm *BoardManager) GetBoardList() {}

//添加板子
func (bm *BoardManager) AddBoard() {}

//检查板子是否存在
func (bm *BoardManager) HasBoard(boardName string) bool {
	return true
}

//设置板子
func (bm *BoardManager) SetBoard() {}

//加载板子
func (bm *BoardManager) LoadBoard() {}

//持久化板子
func (bm *BoardManager) SaveBoard() {}
