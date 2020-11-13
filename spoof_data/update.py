import os

while True:
    # Get the name of the next file to update
    print('Input filename ("quit" or newline to quit):')
    filename = input('> ')

    # Check for exit condition
    if len(filename) == 0 or filename == 'quit':
        break

    # Make sure this file exists
    if not os.path.exists(filename):
        print('ERROR: File does not exist')
        continue
    
    # Get vehicle ID to insert into the file
    print('Input vehicle ID:')
    vehicleID = input('> ')

    # Get route ID to insert
    print('Input route ID:')
    routeID = input('> ')

    # Get tracker ID to insert
    print('Input tracker ID:')
    trackerID = input('> ')

    # Extract data from file
    f = open(filename, 'r')
    data = f.readlines()
    f.close()

    # Remove any existing IDs
    i = 0
    while i < len(data):
        line = data[i]
        if 'vehicle_id' in line or 'route_id' in line or 'tracker_id' in line:
            popped = data.pop(i)
            i -= 1
        i += 1

    # Insert new vehicle, route, and tracker IDs
    i = 0
    while i < len(data):
        line = data[i]
        if '{' in line:
            data.insert(i + 1, '    "vehicle_id": {},\n'.format(vehicleID))
            data.insert(i + 2, '    "route_id": {},\n'.format(routeID))
            data.insert(i + 3, '    "tracker_id": "{}",\n'.format(trackerID))
            i += 3
        i += 1

    # Write the new data
    f = open(filename, 'w')
    f.writelines(data)
    f.close()

    print('Successfully inserted IDs into {}'.format(filename))
    
