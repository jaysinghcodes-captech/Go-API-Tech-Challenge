package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/go-chi/httplog/v2"
	_ "github.com/lib/pq"
)

func New(ctx context.Context, connectionString string, logger *httplog.Logger, retryDuration time.Duration) (*sql.DB, error) {
	logger.Info("connecting to database")
	retryCount := 0
	db, err := retryResult(ctx, retryDuration, func() (*sql.DB, error) {
		retryCount++
		return sql.Open("postgres", connectionString)
	})
	if err != nil {
		return nil, fmt.Errorf(
			"[in database.New] Failed to connect to database with retry duration of %s and %d attempts: %w",
			retryDuration,
			retryCount,
			err,
		)
	}
	if db == nil {
		return nil, fmt.Errorf("[in database.New] db is nil")
	}
	logger.Info("Successfully connected to database", "retryCount", retryCount)

	logger.Info("Attempting to ping database")
	retryCount = 0
	err = retry(ctx, retryDuration, func() error {
		retryCount++
		return db.Ping()
	})
	if err != nil {
		if err := db.Close(); err != nil {
			return nil, fmt.Errorf("[in database.New] Failed to ping database and close connection: %w", err)
		}
		return nil, fmt.Errorf(
			"[in database.New] Failed to ping database with retry duration of %s and %d attempts: %w",
			retryDuration,
			retryCount,
			err,
		)
	}
	logger.Info("Successfully pinged database", "retryCount", retryCount)
	logger.Info("database connection established")

	return db, nil
}

func retry(ctx context.Context, maxDuration time.Duration, retryFunc func() error) error {
	_, err := retryResult(ctx, maxDuration, func() (int, error) {
		return 1, retryFunc()
	})

	return err
}

func retryResult[T any](ctx context.Context, maxDuration time.Duration, retryFunc func() (T, error)) (T, error) {
	var (
		returnData T
		err        error
	)
	if maxDuration <= 0 {
		return returnData, errors.New("invalid retry duration supplied")
	}

	const maxBackoffMilliseconds = 5_000.0

	ctx, cancelFunc := context.WithTimeout(ctx, maxDuration)
	defer cancelFunc()

	go func() {
		counter := 1.0
		for {
			counter++
			returnData, err = retryFunc()
			if err != nil {
				waitMilliseconds := math.Min(
					math.Pow(counter, 2)+float64(rand.Intn(10)),
					maxBackoffMilliseconds,
				)
				time.Sleep(time.Duration(waitMilliseconds) * time.Millisecond)
				continue
			}
			cancelFunc()
			return

		}
	}()

	<-ctx.Done()
	return returnData, err 
}