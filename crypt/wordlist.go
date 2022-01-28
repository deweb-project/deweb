package crypt

import (
	_ "embed"
	"math/rand"
	"strings"
)

//go:embed words.txt
var words_int string

var Words []string

func init() {
	Words = strings.Split(words_int, "\n")
}

func GetWord() string {
	return Words[rand.Intn(len(Words))]
}
