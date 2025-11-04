package auth

import (
	"github.com/go-chi/jwtauth"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var JWT_SECRET_KEY string

func init() {
	v := viper.New()
	v.SetConfigFile("config.yaml")

	if err := v.ReadInConfig(); err != nil {
		log.Printf("read config: %v", err)
	}

	if err := v.UnmarshalKey("jwt_secret_key", &JWT_SECRET_KEY); err != nil {
		log.Printf("unmarshal: %v", err)
	}
}

func GenerateAuthToken() *jwtauth.JWTAuth {
	tokenAuth := jwtauth.New("HS256", []byte(JWT_SECRET_KEY), nil)
	return tokenAuth
}
