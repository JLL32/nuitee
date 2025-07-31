package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	router.HandlerFunc(http.MethodGet, "/v1/hotels", app.listHotelsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:hotelID", app.getHotelHandler)

	// router.HandlerFunc(http.MethodGet, "/v1/hotels/:hotelID/reviews", app.listReviewsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/hotels/:hotelID/reviews/:reviewID", app.getReviewHandler)

	return app.metrics(app.recoverPanic(app.rateLimit(router)))
}
