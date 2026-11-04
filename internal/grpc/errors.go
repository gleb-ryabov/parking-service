package grpc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"parking/internal/service"
)

// MapError maps service errors to gRPC status errors.
func MapError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, service.ErrNotFound):
		return status.Error(codes.NotFound, service.ErrNotFound.Error())
	case errors.Is(err, service.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, service.ErrAlreadyExists.Error())
	case errors.Is(err, service.ErrConstraintViolated):
		return status.Error(codes.FailedPrecondition, service.ErrConstraintViolated.Error())
	case errors.Is(err, service.ErrInvalidData):
		return status.Error(codes.InvalidArgument, service.ErrInvalidData.Error())
	default:
		return status.Error(codes.Internal, service.ErrInternal.Error())
	}
}
