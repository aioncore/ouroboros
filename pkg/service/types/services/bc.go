package services

type Block interface {
	GetBlockID() string
	GetHash() string
}

type BlockHeader struct {
	id   string
	hash string
}

func (bh *BlockHeader) GetBlockID() string {
	return bh.id
}

func (bh *BlockHeader) GetHash() string {
	return bh.hash
}
