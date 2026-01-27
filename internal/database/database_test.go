package database

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/dlsu-lscs/lscs-core-api/internal/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testCfg *config.Config

func mustStartMySQLContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	var (
		dbName = "database"
		dbPwd  = "password"
		dbUser = "user"
	)

	dbContainer, err := mysql.Run(context.Background(),
		"mysql:8.0.36",
		mysql.WithDatabase(dbName),
		mysql.WithUsername(dbUser),
		mysql.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(wait.ForLog("port: 3306  MySQL Community Server - GPL").WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "3306/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	// create test config
	testCfg = &config.Config{
		DBDatabase: dbName,
		DBPassword: dbPwd,
		DBUsername: dbUser,
		DBHost:     dbHost,
		DBPort:     dbPort.Port(),
		JWTSecret:  "test-secret",
		GoEnv:      "test",
		LogLevel:   "info",
	}

	return dbContainer.Terminate, err
}

func TestMain(m *testing.M) {
	teardown, err := mustStartMySQLContainer()
	if err != nil {
		log.Fatalf("could not start mysql container: %v", err)
	}

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown mysql container: %v", err)
	}
}

func TestNew(t *testing.T) {
	// reset singleton for test
	dbInstance = nil
	srv := New(testCfg)
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	// reset singleton for test
	dbInstance = nil
	srv := New(testCfg)

	stats := srv.Health()

	if stats["status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["status"])
	}

	if _, ok := stats["error"]; ok {
		t.Fatalf("expected error not to be present")
	}

	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}

func TestClose(t *testing.T) {
	// reset singleton for test
	dbInstance = nil
	srv := New(testCfg)

	if srv.Close() != nil {
		t.Fatalf("expected Close() to return nil")
	}
}
