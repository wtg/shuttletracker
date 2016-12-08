tar -xvzf sample_database.tar.gz
cd .tempo_st_A08lvSv0
mongoimport --db shuttle_tracking --collection routes --upsert --file st_routes.json
mongoimport --db shuttle_tracking --collection stops --upsert --file st_stops.json
mongoimport --db shuttle_tracking --collection updates --upsert --file st_updates.json
mongoimport --db shuttle_tracking --collection vehicles --upsert --file st_vehicles.json
cd ..
rm -rf .tempo_st_A08lvSv0