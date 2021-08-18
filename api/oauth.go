package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"git.fuyu.moe/fuyu/router"
	"github.com/jobstoit/website/model"
	"golang.org/x/oauth2"
)

type oa struct {
	config      *oauth2.Config
	state       string
	userInfoURL string
}

func newOauth(port int, stateString, clientID, clientSecret, endpoint string) *oa {
	x := new(oa)

	res, err := http.Get(fmt.Sprintf("%s/.well-known/openid-configuration", endpoint))
	if err != nil {
		log.Fatalf("couldn't get the authorization endpoint: %v", err)
	}

	defer res.Body.Close() // nolint: errcheck

	var payload struct {
		AuthorizationEndpoint string `json:"authorization_endpoint"`
		TokenEndpoint         string `json:"token_endpoint"`
		UserinfoEndpoint      string `json:"userinfo_endpoint"`
	}

	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		log.Fatalf("couldn't read the authorization endpoint configuration: %v", err)
	}

	x.state = stateString
	x.config = &oauth2.Config{
		RedirectURL:  fmt.Sprintf("http://localhost:%d/callback", port),
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   payload.AuthorizationEndpoint,
			TokenURL:  payload.TokenEndpoint,
			AuthStyle: oauth2.AuthStyleInHeader,
		},
	}

	x.userInfoURL = payload.UserinfoEndpoint

	return x
}

func (x oa) login(ctx *router.Context) error {
	url := x.config.AuthCodeURL(x.state)
	return ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (x oa) oauthCallback(ctx *router.Context) error {
	state, code := ctx.Request.FormValue("state"), ctx.Request.FormValue("code")

	if state != x.state {
		return ctx.Redirect(http.StatusUnauthorized, "/")
	}

	token, err := x.config.Exchange(ctx.Request.Context(), code)
	if err != nil {
		return ctx.Redirect(http.StatusUnauthorized, "/")
	}

	buf := new(bytes.Buffer)
	if err := json.NewDecoder(buf).Decode(token); err != nil {
		return ctx.Redirect(http.StatusUnauthorized, "/")
	}

	http.SetCookie(ctx.Response, &http.Cookie{
		Name:  `token`,
		Value: buf.String(),
	})

	return ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func (x oa) GetUserInfo(ctx *router.Context) (model.UserInfo, error) {
	var user model.UserInfo
	token, err := ctx.Request.Cookie(`token`)
	if err != nil {
		return user, err
	}

	tk := new(oauth2.Token)
	if err := json.Unmarshal([]byte(token.Value), tk); err != nil {
		return user, err
	}

	cli := x.config.Client(ctx.Request.Context(), tk)
	res, err := cli.Get(x.userInfoURL)
	if err != nil {
		return user, err
	}

	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}
