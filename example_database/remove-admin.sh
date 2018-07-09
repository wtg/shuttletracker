# Set the environment variable, RCS_ID,
# by entering `export RCS_ID="<your_rcs_id>"`
# into a terminal, then run the script

#! /bin/bash

echo "Remove admin user from the database"

psql shuttletracker << EOF
    DELETE FROM users WHERE username='$RCS_ID';
EOF