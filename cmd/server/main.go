package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"travelplanner/internal/ai"
	"travelplanner/internal/app"
	"travelplanner/internal/domain"
	"travelplanner/internal/storage"
)

type Response struct {
	Response string `json:"response"`
	Status   string `json:"status"`
}

func writeJSON(w http.ResponseWriter, status int, resp Response) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	store, storeErr := storage.NewStore(os.Getenv("SQLITE_DB_PATH"))
	if storeErr != nil {
		fmt.Fprintln(os.Stderr, storeErr)
		os.Exit(1)
	}
	defer store.Close()

	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, errors.New("deepseek api key is not set"))
		os.Exit(1)
	}

	client := ai.NewClient(apiKey)
	application := app.New(store, client)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /plan", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		budgetParam := request.URL.Query().Get("budget")
		if budgetParam == "" {
			writeJSON(writer, http.StatusBadRequest, Response{Response: "budget is empty", Status: "error"})
			return
		}

		budget, budgetConvertError := strconv.Atoi(budgetParam)
		if budgetConvertError != nil {
			writeJSON(writer, http.StatusBadRequest, Response{Response: "budget must be a number", Status: "error"})
			return
		}

		holidayType := request.URL.Query().Get("type")
		holidayNature := request.URL.Query().Get("nature")

		holidayParams := domain.HolidayParams{
			Budget:      budget,
			HolidayType: domain.HolidayType(holidayType),
			Nature:      domain.HolidayNature(holidayNature),
		}
		holidayParamsError := holidayParams.Validate()
		if holidayParamsError != nil {
			writeJSON(writer, http.StatusBadRequest, Response{Response: holidayParamsError.Error(), Status: "error"})
			return
		}

		response, aiResponseError := application.GetOrComputeResponse(request.Context(), holidayParams)
		if aiResponseError != nil {
			writeJSON(writer, http.StatusInternalServerError, Response{Response: aiResponseError.Error(), Status: "error"})
			return
		}

		writeJSON(writer, http.StatusOK, Response{Response: response, Status: "success"})
	})

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		fmt.Fprintln(os.Stderr, errors.New("server port env variable is not set"))
		os.Exit(1)
	}

	httpServerError := http.ListenAndServe("localhost:"+serverPort, mux)
	if httpServerError != nil {
		fmt.Fprintln(os.Stderr, httpServerError)
		os.Exit(1)
	}

}
