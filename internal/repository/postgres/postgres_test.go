// +build linux
// Run tests only on unix, dockertest compatibility errors with windows platform

package postgres

import (
	"context"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/require"
	"github.com/vrazdalovschi/url-shortener/internal/repository"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/docker"
)

var (
	svc repository.Repository
)

var (
	user     = "postgres"
	password = "secret"
	db       = "postgres"
	port     = "5433"
)

//Test full flow: Create, Read, Increment, Read, Delete.
func TestPostgres(t *testing.T) {
	shortenedId := "WoW21Wow"
	url := "http://google.com"

	//Save record
	err := svc.Save(context.Background(), "apiKey", url, shortenedId, time.Now().AddDate(1, 0, 0).Format("2006-01-02"))
	require.NoError(t, err)

	//Read originalUrl
	res, err := svc.Load(context.Background(), shortenedId)
	require.Equal(t, url, res)
	require.NoError(t, err)

	//Read stats default
	stats, err := svc.Stats(context.Background(), shortenedId)
	require.NoError(t, err)
	expectedStats := &repository.StatsResponse{
		ShortenedId:  shortenedId,
		Redirects:    0,
		LastRedirect: "",
	}
	require.Equal(t, expectedStats, stats)

	//Increment stats
	err = svc.Increment(context.Background(), shortenedId)
	require.NoError(t, err)

	//Check updated stats
	stats, err = svc.Stats(context.Background(), shortenedId)
	require.NoError(t, err)
	parsedTime, err := time.Parse(time.RFC3339, stats.LastRedirect)
	require.NoError(t, err)
	require.WithinDuration(t, time.Now(), parsedTime, time.Second*3)
	require.Equal(t, shortenedId, stats.ShortenedId)
	require.Equal(t, 1, stats.Redirects)

	//Describe shortenedId
	actualDescribe, err := svc.Describe(context.Background(), shortenedId)
	require.NoError(t, err)

	tt := time.Now().AddDate(1, 0, 0)
	expectedExpiryDate := time.Date(tt.Year(), tt.Month(), tt.Day(), 0, 0, 0, 0, time.UTC)

	expectedDescription := &repository.ShortenedIdResponse{
		ApiKey:      "apiKey",
		OriginalURL: url,
		ShortenedId: shortenedId,
		ExpiryDate:  expectedExpiryDate.Format(time.RFC3339),
	}
	require.Equal(t, expectedDescription, actualDescribe)

	//Remove
	err = svc.Delete(context.Background(), shortenedId)
	require.NoError(t, err)

	//Check if removed
	res, err = svc.Load(context.Background(), shortenedId)
	require.Equal(t, "", res)
	require.Error(t, err)

	//Read stats
	stats, err = svc.Stats(context.Background(), shortenedId)
	require.Error(t, err)
	require.Nil(t, stats)
}

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.3",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + db,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err.Error())
	}

	if err = pool.Retry(func() error {
		cfg := repository.Configuration{
			Host:     "localhost",
			Port:     port,
			User:     user,
			Password: password,
			DbName:   db,
		}
		svc, err = NewRepository(cfg)
		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err.Error())
	}

	defer func() {
		log.Println(svc.Close())
	}()

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

/*
func mi() {
	_ = &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			&migrate.Migration{
				Id:   "1",
				Up:   []string{"CREATE TABLE people (id int)", "INSERT INTO "},
				Down: []string{"DROP TABLE people"},
			},
		},
	}

	//migrate.Exec(db, "postgres", migrations, migrate.Up)
}
*/
