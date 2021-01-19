package main

import (
	"testing"
)

func init() {
	//fmt.Println("ininininit")
	//
	//cmd := "cd ./internal/tst && go generate"
	//out, err := exec.Command("bash", "-c", cmd).Output()
	//if err != nil {
	//	log.Fatalf("Failed to execute command: %s", cmd)
	//}
	//fmt.Println(string(out))

	//err := exec.Command("cd", "./internal/t", "&&", "go", "generate").Run()
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func TestMain(m *testing.M) {
	m.Run()
}

func TestADFG(t *testing.T) {
	//flag.CommandLine.Set("target", strconv.Itoa(tt.i)) // -target=iと指定したかの様に設定できる
	main()

	// todo : 意図通りに出力があるか確認する
	// todo : 問題なければ、次のテスト用にコード生成する

	//err := generator.Run(
	//	".",
	//	".",
	//	"resolver.msgpackgen.go",
	//	2,
	//	false,
	//	false,
	//)
	//if err != nil {
	//	panic(err)
	//}
}
