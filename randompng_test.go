package png_test

import (
	//	. "fmt"
	"io/ioutil"
	"os"
	"png"
	"testing"
)

func Test_RandomPixel(t *testing.T) {
	l := uint32(256)
	w := uint32(256)
	pixel := png.RandomPixel(l, w)
	//  []byte{0, 255, 0, 0, 0, 255, 0, 0, 0, 255, 0, 0, 255, 0, 0, 0, 255, 255, 0, 0, 0, 0, 0, 255, 255, 0, 0, 0, 255, 0}
	p := png.InitSimplePng(l, w, png.Compress(pixel))
	os.Remove("test.png")
	if err := ioutil.WriteFile("test.png", p.ToBytes(), 0777); err != nil {
		t.Error("Fail")
	} else {
		t.Log("Pass")
	}
}
