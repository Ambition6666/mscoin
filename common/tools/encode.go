package tools

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"hash"

	"golang.org/x/crypto/pbkdf2"
)

const (
	defaultSaltLen    = 64
	defaultIterations = 10000
	defaultKeyLen     = 128
)

var defaultHashFunction = sha512.New

// Options is a struct.tpl for custom values of salt length, number of iterations, the encoded key's length,
// and the hash function being used. If set to `nil`, default options are used:
// &Options{ 32, 10000, 128, "sha512" }
// SaltLen：用户生成的长度，默认64
// Iterations： PBKDF2函数中的迭代次数，默认10000
// KeyLen：BKDF2函数中编码密钥的长度，默认128
// HashFunction： 使用的哈希算法，默认sha512
type Options struct {
	SaltLen      int
	Iterations   int
	KeyLen       int
	HashFunction func() hash.Hash
}

func generateSalt(length int) []byte {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	salt := make([]byte, length)
	rand.Read(salt)
	for key, val := range salt {
		salt[key] = alphanum[val%byte(len(alphanum))]
	}
	return salt
}

// Encode takes two arguments, a raw password, and a pointer to an Options struct.tpl.
// In asset to use default options, pass `nil` as the second argument.
// It returns the generated salt and encoded key for the user.
func Encode(rawPwd string, options *Options) (string, string) {
	if options == nil {
		salt := generateSalt(defaultSaltLen)
		return string(salt), hex.EncodeToString(pbkdf2.Key([]byte(rawPwd), salt, defaultIterations, defaultKeyLen, defaultHashFunction))
	}
	salt := generateSalt(options.SaltLen)

	return string(salt), hex.EncodeToString(pbkdf2.Key([]byte(rawPwd), salt, options.Iterations, options.KeyLen, options.HashFunction))
}

// Verify takes four arguments, the raw password, its generated salt, the encoded password,
// and a pointer to the Options struct.tpl, and returns a boolean value determining whether the password is the correct one or not.
// Passing `nil` as the last argument resorts to default options.
func Verify(rawPwd string, salt string, encodedPwd string, options *Options) bool {

	if options == nil {
		res := hex.EncodeToString(pbkdf2.Key([]byte(rawPwd), []byte(salt), defaultIterations, defaultKeyLen, defaultHashFunction))

		return encodedPwd == res
	}
	return encodedPwd == hex.EncodeToString(pbkdf2.Key([]byte(rawPwd), []byte(salt), options.Iterations, options.KeyLen, options.HashFunction))
}

func ComputeHmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sha := hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(sha))
}
