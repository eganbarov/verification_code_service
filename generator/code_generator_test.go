package generator

import (
	"testing"
	"unicode/utf8"
)

func TestCodeGenerator_GenerateCode(t *testing.T) {
	codeGenerator := &CodeGenerator{}

	got := codeGenerator.GenerateCode()

	runeLength := utf8.RuneCountInString(got)
	expectedCodeLength := 6
	if runeLength != expectedCodeLength {
		t.Errorf("incorrect code length: expected %d, got %d", expectedCodeLength, runeLength)
	}
}
