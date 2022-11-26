package main

import (
	"fmt"
	gt "github.com/bas24/googletranslatefree"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranslate(t *testing.T) {
	const text string = `Hello, World!`
	result, err := gt.Translate(text, "en", "pl")
	assert.Nil(t, err)
	fmt.Println(result)
}
