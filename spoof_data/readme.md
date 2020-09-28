# Spoof Data

The update spoofing feature can be enabled in the configuration to create fake
(spoofed) updates at configurable intervals, simulating controlled shuttle
updates rather than getting live updates from the data feed. The spoofed updates
are read from the JSON files in this directory (one file per vehicle). The
updates for each vehicle then will be created sequentially, one for each
vehicle each interval defined in the configuration (default 10s).


## Editing the Data

The starter data in this folder simulates one vehicle on each route running
full, standard laps around each route. The data can be edited as necessary, and
can include an unlimited number of vehicles to be simulated, as long as each
vehicle's updates are confined to its own JSON file and each update defines the
following attributes:
- `latitude`
- `longitude`
- `heading`
- `speed`
- `vehicle_id*`
- `route_id*`
- `tracker_id*`

To use the starter data, it must be edited to define the starred attributes.
Valid values for these attributes will vary based on how your Shuttle Tracker
installation is configured. Values for `vehicle_id` and `route_id` can be found
in the admin panel of your installation. Values for `tracker_id` can be found at
the Shuttle Tracker's [history endpoint](https://shuttles.rpi.edu/history).
Each vehicle must be given a unique tracker ID (in other words, one tracker ID
you find in the historical data per JSON file).


## Using the Update Script

You can use the `update.py` script to automate the addition of these
attributes. Simply run the script, insert the filename and values when
prompted, and the JSON file will be edited to include the proper IDs.

