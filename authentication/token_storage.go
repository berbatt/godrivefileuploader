package authentication

import (
	"encoding/json"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"os"
)

type TokenStorage struct {
	Token *oauth2.Token
}

func (ts *TokenStorage) saveToken(pathTokenFile string) error {
	return storeToken(ts.Token, pathTokenFile)
}

func (ts *TokenStorage) loadToken(pathTokenFile string) (*oauth2.Token, error) {
	token, err := loadToken(pathTokenFile)
	if err != nil {
		return nil, err
	}
	ts.Token = token
	return ts.Token, nil
}

func (ts *TokenStorage) isTokenExists() bool {
	return ts.Token != nil
}

func loadToken(pathTokenFile string) (*oauth2.Token, error) {
	// Open the token file
	file, err := os.Open(pathTokenFile)
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
	err = json.NewDecoder(file).Decode(token)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to load token")
	}
	return token, nil
}

func storeToken(token *oauth2.Token, pathTokenFile string) error {
	// Create or open the token file
	file, err := os.Create(pathTokenFile)
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
