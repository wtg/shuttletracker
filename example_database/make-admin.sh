#! /bin/bash
mongo << EOF
db.users.insert( {'username': "$1"} );
quit()
EOF
