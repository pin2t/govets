package main

import "testing"
import "golang.org/x/tools/go/analysis/analysistest"

func TestGroupdeclVet(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer, "clean", "grouped", "generated")
}
