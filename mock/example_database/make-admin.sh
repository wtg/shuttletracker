# Set the environment variable, RCS_ID,
# by entering `export RCS_ID=<your_rcs_id>`
# into a terminal, then run the script

#! /bin/bash

echo "Add admin user to the database"

psql shuttletracker << EOF
    INSERT INTO users (id, username) VALUES(0, '$RCS_ID');
EOF
