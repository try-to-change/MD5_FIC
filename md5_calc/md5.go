//go:generate go run gen.go -output md5block.go

package md5_calc

import (
	"crypto"
	"encoding/binary"
	"errors"
	"hash"
)

func init() {
	crypto.RegisterHash(crypto.MD5, New)
}

// Size The size of an MD5 checksum in bytes.
const Size = 16

// BlockSize The blocksize of MD5 in bytes.
const BlockSize = 64

const (
	init0 = 0x67452301
	init1 = 0xEFCDAB89
	init2 = 0x98BADCFE
	init3 = 0x10325476
)

// digest represents the partial evaluation of a checksum.
type digest struct {
	s   [4]uint32
	x   [BlockSize]byte
	nx  int
	len uint64
}

func (d *digest) Reset() {
	d.s[0] = init0
	d.s[1] = init1
	d.s[2] = init2
	d.s[3] = init3
	d.nx = 0
	d.len = 0
}

const (
	magic         = "md5\x01"
	marshaledSize = len(magic) + 4*4 + BlockSize + 8
)

func (d *digest) MarshalBinary() ([]byte, error) {
	b := make([]byte, marshaledSize)
	//使用内置的编码/二进制函数简化读取和写入二进制数据的代码。
	binary.BigEndian.PutUint32(b[4:], d.s[0])
	binary.BigEndian.PutUint32(b[8:], d.s[1])
	binary.BigEndian.PutUint32(b[12:], d.s[2])
	binary.BigEndian.PutUint32(b[16:], d.s[3])
	copy(b[20:], d.x[:d.nx])
	binary.BigEndian.PutUint64(b[20+BlockSize:], d.len)
	copy(b[:4], []byte(magic))
	return b, nil
}

func (d *digest) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic) || string(b[:len(magic)]) != magic {
		return errors.New("invalid hash state identifier")
	}
	if len(b) != marshaledSize {
		return errors.New("invalid hash state size")
	}
	d.s[0] = binary.BigEndian.Uint32(b[4:])
	d.s[1] = binary.BigEndian.Uint32(b[8:])
	d.s[2] = binary.BigEndian.Uint32(b[12:])
	d.s[3] = binary.BigEndian.Uint32(b[16:])
	copy(d.x[:], b[20:20+BlockSize])
	d.len = binary.BigEndian.Uint64(b[20+BlockSize:])
	d.nx = int(d.len % BlockSize)
	return nil
}

// New returns a new hash.Hash computing the MD5 checksum.
func New() hash.Hash {
	d := new(digest)
	d.Reset()
	return d
}

func (d *digest) Size() int { return Size }

func (d *digest) BlockSize() int { return BlockSize }

func (d *digest) Write(p []byte) (nn int, err error) {
	nn = len(p)
	d.len += uint64(nn)
	if d.nx > 0 {
		n := copy(d.x[d.nx:], p)
		d.nx += n
		if d.nx == BlockSize {
			block(d, d.x[:])
			d.nx = 0
		}
		p = p[n:]
	}
	if len(p) >= BlockSize {
		n := len(p) &^ (BlockSize - 1)
		block(d, p[:n])
		p = p[n:]
	}
	if len(p) > 0 {
		d.nx = copy(d.x[:], p)
	}
	return
}

func (d *digest) Sum(in []byte) []byte {
	// 复制一份d，这样调用时就可以继续写入求和。
	d0 := *d
	hash := d0.checkSum()
	return append(in, hash[:]...)
}

func (d *digest) checkSum() [Size]byte {
	// 填充，添加1位和0位，直到56字节mod 64。
	tmp := [1 + 63 + 8]byte{0x80}
	pad := (55 - d.len) % 64
	binary.LittleEndian.PutUint64(tmp[1+pad:], d.len<<3)
	d.Write(tmp[:1+pad+8])

	if d.nx != 0 {
		panic("d.nx != 0")
	}

	var digest [Size]byte
	binary.LittleEndian.PutUint32(digest[0:], d.s[0])
	binary.LittleEndian.PutUint32(digest[4:], d.s[1])
	binary.LittleEndian.PutUint32(digest[8:], d.s[2])
	binary.LittleEndian.PutUint32(digest[12:], d.s[3])
	return digest
}

// Sum 返回数据的MD5校验和
func Sum(data []byte) [Size]byte {
	var d digest
	d.Reset()
	d.Write(data)
	return d.checkSum()
}
