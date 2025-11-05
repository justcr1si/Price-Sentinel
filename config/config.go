package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"price_sentinel/internal/api/auth"
	"price_sentinel/internal/models"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
	// "github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string   `mapstructure:"env"`
	Server   Server   `mapstructure:"server"`
	Database Database `mapstructure:"database"`
}

type Server struct {
	Router      *chi.Mux
	DB          *sql.DB
	AuthToken   *jwtauth.JWTAuth
	Address     string        `mapstructure:"address"`
	Timeout     time.Duration `mapstructure:"timeout"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}

type Database struct {
	Port       int    `mapstructure:"port"`
	Host       string `mapstructure:"host"`
	Name       string `mapstructure:"name"`
	User       string `mapstructure:"postgres"`
	Password   string `mapstructure:"password"`
	ConnString string `mapstructure:"conn_string"`
}

func (server *Server) MountMiddleware() {
	server.Router.Use(middleware.Logger)
}

func (server *Server) MountAuthHandlers() {
	server.Router.Route("/auth", func(authRouter chi.Router) {
		authRouter.Post("/signup", auth.SignUp)
		authRouter.Post("/login", auth.Login)
		authRouter.Post("/refresh", auth.Refresh)

		authRouter.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(server.AuthToken))
			r.Use(jwtauth.Authenticator)
			r.Get("/{id}", server.GetUser)
		})
	})
}

func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {

	// We get the 'id' from URL parameters of the request
	id := chi.URLParam(r, "id")
	userId, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please provide the correct input!!"))
		return
	}

	/* After the Verifier and Authenticator have successful validated this request
	 * We destructure the claims from the request and get the userId from claims
	 * We then check whether the userId from claims is same as the userId for which
	 * the request has been hit (from url params), if not that means user is using
	 * different JWT token and hence unauthorized.
	 */
	_, claims, _ := jwtauth.FromContext(r.Context())
	userIdFromClaims := int(claims["id"].(float64))

	if userId != userIdFromClaims {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("You're not authorized >("))
		return
	}

	query := `SELECT * FROM User WHERE id = ?`

	rows, err := server.DB.Query(query, userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please provide the correct input!!"))
		return
	}

	var user *models.User
	for rows.Next() {
		user, err = models.ScanRow(rows)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something bad happened on the server :("))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func CreateServer(db *sql.DB) *Server {
	return &Server{
		Router:    chi.NewRouter(),
		DB:        db,
		AuthToken: auth.GenerateAuthToken(),
	}
}

func Load() (*Config, error) {
	return LoadFrom("config.yaml")
}

func LoadFrom(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var c Config
	if err := v.Unmarshal(&c, func(dc *mapstructure.DecoderConfig) {
		dc.DecodeHook = mapstructure.StringToTimeDurationHookFunc()
	}); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	c.Database.ConnString = fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)

	return &c, nil
}

func MustLoad() *Config {
	c, err := Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	return c
}
