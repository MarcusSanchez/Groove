package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Env struct {
	Port          string
	IsProd        bool
	Secure        bool
	SameSite      string
	FrontendURL   string
	BackendURL    string
	PgURI         string
	SpotifyClient string
	SpotifySecret string
}

func ProvideEnvVars() *Env {
	// if PROD is already set, we don't need to load the .env file.
	// if in production, variables should already be set.
	if os.Getenv("PROD") == "" {
		err := godotenv.Load("./pkgs/env/.env")
		if err != nil {
			panic("Environment error: " + err.Error())
		}
	}

	env := &Env{
		Port:   os.Getenv("PORT"),
		IsProd: os.Getenv("PROD") == "true",
		// (variable for cookies) will be secure(true) if not in development;
		// this is because localhost doesn't support https.
		Secure: os.Getenv("SECURE") == "true",
		// (variable for cookies) will be Lax if in production, None if in development;
		// this is because during development, the frontend and backend are on different ports.
		SameSite: os.Getenv("SAME_SITE"),
		PgURI:    os.Getenv("PG_URI"),
		// (required only for development) used to set CORS.
		FrontendURL:   os.Getenv("FRONTEND_URL"),
		BackendURL:    os.Getenv("BACKEND_URL"),
		SpotifyClient: os.Getenv("SPOTIFY_CLIENT"),
		SpotifySecret: os.Getenv("SPOTIFY_SECRET"),
	}

	env.validate()
	return env
}

func (Env) validate() {
	var errors []string
	variables := []string{"PORT", "PROD", "SECURE", "SAME_SITE", "PG_URI", "SPOTIFY_CLIENT", "SPOTIFY_SECRET", "BACKEND_URL"}
	for _, variable := range variables {
		if os.Getenv(variable) == "" {
			errors = append(errors, variable+" is not set")
		}
	}
	if len(errors) > 0 {
		statement := fmt.Sprintln("Environment error(s): ", errors)
		panic(statement)
	}
}
