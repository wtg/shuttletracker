package mock

// ModelService implements shuttletracker.ModelService.
type ModelService struct {
	VehicleService
	RouteService
	StopService
	LocationService
}
