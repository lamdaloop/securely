
package storage

import (
    "encoding/gob"
    "os"
    "path/filepath"
    "secretbox/models"
)

const SecretDir = "./secrets"

func SaveSecret(secret models.Secret) error {
    path := filepath.Join(SecretDir, secret.ID + ".bin")
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := gob.NewEncoder(file)
    return encoder.Encode(secret)
}

func LoadSecret(id string) (*models.Secret, error) {
    path := filepath.Join(SecretDir, id + ".bin")
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var secret models.Secret
    decoder := gob.NewDecoder(file)
    err = decoder.Decode(&secret)
    return &secret, err
}

func DeleteSecret(id string) error {
    path := filepath.Join(SecretDir, id + ".bin")
    return os.Remove(path)
}
