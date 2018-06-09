package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"gitlab.com/Cacophony/dhelpers"
	"gitlab.com/Cacophony/dhelpers/apihelper"
	"gitlab.com/Cacophony/dhelpers/cache"
)

// New creates a new mux Web Service for reporting information about the SqsProcessor
func New() http.Handler {
	router := chi.NewRouter()

	// setup middleware
	router.Use(middleware.Recoverer)
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: cache.GetLogger(), NoColor: false})
	router.Use(middleware.Logger)
	router.Use(middleware.DefaultCompress)
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.HandleFunc("/stats", getStats)

	return router
}

func getStats(w http.ResponseWriter, r *http.Request) {
	// gather data
	var result apihelper.SqsProcessorStatus
	result.Service = apihelper.GenerateServiceInformation()
	result.Available = true

	// return result
	err := render.Render(w, r, result)
	dhelpers.LogError(err)
}
