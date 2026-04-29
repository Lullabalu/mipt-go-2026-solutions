//go:build !solution

package blowfish

// #cgo pkg-config: libcrypto
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <openssl/blowfish.h>
import "C"

import (
	"unsafe"
)

type Blowfish struct {
	Key C.BF_KEY
}

func New(key []byte) *Blowfish {
	bl := Blowfish{}
	C.BF_set_key(&bl.Key, C.int(len(key)), (*C.uchar)(unsafe.Pointer(&key[0])))
	return &bl
}

func (b *Blowfish) Encrypt(out, in []byte) {
	C.BF_ecb_encrypt((*C.uchar)(unsafe.Pointer(&in[0])), (*C.uchar)(unsafe.Pointer(&out[0])), &b.Key, C.BF_ENCRYPT)
}

func (b *Blowfish) Decrypt(out, in []byte) {
	C.BF_ecb_encrypt((*C.uchar)(unsafe.Pointer(&in[0])), (*C.uchar)(unsafe.Pointer(&out[0])), &b.Key, C.BF_DECRYPT)
}

func (b *Blowfish) BlockSize() int {
	return 8
}
