package saver_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/lib/pq"

	entitiesDB "github.com/akrillis/nakamatest/internal/entities/db"
	"github.com/akrillis/nakamatest/internal/storage/saver"
)

func TestDB_Save(t *testing.T) {
	db, shutdownFn := setup()
	defer shutdownFn()

	ctx := context.Background()
	s := saver.NewDB(ctx, db)

	require.NoError(t, s.Save("type0/version0.json", `{"f0": "v0", "f1": 1}`))
	require.NoError(t, s.Save("type0/version0.json", `{"f0": "v0", "f1": 1}`))
	require.NoError(t, s.Save("type1/version1.json", `{"f1": "v1", "f2": 2}`))
}

func setup() (*sql.DB, func()) {
	ctx := context.Background()

	role := "postgres"
	host, port, stopContainer := setupDocker()
	dsn := fmt.Sprintf("postgres://%s@%s:%d/postgres?sslmode=disable", role, host, port)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		logrus.WithField("dsn", dsn).WithError(err).Fatalf("failed to connect postgres: %s", err.Error())
	}

	if err = conn.PingContext(ctx); err != nil {
		logrus.WithError(err).Fatal("failed to ping postgres")
	}

	_, err = conn.Exec("CREATE DATABASE test__database123456789")
	if err != nil {
		logrus.WithError(err).Fatal("failed to create db for conversions")
	}

	dsnTarget := fmt.Sprintf("postgres://%s@%s:%d/test__database123456789?sslmode=disable", role, host, port)
	db, err := sql.Open("postgres", dsnTarget)
	if err != nil {
		logrus.
			WithError(err).
			Error("failed to pg open connection")
	}

	if _, err = db.ExecContext(ctx, entitiesDB.StmtCreateTable); err != nil {
		logrus.
			WithError(err).
			Error("failed to create target table")
	}

	shutdownFn := func() {
		_, err = conn.Exec("DROP DATABASE test__database123456789")
		if err != nil {
			logrus.WithError(err).Error("failed to drop test_task")
		}
		err = conn.Close()
		if err != nil {
			logrus.WithError(err).Fatal("failed to close connection")
		}
		stopContainer()
	}

	return db, shutdownFn
}

func setupDocker() (string, int, func()) {
	ctx := context.Background()
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:14",
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor:   wait.ForListeningPort("5432/tcp"),
			Env: map[string]string{
				"POSTGRES_HOST_AUTH_METHOD": "trust",
			},
		},
		Started: true,
	})
	if err != nil {
		logrus.WithError(err).Fatalf("failed to create postgres container")
	}

	host, err := pgContainer.Host(ctx)
	if err != nil {
		logrus.WithError(err).Fatal("failed to get host")
	}

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		logrus.WithError(err).Fatal("failed to map port")
	}

	return host, port.Int(), func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			logrus.WithError(err).Error("failed to terminate pg container")
		}
	}
}
