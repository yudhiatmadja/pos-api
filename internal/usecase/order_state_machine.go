package usecase

import "pos-api/internal/domain"

func isValidTransition(current, next domain.OrderStatus) bool {
	switch current {
	case domain.OrderStatusNew:
		return next == domain.OrderStatusAccepted || next == domain.OrderStatusVoided
	case domain.OrderStatusAccepted:
		return next == domain.OrderStatusCooking || next == domain.OrderStatusVoided
	case domain.OrderStatusCooking:
		return next == domain.OrderStatusReady || next == domain.OrderStatusVoided
	case domain.OrderStatusReady:
		return next == domain.OrderStatusDone || next == domain.OrderStatusVoided
	case domain.OrderStatusDone, domain.OrderStatusVoided:
		return false // Terminal states
	default:
		return false
	}
}
