package payoutmanagementsystem

import (
	"testing"
)

func TestMainPrintHello(t *testing.T) {
	str1 := "Hello World"
	print := Print(str1)
	if print != str1 {
		t.Errorf( "%s should be printed",str1 )
	}
}
