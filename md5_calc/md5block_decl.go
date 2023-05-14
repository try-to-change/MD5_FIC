//go:build amd64 || 386 || arm || ppc64le || ppc64 || s390x || arm64

package md5_calc

const haveAsm = true

func block(dig *digest, p []byte) {
	blockGeneric(dig, p)
}
