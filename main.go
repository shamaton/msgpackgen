package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shamaton/msgpackgen/internal/gen"
)

var (
	out     = flag.String("output", "", "output directory")
	input   = flag.String("input", ".", "input directory")
	strict  = flag.Bool("strict", false, "strict mode")
	verbose = flag.Bool("v", false, "verbose diagnostics")
	pointer = flag.Int("pointer", 1, "pointer level to consider")
)

func main() {

	flag.Parse()

	_, err := os.Stat(*input)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *out == "" {
		*out = *input
	}

	if *pointer < 1 {
		*pointer = 1
	}

	g := gen.NewGenerator()
	g.Initialize(*input, *out, *pointer, *strict, *verbose)

	// todo : この呼び方やめる
	files := g.Dirwalk(*input)
	fmt.Println(files)

	// 最初にgenerate対象のパッケージをすべて取得
	// できればコードにエラーがない状態を知りたい

	// todo : 構造体の解析時にgenerate対象でないパッケージを含んだ構造体がある場合
	// 出力対象にしない

	// todo : 出力対象にしない構造体が見つからなくなるまで実行する

	// todo : エラーハンドリング

	g.GetPackages(files)
	g.CreateAnalyzedStructs()
	g.Generate()
}
