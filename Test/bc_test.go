package Test

import (
	"common/bc"
	"fmt"
	"testing"
)

func TestNewWallet(t *testing.T) {
	wallet, err := bc.NewWallet()
	if err != nil {
		panic(err)
	}
	address := wallet.GetTestAddress()
	fmt.Println(string(address))
}
