#! /bin/bash
mongo << EOF
db.users.remove( {'username': "$1"} );
quit()
EOF
