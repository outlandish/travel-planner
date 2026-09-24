package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"travelplanner/internal/ai"
	"travelplanner/internal/app"
	"travelplanner/internal/domain"
	"travelplanner/internal/storage"
)

func parseFlags() (domain.HolidayParams, error) {
	budget := flag.Int("budget", 0, "Budget option")
	holidayType := flag.String("type", "", "Holiday type option")
	holidayNature := flag.String("nature", "", "Holiday nature option")
	flag.Parse()

	holidayParams := domain.HolidayParams{
		Budget:      *budget,
		HolidayType: domain.HolidayType(*holidayType),
		Nature:      domain.HolidayNature(*holidayNature),
	}

	err := holidayParams.Validate()
	if err != nil {
		return domain.HolidayParams{}, err
	}

	return holidayParams, nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	holidayParams, err := parseFlags()
	if err != nil {
		return err
	}

	store, storeErr := storage.NewStore(os.Getenv("SQLITE_DB_PATH"))
	if storeErr != nil {
		return storeErr
	}
	defer store.Close()

	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return errors.New("deepseek api key is not set")
	}

	client := ai.NewClient(apiKey)
	application := app.New(store, client)

	response, responseError := application.GetOrComputeResponse(ctx, holidayParams)
	if responseError != nil {
		return responseError
	}

	fmt.Println(response)

	return nil
}
