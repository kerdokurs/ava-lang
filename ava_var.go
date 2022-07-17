package main

type AvaVar struct {
	Type    AvaType
	IsConst bool
	IsRef   bool
	Value   AvaVal
}
