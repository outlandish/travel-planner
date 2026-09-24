package app

import (
	"context"
	"travelplanner/internal/ai"
	"travelplanner/internal/domain"
	"travelplanner/internal/prompt"
	"travelplanner/internal/storage"
)

type App struct {
	store    *storage.Store
	aiClient *ai.Client
}

func New(store *storage.Store, client *ai.Client) *App {
	return &App{store, client}
}

func (a *App) GetOrComputeResponse(ctx context.Context, params domain.HolidayParams) (string, error) {
	cacheKey := storage.CacheKey(params)
	response, found, err := a.store.GetResponse(ctx, cacheKey)
	if err != nil {
		return "", err
	}

	if found {
		return response, nil
	}

	promptText, err := prompt.Build(params)
	if err != nil {
		return "", err
	}

	aiResponse, clientErr := a.aiClient.Complete(ctx, promptText)
	if clientErr != nil {
		return "", clientErr
	}

	a.store.SaveResponse(ctx, cacheKey, aiResponse)

	return aiResponse, nil
}
