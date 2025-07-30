package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/hotels", app.listHotelsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:hotelID", app.getHotelHandler)

	router.HandlerFunc(http.MethodGet, "/v1/hotels/:hotelID/reviews", app.listReviewsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:hotelID/reviews/:reviewID", app.getReviewHandler)

	return app.recoverPanic(app.rateLimit(router))
}
