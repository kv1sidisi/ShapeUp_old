package models

type User struct {
	ID       []byte
	Username string
	PassHash []byte
}
