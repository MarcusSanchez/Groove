package db

import (
	OAuthState "GrooveGuru/ent/oauthstate"
	Session "GrooveGuru/ent/session"
	"fmt"
	"time"
)

// SpawnSessionCleaner deletes expired sessions every 24 hours.
// It is called in the init function of database.go.
// Required as sessions expire, the database still stores them.
//
// (there is no security risk to lazy clear the database, cookies expire at the
// same time therefore the session will be lost regardless)
func SpawnSessionCleaner() {
	for {
		time.Sleep(24 * time.Hour)
		affected, err := client.Session.
			Delete().
			Where(Session.ExpirationLT(time.Now())).
			Exec(ctx)
		if err != nil {
			logError("SessionCleaner[CRON]", "Worker", err)
		} else {
			fmt.Printf(
				"[%s] [SUCCESS] Session Cleared (affected: %d)\n",
				time.Now().Format("15:04:05"),
				affected,
			)
		}
	}
}

// SpawnOAuthStoreCleaner deletes expired states every 24 hours.
// It is called in the init function of database.go.
// Required as states expire without being fulfilled, meaning the database still stores them.
func SpawnOAuthStoreCleaner() {
	for {
		time.Sleep(24 * time.Hour)
		affected, err := client.OAuthState.
			Delete().
			Where(OAuthState.ExpirationLT(time.Now())).
			Exec(ctx)
		if err != nil {
			logError("OAuthStoreCleaner[CRON]", "Worker", err)
		} else {
			fmt.Printf(
				"%s [SUCCESS] OAuthStore Cleared (affected: %d)\n",
				time.Now().Format("15:04:05"),
				affected,
			)
		}
	}
}

/** helpers **/

// logError formats and prints an error with context.
func logError(fn, context string, err error) {
	fmt.Printf(
		"%s [ERROR] [Function: %s (Context: %s)] %s\n",
		time.Now().Format("15:04:05"),
		fn, context, err.Error(),
	)
}
