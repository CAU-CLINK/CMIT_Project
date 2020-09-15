package main

type Blockchain struct {
	blocks []*Block
}

func (bc *Blockchain) AddNewBlock(Merkleroot string) {
	prevBlock := bc.blocks(len(bc.blocks) - 1)
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks := append(bc.blocks, newBlock)
}