package shuttletracker

type ModelService interface {
	VehicleService
	RouteService
	StopService
	LocationService
}
