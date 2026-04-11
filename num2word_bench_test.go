package propisyu

import (
	"math"
	"testing"
)

// BenchmarkIntToWordsMaxInt measures the cost of converting the largest
// supported integer. Before the O(n)-prepend → linear-append refactor,
// convertPositiveUint64ToWords prepended to its parts slice on every
// iteration, which is O(triads²). For math.MaxInt (~7 triads) this is
// O(49) small-slice allocations per call. After the refactor it is one
// preallocated slice and an in-place reverse: O(triads).
//
// Use this benchmark as a baseline — run `go test -bench BenchmarkIntToWords
// -benchmem ./...` to get allocs/op and ns/op numbers; regressions should
// show up as higher allocs/op.
func BenchmarkIntToWordsMaxInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = IntToWords(math.MaxInt)
	}
}

// BenchmarkIntToWordsSmall pins the hot path for small numbers (3 digits,
// 1 triad) which is what typical invoice / money-formatting workloads hit.
func BenchmarkIntToWordsSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = IntToWords(321)
	}
}

// BenchmarkIntToWordsCompound exercises the full multi-triad path with a
// realistic-large number used in the README (trillions range).
func BenchmarkIntToWordsCompound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = IntToWords(6_453_345_242_432)
	}
}
