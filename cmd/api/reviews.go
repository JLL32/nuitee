package main

import (
	"errors"
	"net/http"

	"github.com/JLL32/nuitee/internal/data"
)

func (app *application) getReviewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "reviewID")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	hotelID, err := app.readIDParam(r, "hotelID")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}


	review, err := app.models.Reviews.Get(hotelID, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"review": review}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listReviewsHandler(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to fetch a list of reviews
	// ...
}
