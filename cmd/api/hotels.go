package main

import (
	"errors"
	"net/http"

	"github.com/JLL32/nuitee/internal/data"
	"github.com/JLL32/nuitee/internal/validator"
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
	var input struct {
		Search string
		Filters     data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Search = app.readString(qs, "search", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "hotel_id")
	input.Filters.SortSafelist = []string{"hotel_id", "name", "country", "city", "rating", "starts", "-hotel_id", "-name", "-country", "-city", "-rating", "-starts"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	hotels, metadata, err := app.models.Hotels.GetAll(input.Search, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "hotels": hotels}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
