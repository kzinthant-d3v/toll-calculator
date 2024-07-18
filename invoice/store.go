package main

import (
	"fmt"

	"github.com/kzinthant-d3v/toll-calculator/types"
)

type MemoryStore struct {
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) Insert(distance types.Distance) error {
	fmt.Println("Inserting distance to memory store", distance)
	return nil
}
