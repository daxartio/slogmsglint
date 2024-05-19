package main

import (
	"github.com/daxartio/slogmsglint"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(slogmsglint.Analyzer)
}
