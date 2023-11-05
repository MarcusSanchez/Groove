package env

import (
	"github.com/joho/godotenv"
	"os"
)

var (
	Port     string
	IsProd   bool
	Secure   bool
	SameSite string
	PgURI    string
)

func init() {
	err := godotenv.Load("./env/.env")
	if err != nil {
		panic("Environment error: " + err.Error())
	}

	Port = os.Getenv("PORT")
	IsProd = os.Getenv("PROD") == "true"
	// (variable for cookies) will be secure(true) if not in development;
	// this is because localhost doesn't support https.
	Secure = os.Getenv("SECURE") == "true"
	// (variable for cookies) will be Lax if in production, None if in development;
	// this is because during development, the frontend and backend are on different ports.
	SameSite = os.Getenv("SAME_SITE")
	PgURI = os.Getenv("PG_URI")
}
