package t_test

import (
	"html/template"
	"log"
	"os"
	"testing"

	"github.com/shamaton/msgpackgen/internal/t"
)

func TestMain(m *testing.M) {

	t.RegisterGeneratedResolver()

	// 開始処理
	log.Print("setup")
	// パッケージ内のテストの実行
	code := m.Run()
	// 終了処理
	log.Print("tear-down")

	resetGeneratedCode()

	// テストの終了コードで exit
	os.Exit(code)
}

func resetGeneratedCode() {
	tpl := template.Must(template.New("").Parse(`package t

import "fmt"

func RegisterGeneratedResolver() {
	fmt.Println("this is dummy.")
}
`))

	file, err := os.OpenFile("./resolver.msgpackgen.go", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = tpl.Execute(file, tpl)
	if err != nil {
		log.Fatal(err)
	}
}

func TestA(t *testing.T) {

}
