package main

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"html/template"
	"image/color"
	"net/http"
	"os"
	"time"

	"github.com/apex/gateway/v2"
	"github.com/apex/log"
	jsonhandler "github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Record struct {
	ID      string     `dynamodbav:"id" json:"id"`
	Created time.Time  `dynamodbav:"created,unixtime" json:"created"`
	Expires *time.Time `dynamodbav:"expires,unixtime,omitempty" json:"expires"`
	Color   string     `dynamodbav:"color" json:"color"`
}

type server struct {
	router *http.ServeMux
	client *dynamodb.Client
	auth   *Authenticator
	store  *sessions.CookieStore
}

func (record *Record) TimeSinceCreation() string {
	return time.Since(record.Created).String()
}

func (record *Record) TimeUntilExpiry() string {
	if record.Expires != nil {
		return time.Until(*record.Expires).String()
	}
	return ""
}

func (record *Record) TransparentBG() template.CSS {
	var c color.RGBA
	var err error
	switch len(record.Color) {
	case 7:
		_, err = fmt.Sscanf(record.Color, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(record.Color, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")
	}
	if err != nil {
		log.WithError(err).Fatal("converting to rgba")
	}
	// return fmt.Sprintf("rgba(%d, %d, %d, .5)", c.R, c.G, c.B)
	log.WithFields(log.Fields{
		"r":   c.R,
		"g":   c.G,
		"b":   c.B,
		"hex": record.Color,
	}).Debug("converted to rgba")
	return template.CSS(fmt.Sprintf("rgba(%d, %d, %d, 0.5)", c.R, c.G, c.B))
}

func newServer(local bool) *server {
	s := &server{router: &http.ServeMux{}}

	gob.Register(map[string]interface{}{})
	s.store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	auth, err := New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	s.auth = auth

	if local {
		log.SetHandler(text.Default)
		log.Info("local mode")
		s.client = dynamoLocal()
	} else {
		log.SetHandler(jsonhandler.Default)
		log.Info("cloud mode")
		s.client = dynamoCloud()
	}

	s.router.Handle("/", s.list())
	s.router.Handle("/login", s.login())
	s.router.Handle("/logout", s.logout())
	s.router.Handle("/user", s.user())
	s.router.Handle("/callback", s.callback())
	s.router.Handle("/add", s.add())

	// k

	return s
}

func main() {
	_, awsDetected := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME")
	log.WithField("awsDetected", awsDetected).Info("starting up")
	s := newServer(!awsDetected)

	var err error

	if awsDetected {
		log.Info("starting cloud server")
		err = gateway.ListenAndServe("", s.router)
	} else {
		err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), s.router)
	}
	log.Info(".... starting")
	log.WithError(err).Fatal("error listening")
	log.Info(".... ending")
}

// Authenticator is used to authenticate our users.
type Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

// New instantiates the *Authenticator.
func New() (*Authenticator, error) {
	provider, err := oidc.NewProvider(
		context.Background(),
		"https://"+os.Getenv("AUTH0_DOMAIN")+"/",
	)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
	}, nil
}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}
