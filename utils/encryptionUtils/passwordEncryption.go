package encryptionUtils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
)

// Random generate string
func GetRandomQrCodeSecurityToken(n int) string {
	const charsToUse = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz- _+"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = charsToUse[b%byte(len(charsToUse))]
	}
	return string(bytes)
}

func GetRandomString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

// This code was copied from: http://code.google.com/p/go/source/browse/pbkdf2/pbkdf2.go?repo=crypto
func PBKDF2(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return dk[:keyLen]
}

func EncryptPasswordForDB(password string, passwordSaltLength int) string {
	salt := GetUserPasswordSalt(passwordSaltLength)
	pwd := EncodePassword(password, salt)
	// save salt and encode password, use $ as split char
	return fmt.Sprintf("%s$%s", salt, pwd)
}

// use pbkdf2 encode password
func EncodePassword(rawPwd string, salt string) string {
	pwd := PBKDF2([]byte(rawPwd), []byte(salt), 10000, 50, sha256.New)
	return hex.EncodeToString(pwd)
}

// return a user salt token
func GetUserPasswordSalt(passwordSaltLength int) string {
	return GetRandomString(passwordSaltLength)
}

func VerifyPassword(rawPwd, encodedPwd string, passwordSaltLength int) bool {
	// split
	var salt, encoded string
	if len(encodedPwd) > (passwordSaltLength + 1) {
		salt = encodedPwd[:passwordSaltLength]
		encoded = encodedPwd[(passwordSaltLength + 1):]
	}

	return EncodePassword(rawPwd, salt) == encoded
}
