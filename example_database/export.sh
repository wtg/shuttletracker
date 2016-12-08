dbname=shuttle_tracking
currenttime=$(date "+%Y%m%d")
mongoexport --db shuttle_tracking --collection routes --out .tempo_st_A08lvSv0/st_routes.json
mongoexport --db shuttle_tracking --collection stops --out .tempo_st_A08lvSv0/st_stops.json
mongoexport --db shuttle_tracking --collection updates --out .tempo_st_A08lvSv0/st_updates.json
mongoexport --db shuttle_tracking --collection vehicles --out .tempo_st_A08lvSv0/st_vehicles.json
tar -czvf "$dbname${currenttime}.tar.gz" .tempo_st_A08lvSv0
rm -rf .tempo_st_A08lvSv0

