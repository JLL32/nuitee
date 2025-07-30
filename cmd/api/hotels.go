package main

import (
	"errors"
	"net/http"

	"github.com/JLL32/nuitee/internal/data"
)

func (app *application) getHotelHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "hotelID")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	hotel, err := app.models.Hotels.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"hotel": hotel}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listHotelsHandler(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to fetch a list of hotels
	// ...
}
