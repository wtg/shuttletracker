#! /bin/bash
mongo << EOF
db.users.remove( {'username': "rolleg"} );
quit()
EOF
