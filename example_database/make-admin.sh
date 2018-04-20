#! /bin/bash
mongo << EOF
use shuttle_tracking
db.users.insert({'username': "rolleg"});
quit()
EOF
