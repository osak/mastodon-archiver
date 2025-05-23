package secret

import (
	"github.com/perimeterx/marshmallow"
	"os"
)

type Secret struct {
	DbHost string `json:"db_host"`
	DbPort int `json:"db_port"`
	DbUser string `json:"db_user"`
	DbPassword string `json:"db_password"`
	DbName string `json:"db_name"`
	MastodonAccessToken string `json:"mastodon_access_token"`
}

func LoadSecret(path string) (*Secret, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	secret := &Secret{}
	_, err = marshmallow.Unmarshal(buf, secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}
