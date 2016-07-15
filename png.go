package png

import (
	"compress/zlib"
	. "fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
)

type SimplePng struct {
	//var pngSignature = [8]byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}
	//In png file the first 8 bytes are filled with fixed value
	PngSignature [8]byte
	Ihdr         [25]byte
	//var srgb = [13]byte{0x00, 0x00, 0x00, 0x01, 0x73, 0x52,
	//                   0x47, 0x42, 0x00, 0xae, 0xce, 0x1c, 0xe9}
	//In my case, sRGB is fixed to make the thing simple
	Srgb [13]byte
	//var gAMA = [16]byte{0x00, 0x00, 0x00, 0x04, 0x67, 0x41, 0x4d, 0x41,
	//                   0x00, 0x00, 0xb1, 0x8f, 0x0b, 0xfc, 0x61, 0x05}
	//In my case, gMAM is fixed to make the thing simple
	Gama [16]byte
	//var pHYS = [21]byte{0x00, 0x00, 0x00, 0x09, 0x70, 0x48, 0x59,
	//                    0x73, 0x00, 0x00, 0x0e, 0xc3, 0x00, 0x00,
	//                    0x0e, 0xc3, 0x01, 0xc7, 0x6f, 0xa8, 0x64}
	//In my case, pHYS is fixed to make the thing simple
	Phys [21]byte
	//IDATA
	Idat []byte
	//var iend = [8]byte{0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82}
	//In png file the last 8 bytes are filled with fixed value
	Iend [8]byte
}

func (sp SimplePng) Write(fn string) {
	os.Remove(fn)
	var tmp []byte
	tmp = append(tmp, sp.PngSignature[:]...)
	tmp = append(tmp, sp.Ihdr[:]...)
	tmp = append(tmp, sp.Srgb[:]...)
	tmp = append(tmp, sp.Gama[:]...)
	tmp = append(tmp, sp.Phys[:]...)
	tmp = append(tmp, sp.Idat...)
	tmp = append(tmp, sp.Iend[:]...)
	Printf("%d\n", tmp)
	ioutil.WriteFile(fn, tmp, 0777)
}

//InitSimplePng will fill the PngSignature, sRGB, gAMA, pHYS, iEND with a default value
func InitSimplePng(l uint32, w uint32, pixelCompressed []byte) (sp SimplePng) {
	sp.PngSignature = [8]byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}
	sp.Ihdr = createIhdr(l, w)
	sp.Srgb = [13]byte{0x00, 0x00, 0x00, 0x01, 0x73, 0x52,
		0x47, 0x42, 0x00, 0xae, 0xce, 0x1c, 0xe9}
	sp.Gama = [16]byte{0x00, 0x00, 0x00, 0x04, 0x67, 0x41, 0x4d, 0x41,
		0x00, 0x00, 0xb1, 0x8f, 0x0b, 0xfc, 0x61, 0x05}
	sp.Phys = [21]byte{0x00, 0x00, 0x00, 0x09, 0x70, 0x48, 0x59,
		0x73, 0x00, 0x00, 0x0e, 0xc3, 0x00, 0x00,
		0x0e, 0xc3, 0x01, 0xc7, 0x6f, 0xa8, 0x64}
	sp.Idat = createIdat(pixelCompressed)
	sp.Iend = [8]byte{0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82}
	return sp
}

func createIhdr(l, w uint32) [25]byte {
	var ihdr [25]byte
	ihdr[0] = byte(13 >> 24) // 0x00
	ihdr[1] = byte(13 >> 16) // 0x00
	ihdr[2] = byte(13 >> 8)  // 0x00
	ihdr[3] = byte(13 >> 0)  // 0x13
	ihdr[4] = byte('I')
	ihdr[5] = byte('H')
	ihdr[6] = byte('D')
	ihdr[7] = byte('R')
	ihdr[8] = byte(l >> 24)
	ihdr[9] = byte(l >> 16)
	ihdr[10] = byte(l >> 8)
	ihdr[11] = byte(l >> 0)
	ihdr[12] = byte(w >> 24)
	ihdr[13] = byte(w >> 16)
	ihdr[14] = byte(w >> 8)
	ihdr[15] = byte(w >> 0)
	ihdr[16] = byte(0x08)
	ihdr[17] = byte(0x02)
	ihdr[18] = byte(0x00)
	ihdr[19] = byte(0x00)
	ihdr[20] = byte(0x00)
	crc := crc32.ChecksumIEEE(ihdr[4:20])
	ihdr[21] = byte(crc >> 24)
	ihdr[22] = byte(crc >> 16)
	ihdr[23] = byte(crc >> 8)
	ihdr[24] = byte(crc >> 0)
	//	for _, b := range ihdr {
	//		Printf("%x ", b)
	//	}
	return ihdr
}

func createIdat(pixel []byte) []byte {
	l := 4 + 4 + len(pixel) + 4
	var idat []byte
	idat = make([]byte, l)
	idat[0] = byte(uint(len(pixel)) >> 24)
	idat[1] = byte(uint(len(pixel)) >> 16)
	idat[2] = byte(uint(len(pixel)) >> 8)
	idat[3] = byte(uint(len(pixel)) >> 0)
	idat[4] = byte('I')
	idat[5] = byte('D')
	idat[6] = byte('A')
	idat[7] = byte('T')
	for i, b := range pixel {
		idat[8+i] = b
	}
	crc := crc32.ChecksumIEEE(idat[4 : l-4])
	idat[l-4] = byte(crc >> 24)
	idat[l-3] = byte(crc >> 16)
	idat[l-2] = byte(crc >> 8)
	idat[l-1] = byte(crc >> 0)
	return idat
}

//func Compress use for compress the pixel data for IDAT chunk
func Compress(byte1 []byte) []byte {
	var in1 bytes.Buffer
	w1 := zlib.NewWriter(&in1)
	w1.Write(byte1)
	w1.Close()
	return in1.Bytes()
}

//func Decompress use for decompress the IDAT chunk to the pixel data
func Decompress(byte1 []byte) []byte {
	var b1 *bytes.Buffer
	b1 = bytes.NewBuffer(byte1)
	var out1 bytes.Buffer
	r1, _ := zlib.NewReader(b1)
	io.Copy(&out1, r1)
	r1.Close()
	return out1.Bytes()
}
