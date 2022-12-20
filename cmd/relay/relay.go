package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fiatjaf/relayer"
	"github.com/fiatjaf/relayer/storage/postgresql"
	"github.com/kelseyhightower/envconfig"
	"github.com/nbd-wtf/go-nostr"
	"github.com/rs/zerolog"
)

type Relay struct {
	PostgresDatabase string `envconfig:"POSTGRESQL_DATABASE"`
	LogLevel         string `envconfig:"LOG_LEVEL"`
	storage          *postgresql.PostgresBackend
}

func (r *Relay) Name() string {
	return "FuseRelay"
}

func (r *Relay) Storage() relayer.Storage {
	return r.storage
}

func (r *Relay) Init() error {
	err := envconfig.Process("", r)
	if err != nil {
		return fmt.Errorf("failed to read from env")
	}

	// Delete old events every hour
	go func() {
		db := r.Storage().(*postgresql.PostgresBackend)
		for {
			time.Sleep(60 * time.Minute)
			db.DB.Exec(`DELETE FROM event WHERE created_at < $1`, time.Now().AddDate(0, -3, 0))
		}
	}()

	return nil
}

func (r *Relay) OnInitialized() {
	relayer.Log.Debug().Msg("Relay has been initalized")
}

func (r *Relay) AcceptEvent(evt *nostr.Event) bool {
	jsonb, _ := json.Marshal(evt)
	// Reject large events
	return len(jsonb) <= 10000
}

func (r *Relay) BeforeSave(evt *nostr.Event) {
	relayer.Log.Debug().Msgf("Saving Event: %v", evt)
}

func main() {
	r := Relay{}
	if err := envconfig.Process("", &r); err != nil {
		relayer.Log.Fatal().Err(err).Msg("failed to read from env")
	}

	level, err := zerolog.ParseLevel(r.LogLevel)
	if err == nil {
		zerolog.SetGlobalLevel(level)
	}

	relayer.Log.Debug().Msgf("DB: %v", r.PostgresDatabase)
	relayer.Log.Debug().Msgf("LOG_LEVEL: %v", r.LogLevel)
	r.storage = &postgresql.PostgresBackend{DatabaseURL: r.PostgresDatabase}

	relayer.Start(&r)
}
