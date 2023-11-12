package db

import (
	"GrooveGuru/ent/session"
	"fmt"
	"log"
	"time"
)

func logError(fn, context string, err error) {
	fmt.Printf(
		"%s [ERROR] [Function: %s (Context: %s)] %s\n",
		time.Now().Format("15:04:05"),
		fn, context, err.Error(),
	)
}

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
			Where(session.ExpirationLT(time.Now())).
			Exec(ctx)
		if err != nil {
			logError("SessionCleaner[CRON]", "Worker", err)
		} else {
			log.Printf(
				"[%s] [SUCCESS] Session Cleared (affected: %d)\n",
				time.Now().Format("2006-01-02 15:04:05"),
				affected,
			)
		}
	}
}
