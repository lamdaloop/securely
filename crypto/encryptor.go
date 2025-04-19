
package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "errors"
    "io"
)

func deriveKey(password string) []byte {
    hash := sha256.Sum256([]byte(password))
    return hash[:]
}

func Encrypt(plainText []byte, password string) (cipherText, iv []byte, err error) {
    key := deriveKey(password)

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, nil, err
    }

    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return nil, nil, err
    }

    iv = make([]byte, aesGCM.NonceSize())
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, nil, err
    }

    cipherText = aesGCM.Seal(nil, iv, plainText, nil)
    return cipherText, iv, nil
}

func Decrypt(cipherText, iv []byte, password string) ([]byte, error) {
    key := deriveKey(password)

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    if len(iv) != aesGCM.NonceSize() {
        return nil, errors.New("invalid IV size")
    }

    return aesGCM.Open(nil, iv, cipherText, nil)
}
