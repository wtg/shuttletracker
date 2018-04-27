#! /bin/bash
mongo << EOF
use shuttle_tracking
db.users.insert({'username': "$1"});
quit()
EOF
