package utils

import (
	"sync"
)

type Lazy[T any] struct {
	New   func() T
	once  sync.Once
	value T
}

func (this *Lazy[T]) Value() T {
	if this.New != nil {
		this.once.Do(func() {
			this.value = this.New()
			this.New = nil
		})
	}
	return this.value
}

func NewLazy[T any](newfunc func() T) *Lazy[T] {
	return &Lazy[T]{New: newfunc}
}
