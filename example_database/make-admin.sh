#! /bin/bash
mongo << EOF
db.users.insert( {'name': "$1"} );
quit()
EOF
