package models

type User struct {
	Email               string `json:"email"`
	Username            string `json:"username"`
	Password            string `json:"password"`
	PublicKey           string `json:"publicKey"`
	EncryptedPrivateKey string `json:"encryptedPrivateKey"`
	PrivateKey          []byte `json:"privateKey"`
}
