package log

import (
	"fmt"

	plbytes "github.com/0xPellNetwork/pelldvs-libs/bytes"
)

type LazySprintf struct {
	format string
	args   []any
}

// NewLazySprintf defers fmt.Sprintf until the Stringer interface is invoked.
// This is particularly useful for avoiding calling Sprintf when debugging is not
// active.
func NewLazySprintf(format string, args ...any) *LazySprintf {
	return &LazySprintf{format, args}
}

func (l *LazySprintf) String() string {
	return fmt.Sprintf(l.format, l.args...)
}

type LazyBlockHash struct {
	block hashable
}

type hashable interface {
	Hash() plbytes.HexBytes
}

// NewLazyBlockHash defers block Hash until the Stringer interface is invoked.
// This is particularly useful for avoiding calling Sprintf when debugging is not
// active.
func NewLazyBlockHash(block hashable) *LazyBlockHash {
	return &LazyBlockHash{block}
}

func (l *LazyBlockHash) String() string {
	return l.block.Hash().String()
}
