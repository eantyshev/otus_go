package main

import (
	"fmt"
	"strings"
)

func runeMultiply(x rune, mul int) string {
	if mul == 0 {
		// -1 means no multiplication
		mul = 1
	}
	return strings.Repeat(string(x), mul)
}

func appendLastIfDefined(res *string, lastRune rune, mul int) {
    if lastRune != 0 {
        *res = *res + runeMultiply(lastRune, mul)
    }
}


func Unpack(s string) string {
	var mul int = 0
	var lastRune int32
    var isEscaped bool
	var res string
	for _, x := range s {
		isDigit := '0' <= x && x <= '9'
        switch {
        case isEscaped:
            appendLastIfDefined(&res, lastRune, mul)
            lastRune = x
            isEscaped = false
        case x == '\\':
            isEscaped = true
        case isDigit:
			mul = 10*mul + int(x) - '0'
		case lastRune == 0:
			lastRune = x
        default:
            appendLastIfDefined(&res, lastRune, mul)
            lastRune = x
            mul = 0
        }
	}
    appendLastIfDefined(&res, lastRune, mul)
	return res
}

func main() {
	var s, res string
	fmt.Scanln(&s)
	res = Unpack(s)
	fmt.Println(res)
}
