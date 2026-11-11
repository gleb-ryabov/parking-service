package parking

import (
	"context"
	"net/http"
	"time"
	"uuid"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"parking/internal/domain"
	"parking/internal/http/dto"
	handl "parking/internal/http/handler"
	resp "parking/internal/lib/api/response"
)

// Service defines business operations for parking.
type Service interface {
	// CreateParking creates a new parking.
	CreateParking(ctx context.Context, parking domain.Parking) (*domain.Parking, error)
	// UpdateParking updates parking fields.
	UpdateParking(ctx context.Context, parking domain.Parking) error
	// DeleteParking removes parking by ID.
	DeleteParking(ctx context.Context, id uuid.UUID) error
	// GetParkings returns parkings matching the specified filters.
	GetParkings(ctx context.Context, code, address string) ([]domain.Parking, error)
}

// Handler handles HTTP requests for parking operations.
type Handler struct {
	logger  *zerolog.Logger
	timeout time.Duration

	parkingS Service
}

// New creates handler for parking.
func New(
	logger *zerolog.Logger,
	timeout time.Duration,
	parkingS Service,
) *Handler {
	return &Handler{
		logger:   logger,
		timeout:  timeout,
		parkingS: parkingS,
	}
}

// CreateParking handles a request to create a new parking.
func (h *Handler) CreateParking(c *gin.Context) {
	const op = "http.handler.parking.CreateParking"

	log := h.logger.With().
		Str("op", op).
		Str("client_ip", c.ClientIP()).
		Logger()
	log.Debug().Msg("incoming request")

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
	defer cancel()

	var parking dto.ParkingRequest
	if err := c.ShouldBindJSON(&parking); err != nil {
		log.Error().Err(err).Msg("invalid parking parameter")
		handl.RespondError(
			c,
			handl.ErrInvalidParameter,
			http.StatusBadRequest,
			handl.ErrInvalidParameter.Error(),
		)

		return
	}

	res, err := h.parkingS.CreateParking(ctx, parking.ToDomain())
	if err != nil {
		handl.RespondError(
			c,
			err,
			http.StatusInternalServerError,
			"failed to create parking",
		)

		return
	}
	log.Debug().Str("parking_id", res.ID.String()).Msg("parking created")

	resp.Respond(c, http.StatusCreated, resp.OK(res))
}

// UpdateParking handles a request to update a parking.
func (h *Handler) UpdateParking(c *gin.Context) {
	const op = "http.handler.parking.UpdateParking"

	log := h.logger.With().
		Str("op", op).
		Str("client_ip", c.ClientIP()).
		Logger()
	log.Debug().Msg("incoming request")

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
	defer cancel()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handl.RespondError(
			c,
			handl.ErrInvalidUUID,
			http.StatusBadRequest,
			handl.ErrInvalidUUID.Error(),
		)

		return
	}

	var parkingDTO dto.ParkingRequest
	if err = c.ShouldBindJSON(&parkingDTO); err != nil {
		log.Error().Err(err).Msg("invalid parking parameter")
		handl.RespondError(
			c,
			handl.ErrInvalidParameter,
			http.StatusBadRequest,
			handl.ErrInvalidParameter.Error(),
		)

		return
	}

	parking := parkingDTO.ToDomain()
	parking.ID = id

	if err = h.parkingS.UpdateParking(ctx, parking); err != nil {
		handl.RespondError(
			c,
			err,
			http.StatusInternalServerError,
			"failed to update parking",
		)

		return
	}
	log.Debug().Str("parking_id", id.String()).Msg("parking updated")

	resp.Respond(c, http.StatusOK, resp.OK(id.String()))
}

// DeleteParking handles a request to delete a parking.
func (h *Handler) DeleteParking(c *gin.Context) {
	const op = "http.handler.parking.DeleteParking"

	log := h.logger.With().
		Str("op", op).
		Str("client_ip", c.ClientIP()).
		Logger()
	log.Debug().Msg("incoming request")

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
	defer cancel()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handl.RespondError(
			c,
			handl.ErrInvalidUUID,
			http.StatusBadRequest,
			handl.ErrInvalidUUID.Error(),
		)

		return
	}

	if err = h.parkingS.DeleteParking(ctx, id); err != nil {
		handl.RespondError(
			c,
			err,
			http.StatusInternalServerError,
			"failed to delete parking",
		)

		return
	}
	log.Debug().Str("parking_id", id.String()).Msg("parking deleted")

	resp.Respond(c, http.StatusOK, resp.OK(id.String()))
}

// GetParkings handles a request to retrieve parkings.
func (h *Handler) GetParkings(c *gin.Context) {
	const op = "http.handler.parking.GetParkings"

	log := h.logger.With().
		Str("op", op).
		Str("client_ip", c.ClientIP()).
		Logger()

	log.Debug().Msg("incoming request")

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
	defer cancel()

	code := c.Query("code")
	address := c.Query("address")

	parkings, err := h.parkingS.GetParkings(ctx, code, address)
	if err != nil {
		handl.RespondError(
			c,
			err,
			http.StatusInternalServerError,
			"failed to get parkings",
		)

		return
	}

	log.Debug().Int("count", len(parkings)).Msg("parkings retrieved")

	resp.Respond(c, http.StatusOK, resp.OK(parkings))
}
