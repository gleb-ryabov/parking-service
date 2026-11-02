package parking

import (
	"uuid"

	"google.golang.org/protobuf/types/known/timestamppb"

	parkingv1 "parking/gen/go/proto/parking"
	"parking/internal/domain"
)

// toProtoSession converts a domain ParkingSession to a protobuf ParkingSession.
func toProtoSession(session domain.ParkingSession) *parkingv1.ParkingSession {
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

// toDomainSession converts a protobuf ParkingSession to a domain ParkingSession.
func toDomainSession(pb *parkingv1.ParkingSession) (domain.ParkingSession, error) {
	id, err := uuid.Parse(pb.GetId())
	if err != nil {
		return domain.ParkingSession{}, err
	}

	parkingID, err := uuid.Parse(pb.GetParkingId())
	if err != nil {
		return domain.ParkingSession{}, err
	}

	session := domain.ParkingSession{
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

	return session, nil
}
