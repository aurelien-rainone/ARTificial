package main

import (
	"image/color"
	"math/rand"
	"testing"
)

func checkB(b *testing.B, err error) {
	if err != nil {
		b.Helper()
		b.Fatal("error:", err)
	}
}

// create an imageDNA for testing purposes, with 50 randomly generated polygons
// of the same color.
func createTestCandidate(r, g, b uint8) *imageDNA {
	const numPolys = 50
	rng := rand.New(rand.NewSource(0))
	var img = &imageDNA{
		w:     128,
		h:     128,
		polys: make([]poly, numPolys),
	}
	// add N `numPolys` random polygons
	for i := 0; i < numPolys; i++ {
		img.polys[i] = randomPoly(img, 3, 4, rng)
		img.polys[i].col = color.RGBA{r, g, b, 255}
	}
	return img
}

func BenchmarkRenderImageDNA(b *testing.B) {
	c := createTestCandidate(255, 0, 0)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		c.render()
	}
}

func BenchmarkFitnessEvaluator(b *testing.B) {
	want := 0.0

	refImageFn := "testdata/red.png"
	ref, err := loadPNGAsRGBA(refImageFn)
	if err != nil {
		b.Fatal(err)
	}

	cand := monochromeImage(ref.Bounds().Dx(), ref.Bounds().Dy(), color.RGBA{255, 0, 0, 255})
	evaluator := fitnessEvaluator{img: ref}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StartTimer()
		got := evaluator.Fitness(cand, nil)
		b.StopTimer()

		if got != want {
			b.Fatalf("wrong fitness, want %v, got %v", want, got)
		}
	}
}

func BenchmarkCairoFitnessEvaluator(b *testing.B) {
	want := 0.0

	refImageFn := "testdata/red.png"
	cand := monochromeImage(128, 128, color.RGBA{255, 0, 0, 255})

	evaluator := newCairoEvaluator(refImageFn)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StartTimer()
		got := evaluator.Fitness(cand, nil)
		b.StopTimer()

		if got != want {
			b.Fatalf("wrong fitness, want %v, got %v", want, got)
		}
	}
}
