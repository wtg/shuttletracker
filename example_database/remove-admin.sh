#! /bin/bash
mongo << EOF
db.users.remove( {'name': "$1"} );
quit()
EOF
