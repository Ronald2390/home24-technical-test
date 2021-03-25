package seeder

import (
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

// Seed struct for seeder
type Seed struct {
	version string
	content string
}

type schemaSeeder struct {
	Version int  `db:"version"`
	Dirty   bool `db:"dirty"`
}

type fileData struct {
	path string
	name string
}

var allSeeder = []Seed{}

// SeedUp is function to Seed Seeder Up
func SeedUp(dbConnectionString string, isDevelopment bool) {
	now := time.Now()

	storage, err := sqlx.Open("postgres", dbConnectionString)
	if err != nil {
		log.Fatalln("failed to open database x: ", err)
	}
	if err = storage.Ping(); err != nil {
		log.Fatalln("failed to connect to database :", err)
	}
	defer storage.Close()

	addAllSeeder()
	addStagingSeeder(isDevelopment)

	sort.SliceStable(allSeeder, func(i, j int) bool {
		return allSeeder[i].version < allSeeder[j].version
	})

	s := findOrCreateSchemaSeeder(storage)

	if s == nil {
		s = &schemaSeeder{}
		initSchema(storage)
	}

	newestVersion := getSchemaVersion(allSeeder[len(allSeeder)-1].version)
	if s.Version == newestVersion {
		log.Println("Seeder is on Same Meta Version")
		return
	}

	var oldVersion = s.Version
	for i, seed := range allSeeder {
		version := getSchemaVersion(seed.version)
		if s.Version == version && s.Dirty == true {
			s.Version = version
			_ = sqlx.MustExec(storage, seed.content)
			s.Dirty = false
		} else if s.Version < version && s.Dirty == false {
			s.Version = version
			_ = sqlx.MustExec(storage, seed.content)
			s.Dirty = false
		}
		if i == len(allSeeder)-1 {
			s.Version = version
		}
	}

	if oldVersion < s.Version {
		updateSchema(storage, oldVersion, s)
	} else if oldVersion > s.Version {
		log.Fatalf("Missing Version, Version on DB %d, latest Version exist %d", oldVersion, s.Version)
	}

	log.Printf("Seeding took %v", time.Since(now))
}

func getSchemaVersion(version string) int {
	var result int
	if n, err := strconv.Atoi(version); err == nil {
		result = n
	} else {
		log.Fatal(version, " is not an integer.")
	}
	return result
}

func initSchema(db *sqlx.DB) {
	query := `
	INSERT INTO schema_seeders ("version", "dirty") VALUES
	(0, false);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func updateSchema(db *sqlx.DB, oldVersion int, schema *schemaSeeder) {
	query := `
		UPDATE schema_seeders SET
		"version" = $1,
		"dirty" = $2
		WHERE "version" = $3;
	`

	_, err := db.Exec(query, schema.Version, schema.Dirty, oldVersion)
	if err != nil {
		log.Fatal(err)
	}
}

func findOrCreateSchemaSeeder(db *sqlx.DB) *schemaSeeder {
	flag := true
	dest := &schemaSeeder{}
	err := db.Get(dest, `SELECT version, dirty FROM schema_seeders LIMIT 1`)
	if err != nil {
		_, isErrPQ := err.(*pq.Error)
		if isErrPQ {
			if err.(*pq.Error).Code == "42P01" {
				// ignore
				qry := `
					CREATE TABLE public.schema_seeders(
						"version" BIGINT NOT NULL UNIQUE,
						"dirty" BOOL NOT NULL);`
				_, err = db.Exec(qry)
				if err != nil {
					log.Fatal(err)
				}
				flag = false
			}
		} else {
			return nil
		}
	}
	if !flag {
		err = db.Get(dest, `SELECT version, dirty FROM schema_seeders LIMIT 1`)
		if err != nil {
			return nil
		}
	}
	return dest
}

// // Append data should on order by timestamps created
func addAllSeeder() {
	// MUST BE ON ORDER BY TIME STAMPS AND FK

	// append all seeder
}

func addStagingSeeder(isDevelopment bool) {
	if isDevelopment {
		allSeeder = append(allSeeder, seed_20210326014100)
	}
}
