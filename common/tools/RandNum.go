package tools

import (
	"fmt"
	"math/rand"
	"time"
)

func Rand4Num() string {
	res := rand.Intn(9999)
	if res < 1000 {
		res += 1000
	}

	return fmt.Sprintf("%d", res)
}
func Rand6Num() string {
	res := rand.Intn(999999)
	if res < 100000 {
		res += 100000
	}

	return fmt.Sprintf("%d", res)
}

func Unique(str string) string {
	milli := time.Now().UnixMilli()
	return fmt.Sprintf("%s%s%d", str, Rand6Num(), milli)
}
