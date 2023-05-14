package md5_calc

import (
	"encoding/binary"
	"math/bits"
)

func blockGeneric(dig *digest, p []byte) {
	// 加载当前状态
	a, b, c, d := dig.s[0], dig.s[1], dig.s[2], dig.s[3]

	for i := 0; i <= len(p)-BlockSize; i += BlockSize {
		// 消除p上的边界检查，确保 q 的长度为 BlockSize，并且容量也为 BlockSize
		q := p[i:]
		q = q[:BlockSize:BlockSize]

		// 保存当前状态
		aa, bb, cc, dd := a, b, c, d

		// 加载输入
		x0 := binary.LittleEndian.Uint32(q[4*0x0:])
		x1 := binary.LittleEndian.Uint32(q[4*0x1:])
		x2 := binary.LittleEndian.Uint32(q[4*0x2:])
		x3 := binary.LittleEndian.Uint32(q[4*0x3:])
		x4 := binary.LittleEndian.Uint32(q[4*0x4:])
		x5 := binary.LittleEndian.Uint32(q[4*0x5:])
		x6 := binary.LittleEndian.Uint32(q[4*0x6:])
		x7 := binary.LittleEndian.Uint32(q[4*0x7:])
		x8 := binary.LittleEndian.Uint32(q[4*0x8:])
		x9 := binary.LittleEndian.Uint32(q[4*0x9:])
		xa := binary.LittleEndian.Uint32(q[4*0xa:])
		xb := binary.LittleEndian.Uint32(q[4*0xb:])
		xc := binary.LittleEndian.Uint32(q[4*0xc:])
		xd := binary.LittleEndian.Uint32(q[4*0xd:])
		xe := binary.LittleEndian.Uint32(q[4*0xe:])
		xf := binary.LittleEndian.Uint32(q[4*0xf:])

		// round 1
		a = b + bits.RotateLeft32((((c^d)&b)^d)+a+x0+0xd76aa478, 7)
		d = a + bits.RotateLeft32((((b^c)&a)^c)+d+x1+0xe8c7b756, 12)
		c = d + bits.RotateLeft32((((a^b)&d)^b)+c+x2+0x242070db, 17)
		b = c + bits.RotateLeft32((((d^a)&c)^a)+b+x3+0xc1bdceee, 22)
		a = b + bits.RotateLeft32((((c^d)&b)^d)+a+x4+0xf57c0faf, 7)
		d = a + bits.RotateLeft32((((b^c)&a)^c)+d+x5+0x4787c62a, 12)
		c = d + bits.RotateLeft32((((a^b)&d)^b)+c+x6+0xa8304613, 17)
		b = c + bits.RotateLeft32((((d^a)&c)^a)+b+x7+0xfd469501, 22)
		a = b + bits.RotateLeft32((((c^d)&b)^d)+a+x8+0x698098d8, 7)
		d = a + bits.RotateLeft32((((b^c)&a)^c)+d+x9+0x8b44f7af, 12)
		c = d + bits.RotateLeft32((((a^b)&d)^b)+c+xa+0xffff5bb1, 17)
		b = c + bits.RotateLeft32((((d^a)&c)^a)+b+xb+0x895cd7be, 22)
		a = b + bits.RotateLeft32((((c^d)&b)^d)+a+xc+0x6b901122, 7)
		d = a + bits.RotateLeft32((((b^c)&a)^c)+d+xd+0xfd987193, 12)
		c = d + bits.RotateLeft32((((a^b)&d)^b)+c+xe+0xa679438e, 17)
		b = c + bits.RotateLeft32((((d^a)&c)^a)+b+xf+0x49b40821, 22)

		// round 2
		a = b + bits.RotateLeft32((((b^c)&d)^c)+a+x1+0xf61e2562, 5)
		d = a + bits.RotateLeft32((((a^b)&c)^b)+d+x6+0xc040b340, 9)
		c = d + bits.RotateLeft32((((d^a)&b)^a)+c+xb+0x265e5a51, 14)
		b = c + bits.RotateLeft32((((c^d)&a)^d)+b+x0+0xe9b6c7aa, 20)
		a = b + bits.RotateLeft32((((b^c)&d)^c)+a+x5+0xd62f105d, 5)
		d = a + bits.RotateLeft32((((a^b)&c)^b)+d+xa+0x02441453, 9)
		c = d + bits.RotateLeft32((((d^a)&b)^a)+c+xf+0xd8a1e681, 14)
		b = c + bits.RotateLeft32((((c^d)&a)^d)+b+x4+0xe7d3fbc8, 20)
		a = b + bits.RotateLeft32((((b^c)&d)^c)+a+x9+0x21e1cde6, 5)
		d = a + bits.RotateLeft32((((a^b)&c)^b)+d+xe+0xc33707d6, 9)
		c = d + bits.RotateLeft32((((d^a)&b)^a)+c+x3+0xf4d50d87, 14)
		b = c + bits.RotateLeft32((((c^d)&a)^d)+b+x8+0x455a14ed, 20)
		a = b + bits.RotateLeft32((((b^c)&d)^c)+a+xd+0xa9e3e905, 5)
		d = a + bits.RotateLeft32((((a^b)&c)^b)+d+x2+0xfcefa3f8, 9)
		c = d + bits.RotateLeft32((((d^a)&b)^a)+c+x7+0x676f02d9, 14)
		b = c + bits.RotateLeft32((((c^d)&a)^d)+b+xc+0x8d2a4c8a, 20)

		// round 3
		a = b + bits.RotateLeft32((b^c^d)+a+x5+0xfffa3942, 4)
		d = a + bits.RotateLeft32((a^b^c)+d+x8+0x8771f681, 11)
		c = d + bits.RotateLeft32((d^a^b)+c+xb+0x6d9d6122, 16)
		b = c + bits.RotateLeft32((c^d^a)+b+xe+0xfde5380c, 23)
		a = b + bits.RotateLeft32((b^c^d)+a+x1+0xa4beea44, 4)
		d = a + bits.RotateLeft32((a^b^c)+d+x4+0x4bdecfa9, 11)
		c = d + bits.RotateLeft32((d^a^b)+c+x7+0xf6bb4b60, 16)
		b = c + bits.RotateLeft32((c^d^a)+b+xa+0xbebfbc70, 23)
		a = b + bits.RotateLeft32((b^c^d)+a+xd+0x289b7ec6, 4)
		d = a + bits.RotateLeft32((a^b^c)+d+x0+0xeaa127fa, 11)
		c = d + bits.RotateLeft32((d^a^b)+c+x3+0xd4ef3085, 16)
		b = c + bits.RotateLeft32((c^d^a)+b+x6+0x04881d05, 23)
		a = b + bits.RotateLeft32((b^c^d)+a+x9+0xd9d4d039, 4)
		d = a + bits.RotateLeft32((a^b^c)+d+xc+0xe6db99e5, 11)
		c = d + bits.RotateLeft32((d^a^b)+c+xf+0x1fa27cf8, 16)
		b = c + bits.RotateLeft32((c^d^a)+b+x2+0xc4ac5665, 23)

		// round 4
		a = b + bits.RotateLeft32((c^(b|^d))+a+x0+0xf4292244, 6)
		d = a + bits.RotateLeft32((b^(a|^c))+d+x7+0x432aff97, 10)
		c = d + bits.RotateLeft32((a^(d|^b))+c+xe+0xab9423a7, 15)
		b = c + bits.RotateLeft32((d^(c|^a))+b+x5+0xfc93a039, 21)
		a = b + bits.RotateLeft32((c^(b|^d))+a+xc+0x655b59c3, 6)
		d = a + bits.RotateLeft32((b^(a|^c))+d+x3+0x8f0ccc92, 10)
		c = d + bits.RotateLeft32((a^(d|^b))+c+xa+0xffeff47d, 15)
		b = c + bits.RotateLeft32((d^(c|^a))+b+x1+0x85845dd1, 21)
		a = b + bits.RotateLeft32((c^(b|^d))+a+x8+0x6fa87e4f, 6)
		d = a + bits.RotateLeft32((b^(a|^c))+d+xf+0xfe2ce6e0, 10)
		c = d + bits.RotateLeft32((a^(d|^b))+c+x6+0xa3014314, 15)
		b = c + bits.RotateLeft32((d^(c|^a))+b+xd+0x4e0811a1, 21)
		a = b + bits.RotateLeft32((c^(b|^d))+a+x4+0xf7537e82, 6)
		d = a + bits.RotateLeft32((b^(a|^c))+d+xb+0xbd3af235, 10)
		c = d + bits.RotateLeft32((a^(d|^b))+c+x2+0x2ad7d2bb, 15)
		b = c + bits.RotateLeft32((d^(c|^a))+b+x9+0xeb86d391, 21)

		// add saved state
		a += aa
		b += bb
		c += cc
		d += dd
	}

	// save state
	dig.s[0], dig.s[1], dig.s[2], dig.s[3] = a, b, c, d
}