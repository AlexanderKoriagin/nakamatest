package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"

	entitiesDB "github.com/akrillis/nakamatest/internal/entities/db"
	"github.com/akrillis/nakamatest/internal/server/handlers"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	logger.Debug("FileHandler module is starting...")

	_, err := db.ExecContext(ctx, entitiesDB.StmtCreateTable)
	if err != nil {
		logger.Error("Error creating table: %v", err.Error())
		return fmt.Errorf("error creating table: %w", err)
	}

	if err = initializer.RegisterRpc("FileChecker", handlers.FileChecker); err != nil {
		logger.Error("Error registering RPC: %v", err.Error())
		return fmt.Errorf("error registering RPC: %w", err)
	}

	logger.Info("FileHandler module loaded")

	return nil
}
