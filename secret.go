package main

import (
	"github.com/perimeterx/marshmallow"
	"os"
)

type Secret struct {
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
