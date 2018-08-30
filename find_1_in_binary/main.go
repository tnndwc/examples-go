package main

import (
	"flag"
	"fmt"
	"strconv"
)

func main() {
	var numArgs = flag.Int("num", 0, "number")

	flag.Parse()

	var num, numOf1 = *numArgs, 0

	if num < 0 {
		num *= -1
	}

	for num > 0 {
		numOf1++
		num = num & (num - 1)
	}

	fmt.Println("num: " + strconv.Itoa(*numArgs) + ", numOf1: " + strconv.Itoa(numOf1))
}
