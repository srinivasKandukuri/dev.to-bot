package utils

import (
	"fmt"
	"math/rand"
	"time"
	"unicode"
)

func WaitToEnd(num int64) {
	time.Sleep(time.Duration(num) * time.Second)
}

func AssertEquals(a string, b string, msg string) {
	if a != b {
		fmt.Printf("Failed on -> %s: \"%v\" != \"%v\"\n", msg, a, b)
		panic("Assertion Error")
	}
}

func EllipticalTruncate(text string, maxLen int) string {
	lastSpaceIx := maxLen
	len := 0
	for i, r := range text {
		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}
		len++
		if len > maxLen {
			return text[:lastSpaceIx] + "..."
		}
	}
	return text
}

func GenerateRandomTitle(titles []string) string {
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(titles)
	return titles[n]
}
