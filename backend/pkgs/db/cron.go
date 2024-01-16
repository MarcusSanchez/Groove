package db

import (
	"GrooveGuru/pkgs/ent"
	OAuthState "GrooveGuru/pkgs/ent/oauthstate"
	Session "GrooveGuru/pkgs/ent/session"
	"context"
	"fmt"
	"go.uber.org/fx"
	"time"
)

type Scheduler struct {
	stop    chan struct{}
	done    chan struct{}
	tickers []*time.Ticker
	client  *ent.Client
}

func InvokeScheduler(lc fx.Lifecycle, client *ent.Client) {
	scheduler := &Scheduler{
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
		tickers: []*time.Ticker{},
		client:  client,
	}
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			scheduler.Start()
			return nil
		},
		OnStop: func(context.Context) error {
			scheduler.Shutdown()
			return nil
		},
	})
}

func (s *Scheduler) ticker(d time.Duration) *time.Ticker {
	ticker := time.NewTicker(d)
	s.tickers = append(s.tickers, ticker)
	return ticker
}

func (s *Scheduler) Start() {
	go func() {
		defer close(s.done)

		ticker24h := s.ticker(24 * time.Hour)

		for {
			select {
			case <-ticker24h.C:
				go s.CleanSession()
				go s.CleanOAuthStore()
			case <-s.stop:
				return
			}
		}
	}()
}

func (s *Scheduler) Shutdown() {
	close(s.stop)
	for _, ticker := range s.tickers {
		ticker.Stop()
	}
	<-s.done
}

// CleanSession deletes expired sessions every 24 hours.
// It is called in the init function of database.go.
// Required as sessions expire, the database still stores them.
//
// (there is no security risk to lazy clear the database, cookies expire at the
// same time therefore the session will be lost regardless)
func (s *Scheduler) CleanSession() {
	affected, err := s.client.Session.
		Delete().
		Where(Session.ExpirationLT(time.Now())).
		Exec(context.Background())
	if err != nil {
		logError("SessionCleaner[CRON]", "Worker", err)
	} else {
		fmt.Printf(
			"%s [SUCCESS] Session Cleared (affected: %d)\n",
			time.Now().Format("15:04:05"),
			affected,
		)
	}
}

// CleanOAuthStore deletes expired states every 24 hours.
// It is called in the init function of database.go.
// Required as states expire without being fulfilled, meaning the database still stores them.
func (s *Scheduler) CleanOAuthStore() {
	affected, err := s.client.OAuthState.
		Delete().
		Where(OAuthState.ExpirationLT(time.Now())).
		Exec(context.Background())
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

/** helpers **/

// logError formats and prints an error with context.
func logError(fn, context string, err error) {
	fmt.Printf(
		"%s [ERROR] [Function: %s (Context: %s)] %s\n",
		time.Now().Format("15:04:05"),
		fn, context, err.Error(),
	)
}
