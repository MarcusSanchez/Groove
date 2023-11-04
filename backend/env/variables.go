package env

import (
	"github.com/joho/godotenv"
	"os"
)

var (
	Port   string
	IsProd bool
	Secure bool
	PgURI  string
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
	PgURI = os.Getenv("PG_URI")
}
