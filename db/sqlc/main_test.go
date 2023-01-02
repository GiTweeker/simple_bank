package db

import (
	"database/sql"
	"github.com/techschool/simple-bank/util"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {

	var err error
	config, err := util.LoadConfig("../../")

	if err != nil {
		log.Fatal("Cannot load config ", err)
	}

	testDb, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
