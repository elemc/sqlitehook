package sqlitehook_test

import (
	"database/sql"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/elemc/sqlitehook"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func TestLocalDatabaseAndPrint(t *testing.T) {
	timeout := time.Second * 10
	dsn := "file:test.sqlite3"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		t.Fatalf("Unable to open database: %s", err)
	}
	if db != nil {
		defer func() {
			_ = db.Close()
		}()
	}
	t.Log("Database opened successful")

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	hook, err := sqlitehook.NewSQLiteHook(db, timeout)
	if err != nil {
		t.Fatalf("Unable to initialize hook: %s", err)
	}
	t.Logf("Hook initialized successful")

	log.AddHook(hook)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10000; i++ {
		l := rand.Intn(3)
		logEntry := log.WithField("level_value", l).WithField("database", dsn).WithField("timeout", timeout)
		switch l {
		case 0:
			logEntry.Debug("Debug")
		case 1:
			logEntry.Info("Info")
		case 2:
			logEntry.Warn("Info")
		case 3:
			logEntry.WithError(errors.New("error")).Error("Error")
		}
	}
	t.Log("All works fine")
}
