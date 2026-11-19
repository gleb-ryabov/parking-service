package schedule

import (
	"context"
	"net/http"
	"time"
	"uuid"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"parking/internal/domain"
	"parking/internal/http/dto"
	handl "parking/internal/http/handler"
	"parking/internal/lib/api/parser"
	resp "parking/internal/lib/api/response"
)

// Service defines business operations for schedule.
type Service interface {
	// CreateSchedule creates a new schedule.
	CreateSchedule(ctx context.Context, schedule domain.Schedule) (*domain.Schedule, error)
	// UpdateSchedule updates schedule fields.
	UpdateSchedule(ctx context.Context, schedule domain.Schedule) error
	// DeleteSchedule removes schedule by ID.
	DeleteSchedule(ctx context.Context, id uuid.UUID) error
	// GetSchedulesByParkingID returns schedules for a specific parking with optional active filter.
	GetSchedulesByParkingID(ctx context.Context, parkingID uuid.UUID, active *bool) ([]domain.Schedule, error)
}

// Handler handles HTTP requests for schedule operations.
type Handler struct {
	logger  *zerolog.Logger
	timeout time.Duration

	scheduleS Service
}

// New creates handler for schedule.
func New(
	logger *zerolog.Logger,
	timeout time.Duration,
	scheduleS Service,
) *Handler {
	return &Handler{
		logger:    logger,
		timeout:   timeout,
		scheduleS: scheduleS,
	}
}

// CreateSchedule handles a request to create a new schedule.
// @Summary Create a new schedule
// @Description Creates a new parking schedule with the specified time range, days, and pricing.
// @Tags Schedule
// @Accept json
// @Produce json
// @Param request body dto.ScheduleRequest true "Schedule creation request"
// @Success 201 {object} response.Response{data=domain.Schedule} "Schedule created successfully"
// @Failure 400 {object} response.Response "Invalid request parameters"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/schedule [post]
func (h *Handler) CreateSchedule(c *gin.Context) {
	const op = "http.handler.schedule.CreateSchedule"

	log := h.logger.With().
		Str("op", op).
		Str("client_ip", c.ClientIP()).
		Logger()
	log.Debug().Msg("incoming request")

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
	defer cancel()

	schedule, err := getScheduleFromContext(c)
	if err != nil {
		handl.RespondError(
			c,
			err,
			http.StatusBadRequest,
			"failed to get schedule from request",
		)

		return
	}

	res, err := h.scheduleS.CreateSchedule(ctx, schedule)
	if err != nil {
		handl.RespondError(
			c,
			err,
			http.StatusInternalServerError,
			"failed to create schedule",
		)

		return
	}
	log.Debug().Str("schedule_id", res.ID.String()).Msg("schedule created")

	resp.Respond(c, http.StatusCreated, resp.OK(res))
}

// UpdateSchedule handles a request to update a schedule.
// @Summary Update an existing schedule
// @Description Updates an existing parking schedule with new time range, days, and/or pricing.
// @Tags Schedule
// @Accept json
// @Produce json
// @Param id path string true "Schedule ID (UUID)"
// @Param request body dto.ScheduleRequest true "Schedule update request"
// @Success 200 {object} response.Response{data=string} "Schedule updated successfully, returns schedule ID"
// @Failure 400 {object} response.Response "Invalid schedule ID or request parameters"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/schedule/{id} [put]
func (h *Handler) UpdateSchedule(c *gin.Context) {
	const op = "http.handler.schedule.UpdateSchedule"

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
			err,
			http.StatusBadRequest,
			"failed to get ID from request",
		)

		return
	}

	schedule, err := getScheduleFromContext(c)
	if err != nil {
		handl.RespondError(
			c,
			err,
			http.StatusBadRequest,
			"failed to get schedule from request",
		)
	}

	schedule.ID = id

	if err = h.scheduleS.UpdateSchedule(ctx, schedule); err != nil {
		handl.RespondError(
			c,
			err,
			http.StatusInternalServerError,
			"failed to update schedule",
		)

		return
	}
	log.Debug().Str("schedule_id", id.String()).Msg("schedule updated")

	resp.Respond(c, http.StatusOK, resp.OK(id.String()))
}

// DeleteSchedule handles a request to delete a schedule.
// @Summary Delete a schedule
// @Description Permanently deletes a parking schedule by its ID.
// @Tags Schedule
// @Produce json
// @Param id path string true "Schedule ID (UUID)"
// @Success 200 {object} response.Response{data=string} "Schedule deleted successfully, returns schedule ID"
// @Failure 400 {object} response.Response "Invalid schedule ID"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/schedule/{id} [delete]
func (h *Handler) DeleteSchedule(c *gin.Context) {
	const op = "http.handler.schedule.DeleteSchedule"

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

	if err = h.scheduleS.DeleteSchedule(ctx, id); err != nil {
		handl.RespondError(
			c,
			err,
			http.StatusInternalServerError,
			"failed to delete schedule",
		)

		return
	}
	log.Debug().Str("schedule_id", id.String()).Msg("schedule deleted")

	resp.Respond(c, http.StatusOK, resp.OK(id.String()))
}

// GetSchedules handles a request to retrieve schedules by parking_id.
// @Summary List schedules for a parking
// @Description Retrieves all schedules for a specific parking facility, optionally filtered by active status.
// @Tags Schedule
// @Produce json
// @Param parking_id query string true "Parking ID (UUID)"
// @Param active query boolean false "Filter by active status (true/false)"
// @Success 200 {object} response.Response{data=[]domain.Schedule} "List of schedules"
// @Failure 400 {object} response.Response "Invalid parking_id or active parameter"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/schedule [get]
func (h *Handler) GetSchedules(c *gin.Context) {
	const op = "http.handler.schedule.GetSchedules"

	log := h.logger.With().
		Str("op", op).
		Str("client_ip", c.ClientIP()).
		Logger()

	log.Debug().Msg("incoming request")

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
	defer cancel()

	parkingID, err := parser.ParseRequiredUUID(c.Query("parking_id"))
	if err != nil {
		handl.RespondError(
			c,
			handl.ErrInvalidParameter,
			http.StatusBadRequest,
			"invalid parking_id",
		)

		return
	}

	// Parse active query parameter (optional)
	active, err := parser.ParseBoolPtr(c.Query("active"))
	if err != nil {
		handl.RespondError(
			c,
			handl.ErrInvalidParameter,
			http.StatusBadRequest,
			"active must be true or false",
		)

		return
	}

	schedules, err := h.scheduleS.GetSchedulesByParkingID(ctx, parkingID, active)
	if err != nil {
		handl.RespondError(
			c,
			err,
			http.StatusInternalServerError,
			"failed to get schedules",
		)

		return
	}

	log.Debug().Int("count", len(schedules)).Msg("schedules retrieved")

	resp.Respond(c, http.StatusOK, resp.OK(schedules))
}

func getScheduleFromContext(c *gin.Context) (domain.Schedule, error) {
	var scheduleDTO dto.ScheduleRequest
	if err := c.ShouldBindJSON(&scheduleDTO); err != nil {
		log.Error().Err(err).Msg("invalid schedule parameter")

		return domain.Schedule{}, err
	}

	schedule, err := scheduleDTO.ToDomain()
	if err != nil {
		log.Error().Err(err).Msg("failed to parse schedule")

		return domain.Schedule{}, err
	}

	return schedule, nil
}
