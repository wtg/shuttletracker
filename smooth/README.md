# Smooth Tracking

## Summary

Smooth tracking makes predictions of each shuttle's position every second to
give users an idea of where shuttles are in between updates. This document will
discuss the two main components of smooth tracking as well as improvements that
could be made to the algorithm to make predictions more accurate.

## Components

### Algorithm
[smooth_tracking.go](https://github.com/wtg/shuttletracker/blob/smooth-tracking/smooth/smooth_tracking.go)

The algorithm works as follows:
- Calculate the distance the shuttle would have travelled in the elapsed time
  since the last update based on the previous velocity
- Find the closest point on the shuttle's route to the shuttle's previous
  location
- Iterate over each route point, starting from the closest point, and
  accumulate the distance between each of the route points iterated over
- If a sharp angle is encountered, reduce the predicted distance to account for
  the shuttle potentially slowing down.
- When the accumulated distance over the route points exceeds the distance the
  shuttle should have travelled, place the vehicle at route point of the
  current iteration

### Manager
[manager.go](https://github.com/wtg/shuttletracker/blob/smooth-tracking/smooth/manager.go)

The manager is responsible for creating and pushing out shuttle predictions.
Each second, it will make predictions for each shuttle's location, create
"fake" updates from those predictions and display them on the map, and notify
any subscribers of the fake updates. If debug mode is enabled, it will also
output information about the accuracy of the predictions made so far so that
the algorithm can be tuned. The manager is also responsible for reading the
smooth tracking config options.

## Improvements

The algorithm could be improved by accounting for more factors than just the
shuttle's velocity and turns. Most of these improvements would require shuttles
to slow down when they encounter a certain route feature. This could be
accomplished by reducing the predicted distance by a certain factor.
- Shuttle stops: it's reasonable to assume that shuttles movement will pause at
  stops on their route (in progress)
- Intersections/stop lights: shuttles will, most of the time, pause at
  intersections with stop signs or stop lights (in progress)

