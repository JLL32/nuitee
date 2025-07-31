package main

import (
	"errors"
	"fmt"

	"net/http"

	"github.com/JLL32/nuitee/internal/data"
	"github.com/JLL32/nuitee/internal/validator"
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
	var input struct {
		Search string
		data.Filters
	}

	hotelID, err := app.readIDParam(r, "hotelID")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Search = app.readString(qs, "search", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "hotel_id", "name", "country", "city", "rating", "starts", "-id", "-hotel_id", "-name", "-country", "-city", "-rating", "-starts"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	reviews, metadata, err := app.models.Reviews.GetAll(hotelID, input.Search, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"meta": metadata, "reviews": reviews}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getReviewSummaryHandler(w http.ResponseWriter, r *http.Request) {
	reviewID, err := app.readIDParam(r, "reviewID")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	hotelID, err := app.readIDParam(r, "hotelID")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	review, err := app.models.Reviews.Get(hotelID, reviewID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	prompt := fmt.Sprintf(`summarize the following hotel review in a few sentences:
		headline: %s
		average score: %d
		pros: %s
		cons: %s`, review.Headline, review.AverageScore, review.Pros, review.Cons)

	result, err := app.OpenAIPrompt(prompt)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if result.Choices[0].Message.Content == "" {
		app.serverErrorResponse(w, r, errors.New("empty summary"))
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"summary": result.Choices[0].Message.Content}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
