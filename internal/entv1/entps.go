package entv1

import (
	"database/sql"
	"database/sql/driver"

	"modernc.org/sqlite"
)

func init() {
	sql.Register("sqlite3", sqlite3Driver{Driver: &sqlite.Driver{}})
}

type sqlite3Driver struct {
	*sqlite.Driver
}

type sqlite3DriverConn interface {
	Exec(string, []driver.Value) (driver.Result, error)
}

func (d sqlite3Driver) Open(name string) (conn driver.Conn, err error) {
	conn, err = d.Driver.Open(name)
	if err != nil {
		return
	}
	var sqls = []string{
		"PRAGMA mode=rwc;",
		"PRAGMA foreign_keys=ON;",
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA journal_mode=WAL;",
		"PRAGMA journal_size_limit=104857600;",
		"PRAGMA busy_timeout=5000;",
		"PRAGMA cache=shared;",
		"PRAGMA cache_size=-8000;",
		"PRAGMA mmap_size=134217728;",
		"PRAGMA temp_store=MEMORY;",
		"PRAGMA locking_mode=NORMAL;",
		"PRAGMA cache_spill=ON;",
	}
	for _, stmt := range sqls {
		_, err = conn.(sqlite3DriverConn).Exec(stmt, nil)
		if err != nil {
			_ = conn.Close()
			return
		}
	}
	return
}
