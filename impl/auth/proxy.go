package auth

import (
	"encoding/json"
	"net/http"

	implCache "github.com/aboglioli/big-brother/impl/cache"
	"github.com/aboglioli/big-brother/pkg/cache"
	"github.com/aboglioli/big-brother/pkg/config"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/user"
)

type UserProxy struct {
	cache cache.Cache
}

func NewUserProxy() *UserProxy {
	c := implCache.InMemory()
	return &UserProxy{c}
}

func (p *UserProxy) Validate(token string) (*user.User, error) {
	if data := p.cache.Get(token); data != nil {
		if user, ok := data.(*user.User); ok {
			return user, nil
		}
	}

	user, err := getFromApi(token)
	if err != nil {
		return nil, errors.NewInternal("AUTH_API").SetRef(err)
	}

	p.cache.Set(token, user, implCache.DefaultExpiration)

	return user, nil
}

func getFromApi(token string) (*user.User, error) {
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

	user := &user.User{}
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (p *UserProxy) Invalidate(token string) {
	p.cache.Delete(token)
}
