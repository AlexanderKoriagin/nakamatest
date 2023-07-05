package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/heroiclabs/nakama-common/runtime"

	entitiesHttp "github.com/akrillis/nakamatest/internal/entities/http"
	"github.com/akrillis/nakamatest/internal/service"
	"github.com/akrillis/nakamatest/internal/service/file"
	"github.com/akrillis/nakamatest/internal/storage"
	"github.com/akrillis/nakamatest/internal/storage/saver"
)

func FileChecker(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("got request with payload %s", payload)

	var req entitiesHttp.FileRequest
	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		logger.Error("couldn't parse payload %s: %v", payload, err.Error())
		return "", fmt.Errorf("couldn't parse payload: %w", err)
	}

	req.Validate()
	f := file.NewFile(req.ToPath(), req.Hash)

	content, err := getContent(f, saver.NewDB(ctx, db))
	if err != nil {
		logger.Error("couldn't read file %s: %v", req.ToPath(), err.Error())
		return "", fmt.Errorf("couldn't get file %s content: %w", req.ToPath(), err)
	}

	resp, err := json.Marshal(entitiesHttp.FileResponse{
		Type:    req.Type,
		Version: req.Version,
		Hash:    req.Hash,
		Content: content,
	})
	if err != nil {
		logger.Error("couldn't marshal response: %v", err.Error())
		return "", fmt.Errorf("couldn't marshal response: %w", err)
	}

	return string(resp), nil
}

func getContent(f service.Filer, saver storage.Saver) (string, error) {
	content, err := f.ReadWithCheck()
	if err != nil {
		return "", fmt.Errorf("couldn't get file content: %w", err)
	}

	if err = saver.Save(f.GetPath(), content); err != nil {
		log.Printf("couldn't save content from file %s: %v", f.GetPath(), err.Error())
	}

	return content, nil
}
