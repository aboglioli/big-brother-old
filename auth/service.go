package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aboglioli/big-brother/pkg/config"
	"github.com/aboglioli/big-brother/pkg/errors"
	gocache "github.com/patrickmn/go-cache"
)

var cache = gocache.New(2*time.Minute, 5*time.Minute)

func Validate(token string) (*User, errors.Error) {
	if data, ok := cache.Get(token); ok {
		if user, ok := data.(*User); ok {
			return user, nil
		}
	}

	user, err := getFromApi(token)
	if err != nil {
		return nil, errors.NewInternal().SetCode("AUTH_API").SetMessage(fmt.Sprintf("Couldn't request Auth API: %s", err.Error()))
	}

	cache.Set(token, user, gocache.DefaultExpiration)

	return user, nil
}

func getFromApi(token string) (*User, error) {
	conf := config.Get()

	req, err := http.NewRequest("GET", conf.AuthURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return nil, err
	}
	defer resp.Body.Close()

	user := &User{}
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

func Invalidate(token string) {
	cache.Delete(token)
}
