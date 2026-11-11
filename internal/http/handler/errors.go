package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	resp "parking/internal/lib/api/response"
	svc "parking/internal/service"
)

var (
	// ErrInvalidParameter is returned when a parameter is invalid.
	ErrInvalidParameter = errors.New("invalid parameter")
	// ErrInvalidUUID is returned when UUID is invalid.
	ErrInvalidUUID = errors.New("invalid uuid")
	// ErrInternalServerError is returned when an internal server error occurs.
	ErrInternalServerError = errors.New("internal server error")
)

// RespondError maps service and handler errors to HTTP responses.
func RespondError(c *gin.Context, err error, defaultCode int, defaultMessage string) {
	switch {
	case errors.Is(err, svc.ErrNotFound):
		resp.Respond(c, http.StatusNotFound, resp.Error(svc.ErrNotFound))
	case errors.Is(err, svc.ErrAlreadyExists):
		resp.Respond(c, http.StatusConflict, resp.Error(svc.ErrAlreadyExists))
	case errors.Is(err, svc.ErrInvalidData):
		resp.Respond(c, http.StatusBadRequest, resp.Error(svc.ErrInvalidData))
	case errors.Is(err, svc.ErrConstraintViolated):
		resp.Respond(c, http.StatusConflict, resp.Error(svc.ErrConstraintViolated))
	case errors.Is(err, svc.ErrInternal):
		resp.Respond(c, http.StatusInternalServerError, resp.Error(svc.ErrInternal))
	case errors.Is(err, ErrInvalidParameter):
		resp.Respond(c, http.StatusBadRequest, resp.Error(ErrInvalidParameter))
	case errors.Is(err, ErrInvalidUUID):
		resp.Respond(c, http.StatusBadRequest, resp.Error(ErrInvalidUUID))
	default:
		resp.Respond(c, defaultCode, resp.Error(fmt.Errorf("%s", defaultMessage)))
	}
}
