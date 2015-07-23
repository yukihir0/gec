package gec

type blockList []*block

func (bl blockList) Len() int {
	return len(bl)
}

func (bl blockList) Swap(i, j int) {
	bl[i], bl[j] = bl[j], bl[i]
}

func (bl blockList) Less(i, j int) bool {
	return bl[i].Score < bl[j].Score
}

func newBlockList() (bl blockList) {
	bl = blockList{}
	bl = append(bl, newBlock("", 0.0))
	return
}

func (bl blockList) AppendBlock() blockList {
	bl = append(bl, newBlock("", 0.0))
	return bl
}
