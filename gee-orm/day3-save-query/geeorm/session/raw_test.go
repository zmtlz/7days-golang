package session

import (
	"database/sql"
	"os"
	"testing"

	"geeorm/dialect"
	_ "github.com/mattn/go-sqlite3"
)

var (
	TestDB   *sql.DB
	TestDial dialect.Dialect
)

func setup() {
	TestDB, _ = sql.Open("sqlite3", "gee.db")
	TestDial, _ = dialect.GetDialect("sqlite3")
}

func teardown() {
	_ = TestDB.Close()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func NewSession() *Session {
	return &Session{db: TestDB, dialect: TestDial}
}

func TestSession_Exec(t *testing.T) {
	_, _ = NewSession().Raw("DROP TABLE USER;").Exec()
	_, _ = NewSession().Raw("CREATE TABLE USER(name text);").Exec()
	result, _ := NewSession().Raw("INSERT INTO USER(`name`) values (?), (?)", "Tom", "Sam").Exec()
	if count, err := result.RowsAffected(); err != nil || count != 2 {
		t.Fatal("expect 2, but got", count)
	}
}

func TestSession_QueryRows(t *testing.T) {
	_, _ = NewSession().Raw("DROP TABLE USER;").Exec()
	_, _ = NewSession().Raw("CREATE TABLE USER(name text);").Exec()
	row := NewSession().Raw("SELECT count(*) FROM USER").QueryRow()
	var count int
	if err := row.Scan(&count); err != nil || count != 0 {
		t.Fatal("failed to query db", err)
	}
}
