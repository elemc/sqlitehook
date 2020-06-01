package sqlitehook

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const defaultTimeout = time.Millisecond * 100

type SQLiteHook struct {
	db      *sql.DB
	timeout time.Duration
}

// NewSQLiteHook - create new SQLite3 logrus hook
func NewSQLiteHook(db *sql.DB, timeout time.Duration) (hook *SQLiteHook, err error) {
	if err = db.Ping(); err != nil {
		return
	}
	hook = &SQLiteHook{
		db:      db,
		timeout: timeout,
	}

	// default timeout
	if hook.timeout == 0 {
		hook.timeout = defaultTimeout
	}
	if err = hook.createTable(); err != nil {
		return
	}
	return
}

func (hook *SQLiteHook) Fire(entry *logrus.Entry) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), hook.timeout)
	defer cancel()

	str, err := entry.String()
	if err != nil {
		err = errors.Wrap(err, "unable to read logrus entry")
		return
	}

	query := `
INSERT INTO logs
(
	timestamp,
	level,
	message,
	full_message
)
VALUES ($1, $2, $3, $4)
`
	if _, err = hook.db.ExecContext(ctx, query,
		entry.Time,
		entry.Level.String(),
		entry.Message,
		str,
	); err != nil {
		err = errors.Wrap(err, "unable to insert log entry")
		return
	}

	return
}

func (hook *SQLiteHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *SQLiteHook) createTable() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), hook.timeout)
	defer cancel()

	query := `create table if not exists logs
(
	id integer
		constraint logs_pk
			primary key autoincrement,
	timestamp datetime not null,
	level varchar(10) not null,
	message TEXT not null,
	full_message TEXT not null
);
create index if not exists logs_level_index
	on logs (level);
`
	if _, err = hook.db.ExecContext(ctx, query); err != nil {
		err = errors.Wrap(err, "unable to initialize SQLite3 logrus hook table")
		return
	}

	return
}
