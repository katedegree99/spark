package handler

import "github.com/katedegree/spark/api/pkg/generated"

// Handler composes all feature handlers into a single StrictServerInterface.
type Handler struct {
	*AuthHandler
	*ImageHandler
	*ProfileHandler
	*ThingHandler
}

var _ generated.StrictServerInterface = (*Handler)(nil)
