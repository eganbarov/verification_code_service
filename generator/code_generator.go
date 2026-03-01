package generator

import (
	"math/rand/v2"
	"strconv"
)

type CodeGenerator struct {
}

func (c *CodeGenerator) GenerateCode() string {
	code := rand.IntN(900000) + 100000
	return strconv.Itoa(code)
}
