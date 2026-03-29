package generator

import (
	"crypto/rand"
	"math/big"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
const length = 10

type Generator interface {
	Generate() string
}

type generator struct{}

func New() Generator {
	return &generator{}
}

func (g *generator) Generate() string {
	res := make([]byte, length)

	for i := range res {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		res[i] = alphabet[n.Int64()]
	}

	return string(res)
}
