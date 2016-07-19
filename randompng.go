package png

import (
	//	. "fmt"
	"math/rand"
)

func RandomPixel(l, w uint32) []byte {
	var p []byte
	for i := uint32(0); i < l; i++ {
		p = append(p, 0)
		for j := uint32(0); j < w; j++ {
			p = append(p, byte(rand.Int31n(256)))
			p = append(p, byte(rand.Int31n(256)))
			p = append(p, byte(rand.Int31n(256)))
		}
	}
	return p
}
