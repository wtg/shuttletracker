package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/wtg/shuttletracker"
)

// RouteService implements a mock of shuttletracker.RouteService.
type RouteService struct {
	mock.Mock
}

// CreateRoute creates a Route.
func (rs *RouteService) CreateRoute(route *shuttletracker.Route) error {
	args := rs.Called(route)
	return args.Error(0)
}

// DeleteRoute deletes a Route.
func (rs *RouteService) DeleteRoute(id int) error {
	args := rs.Called(id)
	return args.Error(0)
}

// Route gets a Route.
func (rs *RouteService) Route(id int) (*shuttletracker.Route, error) {
	args := rs.Called(id)
	return args.Get(0).(*shuttletracker.Route), args.Error(1)
}

// ModifyRoute modifies a Route.
func (rs *RouteService) ModifyRoute(route *shuttletracker.Route) error {
	args := rs.Called(route)
	return args.Error(0)
}

// Routes returns all Routes.
func (rs *RouteService) Routes() ([]*shuttletracker.Route, error) {
	args := rs.Called()
	return args.Get(0).([]*shuttletracker.Route), args.Error(1)
}
