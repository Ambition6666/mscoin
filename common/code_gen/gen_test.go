package code_gen

import "testing"

func TestGenModel(t *testing.T) {
	GenProtoMessage("exchange_coins", "coins")
}
