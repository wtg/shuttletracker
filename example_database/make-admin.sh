#! /bin/bash
mongo << EOF
use shuttle_tracking
db.users.insert( {'name': "rolleg"} );
quit()
EOF
