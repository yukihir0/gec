package gec

type BlockList []*Block

func (self BlockList) Len() int {
	return len(self)
}

func (self BlockList) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self BlockList) Less(i, j int) bool {
	return self[i].Score < self[j].Score
}

func NewBlockList() (bl BlockList) {
	bl = BlockList{}
	bl = append(bl, NewBlock("", 0.0))
	return
}

func (self BlockList) AppendBlock() (bl BlockList) {
	self = append(self, NewBlock("", 0.0))
	return self
}
