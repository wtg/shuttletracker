#! /bin/bash
mongo << EOF
use shuttle_tracking
db.users.insert( {'name': "$1"} );
quit()
EOF