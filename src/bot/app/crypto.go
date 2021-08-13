package app

import (
	"crypto/rc4"
	"fw/src/core"
)

// ============================================================================

type Rc4 rc4.Cipher

// ============================================================================

func NewRc4() *Rc4 {
	c, err := rc4.NewCipher([]byte{113, 120, 139, 215, 248, 243, 166, 253, 81, 142, 114, 191, 138, 104, 167, 222})
	if err != nil {
		core.Panic("create rc4 failed:", err)
	}

	return (*Rc4)(c)
}

func (self *Rc4) Encrypt(dst, src []byte) {
	(*rc4.Cipher)(self).XORKeyStream(dst, src)
}

func (self *Rc4) Decrypt(dst, src []byte) {
	(*rc4.Cipher)(self).XORKeyStream(dst, src)
}
