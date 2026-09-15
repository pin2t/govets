package main

import "testing"
import "golang.org/x/tools/go/analysis/analysistest"

func TestVardeclVet(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer, "clean", "short", "generated")
}
