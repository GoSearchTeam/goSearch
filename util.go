package main

func stringToNum(input string) (num int) {
	runes := []rune(input)
	e := 0
	for i := 0; i < len(runes); i++ {
		e += int(runes[i])
	}
	return e
}
