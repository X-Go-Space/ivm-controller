package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"ivm-controller/initEnv"
)

var PublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwLKZjpq1EzpiXbJhkOmm
Ap9sV3P2SRP/KdB2ZLqS+I7ONkCiLpB2EJERhkSfs0zOgZMfSqD7wOtwCGOXxjfv
UcHV81cHwAx4lfqZvt2BrJyQYJrIwqksv1iN4eZTsqJwxH/CjCHI7ns5WclzViYF
6iZd1YBGWSez2df29VaNABWlRD3PFj7CspPcjb9P3Ei2ZxRGsuxwAIhHEqBqw8Ms
uJjA1mBSt38O61Na9vXnjmobI1lQi0jcun93WfYlSpv/+HrsZl+meKQeHHwsSSws
UDXp/LYBxOIx1APpHOsy79PvZCn5mekukFmxPMYbIXhSZkk37O0KXrh62WawHRmI
wwIDAQAB
-----END PUBLIC KEY-----`


var PrivateKeyString = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAwLKZjpq1EzpiXbJhkOmmAp9sV3P2SRP/KdB2ZLqS+I7ONkCi
LpB2EJERhkSfs0zOgZMfSqD7wOtwCGOXxjfvUcHV81cHwAx4lfqZvt2BrJyQYJrI
wqksv1iN4eZTsqJwxH/CjCHI7ns5WclzViYF6iZd1YBGWSez2df29VaNABWlRD3P
Fj7CspPcjb9P3Ei2ZxRGsuxwAIhHEqBqw8MsuJjA1mBSt38O61Na9vXnjmobI1lQ
i0jcun93WfYlSpv/+HrsZl+meKQeHHwsSSwsUDXp/LYBxOIx1APpHOsy79PvZCn5
mekukFmxPMYbIXhSZkk37O0KXrh62WawHRmIwwIDAQABAoIBAF49+wPXff+tajZQ
646n9t0Jgz8yI52R/hVBMuYIqaCOlyPJcIg9dsCbcmqsXT6frc+JWKBzIy0y+FPi
AXScHptppW2hftTaRI91RIQoaSc2WxYkHVO20X+zm5CDySNwdp5jrY5DZ7Xa5i2X
bKURa4scwH+OgNlRpZBBIzLR5ZvIIQyU5sxnW6wQtJ+4qWjAzqwf32oHrtE2aVrV
escEeXt2YtuPLOF29sQqEa6VXnMyW86fIaTeAuAs8lL42krnmu/qWNeyhWR/og44
50tge/PulA+/ZYuYyz2GAnnz5GuV7jooYFhpWk2TECjyJXdTXijeLjD9Ui9yjPdv
E/FVe8ECgYEA4E+4uQLO5r/F71vLad/zCQ6Ef84hugDMyPebPqSmsZSUXLqIyw+r
PPz25KtU0RaJy3eGZ5Jvj6mnqVstVudSvGVGn9/umGkNKoy3vkPtVm912TvEXO9K
R9Gn9+VMVHWq4IlwglM/rpcsIpp5bgILoHa25e17EaAGyzb+jddmxbkCgYEA2+uP
DxnPMRAQ2i6/Zg3jb3R+8agxi4hlf/LWvSlvu9YIbg+fGyee+d+H1TU9wL895SYA
d0TS+jT6SeP8jU/AMSURqUEOMueld6bwH2TeKnkfVa3X4iaYVrXFBP5RmYFk0YNr
6G4KYHL8Lvs5fPcZdVxBmFJ7b1R+JoeBkcy7QFsCgYArDqKCwQs+N+mivJgbRqW8
Q1Ejx0mqDqVAnmbqa2ikBcVE13mSoPtZxaUO1+R8Djt9FwBxuSY5CXPpilr1p4m2
KCqaXb3K+79PP5u1pgxU3yhb/qD+xeAYUSJQ727rd3rJhxhVq+05ckNCkSxl9XaN
4rvQ/vj0tScYswHB8GsF4QKBgAm2FzUlgJ68BOJ9mfoZtudfD5QAR1/QABtsT8s+
ny5+PxUNH4uFbmG+WzMxDK8MQuFxkieyJFbkLAFDTg23bdc9uc/tjYD19bqY5pWc
UKszegzAhn34ElYR5MdZq6TJr/gIg6VZ5p9ntHcmpN091CP4lPTy/3xlEMUGytPz
ZHltAoGBAJhEDFYC4aiW+3+5EQ1JPB0bDVrOA/Mt0iH+pXB7EEQxhij7g9zVEiFQ
otgmVUXwkYMvHID9fQF4EjC5wRl8AEF6rpoR7LA4eAKH2vaKkPj2Iu6OoEQ2Zfp/
yF7WDTva9AwFDUzgh4CdKYKI+PU5Y6GmtJdfDmjW4eoBiENf430/
-----END RSA PRIVATE KEY-----`

func Encrypt(data []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(PublicKey))
	if block == nil {
		return nil, fmt.Errorf("Failed to decode PEM block containing public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("Invalid public key type")
	}

	return rsa.EncryptPKCS1v15(rand.Reader, rsaPubKey, data)
}

func Decrypt(data string) (string, error) {

	block, _ := pem.Decode([]byte(PrivateKeyString))
	if block == nil {
		initEnv.Logger.Error("Failed to decode PEM block containing private key")
		return "", fmt.Errorf("Failed to decode PEM block containing private key")
	}

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		initEnv.Logger.Error("base 64 decrypt failed, err: ", err)
		return "", err
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		initEnv.Logger.Error("x509 decrypt failed, err: ", err)
		return "", err
	}

	byteData, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, decoded)
    if err != nil {
		initEnv.Logger.Error("rsa decrypt failed, err: ", err)
		return "", err
	}
	return string(byteData), nil
}

func GenerateSalt(length int) (string, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	encodedSalt := base64.StdEncoding.EncodeToString(salt)
	return encodedSalt, nil
}

func HashPasswordWithSalt(password, salt string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password + salt))
	hash := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	return hash
}
