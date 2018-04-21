package shuttletracker

// ModelService is a collection of interfaces related to vehicles, routes, stops, and their locations.
type ModelService interface {
	VehicleService
	RouteService
	StopService
	LocationService
}
