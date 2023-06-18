package file_operations

import (
	"encoding/json"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"os"
)

const (
	PathCredentialsFile = "credentials.json"
	PathTokenFile       = "token.json"
)

type TokenStorage struct {
	Token *oauth2.Token
}

func (ts *TokenStorage) SaveToken(token *oauth2.Token) error {
	ts.Token = token
	return storeToken(token)
}

func (ts *TokenStorage) LoadToken() (*oauth2.Token, error) {
	token, err := loadToken()
	if err != nil {
		return nil, nil
	}
	ts.Token = token
	return ts.Token, nil
}

func loadToken() (*oauth2.Token, error) {
	// Open the token file
	file, err := os.Open(PathTokenFile)
	if err != nil {
		return nil, errors.Wrap(err, "token file doesn't exist or cannot be opened")
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			return
		}
	}(file)

	// Deserialize the token from JSON
	token := &oauth2.Token{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(token)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to load token")
	}

	return token, nil
}

func storeToken(token *oauth2.Token) error {
	// Create or open the token file
	file, err := os.Create(PathTokenFile)
	if err != nil {
		return errors.Wrap(err, "Unable to create token file")
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			return
		}
	}(file)

	// Serialize the token to JSON
	encoder := json.NewEncoder(file)
	err = encoder.Encode(token)
	if err != nil {
		return errors.Wrap(err, "unable to store token")
	}
	return nil
}
