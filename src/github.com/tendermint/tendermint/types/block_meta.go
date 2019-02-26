package types

type BlockMeta struct {
	BlockID	BlockID	`json:"block_id"`
	Header	*Header	`json:"header"`
}

func NewBlockMeta(block *Block, blockParts *PartSet) *BlockMeta {
	return &BlockMeta{
		BlockID:	BlockID{block.Hash(), blockParts.Header()},
		Header:		block.Header,
	}
}
