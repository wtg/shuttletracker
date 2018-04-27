#! /bin/bash
mongo << EOF
use shuttle_tracking
db.users.remove( {'username': "$1"} );
quit()
EOF
