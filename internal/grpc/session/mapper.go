package session

import (
	"uuid"

	"google.golang.org/protobuf/types/known/timestamppb"

	parkingv1 "parking/gen/go/proto/parking"
	"parking/internal/domain"
)

// toProtoSession converts a domain Session to a protobuf Session.
func toProtoSession(session *domain.Session) *parkingv1.ParkingSession {
	pb := &parkingv1.ParkingSession{
		Id:          session.ID.String(),
		ParkingId:   session.ParkingID.String(),
		StateNumber: session.StateNumber,
		EnteredAt:   timestamppb.New(session.EnteredAt),
		Cost:        0,
		Paid:        session.Paid,
	}

	if session.ExitedAt != nil {
		pb.ExitedAt = timestamppb.New(*session.ExitedAt)
	}

	if session.Cost != nil {
		pb.Cost = *session.Cost
	}

	return pb
}

// toDomainSession converts a protobuf Session to a domain Session.
func toDomainSession(pb *parkingv1.ParkingSession) (*domain.Session, error) {
	id, err := uuid.Parse(pb.GetId())
	if err != nil {
		return nil, err
	}

	parkingID, err := uuid.Parse(pb.GetParkingId())
	if err != nil {
		return nil, err
	}

	session := domain.Session{
		ID:          id,
		ParkingID:   parkingID,
		StateNumber: pb.GetStateNumber(),
		Paid:        pb.GetPaid(),
	}

	if pb.GetEnteredAt() != nil {
		session.EnteredAt = pb.GetEnteredAt().AsTime()
	}

	if pb.GetExitedAt() != nil {
		exitedAt := pb.GetExitedAt().AsTime()
		session.ExitedAt = &exitedAt
	}

	cost := pb.GetCost()
	session.Cost = &cost

	return &session, nil
}
