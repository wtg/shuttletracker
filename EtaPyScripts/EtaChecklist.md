# ETA Feature Checklist

## Goal 1: History Endpoint
* Code was written last semester
* Waiting on a merge

## Goal 2: Python Implementation
* Implement algorithms in python to make sure it works
### Subgoals
  1. Decide on heuristic(s) to choose historical data to use in average velocity calculations (day of the week, shuttle id, etc)
  2. Write python script to pull historical data from history endpoint (json and parsing), implement with heuristic(s) in Mind
  3. Implement Graph data structure for Stops and routes
  4. Implement distance calculations for stops into Graph
  5. Implement average velocity calculations into Graph (live velocity and historically average velocity calculated by heuristic)
  6. Use above distance and average velocity to assign stop nodes with ETAs

## Goal 3: Test
* Test that python implementation 'works'
* Develop test cases

## Goal 4: Port to Go
* Eventually

## Goal 5: Test
* Test. again.
