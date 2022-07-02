package gommon

import (
	"fmt"
	"os"
	"testing"
)

func TestEval(t *testing.T) {
	fmt.Println("Start TestEval")
	code := "(print 1 (quote 2 3) (+ 4 5))"
	list, err := ReadString(code)
	if err != nil {
		t.Fatalf("Error: %s: %s", code, err.Error())
	}
	fmt.Println("Parsing:", code)
	_, err = list.Eval()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	fmt.Println("End TestEval")
}
