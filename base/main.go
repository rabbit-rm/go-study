package main

import (
	"fmt"
	"strings"
)

type LZ77Token struct {
	Offset int
	Length int
	Next   byte
}

func LZ77Compress(input string) []LZ77Token {
	var result []LZ77Token
	for i := 0; i < len(input); {
		maxLength := 0
		maxOffset := 0
		for offset := 1; offset <= i && offset <= 255; offset++ {
			length := 0
			for i+length < len(input) && length < 255 && input[i+length] == input[i+length-offset] {
				length++
			}
			if length > maxLength {
				maxLength = length
				maxOffset = offset
			}
		}
		i += maxLength
		next := byte(0)
		if i < len(input) {
			next = input[i]
			i++
		}
		result = append(result, LZ77Token{maxOffset, maxLength, next})
	}
	return result
}

func LZ77Decompress(tokens []LZ77Token) string {
	var result strings.Builder
	for _, token := range tokens {
		for i := 0; i < token.Length; i++ {
			b := result.String()[result.Len()-token.Offset]
			result.WriteByte(b)
		}
		result.WriteByte(token.Next)
	}
	return result.String()
}

func main() {
	input := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpcCI6IjE5Mi4xNjguMS4yMzMiLCJmaXhlZCI6IiIsImV4cCI6MTcxNTgyNzEyMH0.EU_A1GDMPOKYq-5WAATxk4pBgsSalj24U2hzojYGPoU"
	tokens := LZ77Compress(input)
	fmt.Println("Compressed:", tokens)
	output := LZ77Decompress(tokens)
	fmt.Println("Decompressed:", output)
}
