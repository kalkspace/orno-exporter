package models

type RegisterConfig struct {
	Num   uint16
	Label string
	Size  int
	Type  BinType
}

type BinType int

const (
	BinTypeUnknown BinType = iota
	BinTypeUint16
	BinTypeFloat
)
