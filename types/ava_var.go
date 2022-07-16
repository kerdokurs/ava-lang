package types

type AvaVar struct {
	Type    AvaType
	IsConst bool
	IsRef   bool
	Value   any
}
