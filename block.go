package main

type Block struct {
	Timestamp      int64
	Merkleroot     []byte
	PrevBlockcHash []byte
	Hash           []byte
}
