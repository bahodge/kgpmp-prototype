package main

import "testing"

const ITERATIONS = 100

func BenchmarkRunCapn1(b *testing.B)       { benchmarkRunCapn(1, b) }
func BenchmarkRunCapn100(b *testing.B)     { benchmarkRunCapn(100, b) }
func BenchmarkRunCapn10000(b *testing.B)   { benchmarkRunCapn(10_000, b) }
func BenchmarkRunCapn100000(b *testing.B)  { benchmarkRunCapn(100_000, b) }
func BenchmarkRunCapn1000000(b *testing.B) { benchmarkRunCapn(1_000_000, b) }

func BenchmarkRunCBOR1(b *testing.B)       { benchmarkRunCBOR(1, b) }
func BenchmarkRunCBOR100(b *testing.B)     { benchmarkRunCBOR(100, b) }
func BenchmarkRunCBOR10000(b *testing.B)   { benchmarkRunCBOR(10_000, b) }
func BenchmarkRunCBOR100000(b *testing.B)  { benchmarkRunCBOR(100_000, b) }
func BenchmarkRunCBOR1000000(b *testing.B) { benchmarkRunCBOR(1_000_000, b) }

func BenchmarkRunJSON1(b *testing.B)       { benchmarkRunJSON(1, b) }
func BenchmarkRunJSON100(b *testing.B)     { benchmarkRunJSON(100, b) }
func BenchmarkRunJSON10000(b *testing.B)   { benchmarkRunJSON(10_000, b) }
func BenchmarkRunJSON100000(b *testing.B)  { benchmarkRunJSON(100_000, b) }
func BenchmarkRunJSON1000000(b *testing.B) { benchmarkRunJSON(1_000_000, b) }

func benchmarkRunCapn(iters int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		RunCapn(iters)
	}
}

func benchmarkRunCBOR(iters int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		RunCBOR(iters)
	}
}

func benchmarkRunJSON(iters int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		RunJSON(iters)
	}
}
