package main_test

import (
	"fmt"
	"log"
	"os/exec"
	"testing"

	ttt "github.com/shamaton/msgpackgen/internal/tst"
)

func init() {
	fmt.Println("ininininit")

	cmd := "cd ./internal/t && go generate"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatalf("Failed to execute command: %s", cmd)
	}
	fmt.Println(string(out))

	//err := exec.Command("cd", "./internal/t", "&&", "go", "generate").Run()
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func TestADFG(t *testing.T) {

	// todo : 意図通りに出力があるか確認する
	// todo : 問題なければ、次のテスト用にコード生成する

	ttt.RegisterGeneratedResolver()
	t.Fatal("fatal!!!")
}
