tar -xvzf $1
cd .tempo_st_A08lvSv0
mongoimport --collection routes --upsert --file st_routes.json
mongoimport --collection stops --upsert --file st_stops.json
mongoimport --collection updates --upsert --file st_updates.json
mongoimport --collection vehicles --upsert --file st_vehicles.json
cd ..
rm -rf .tempo_st_A08lvSv0
