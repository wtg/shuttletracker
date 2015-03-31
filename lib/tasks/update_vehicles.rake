require 'open-uri'
require 'uri'

namespace update_vehicles: :environment do

  desc "Run task to update vehicle locations every few seconds"
  task :auto_update do 
    # Compliments of @bamnet 
    # https://github.com/wtg/shuttle_tracking/blob/master/lib/tasks/auto_update.rake
    interval = BACKEND_UPDATE
    while true
      start_time = Time.now
      Rake::Task['update_vehicles:locations'].execute
      end_time = Time.now

      #How long did the update take?
      run_time = end_time - start_time
      #If the update took super long (or a bogus answer), ignore the previous run
      run_time = 0 if (run_time < 0 || run_time > interval)
      #Figure out how long to sleep 
      sleep_time = interval - run_time
      puts "Sleeping #{sleep_time}"

      sleep sleep_time
    end
  end

  desc "Update vehicle locations."
  task locations: :environment do 
    # Compliments of @bamnet 
    # https://github.com/wtg/shuttle_tracking/blob/master/lib/tasks/sample.rake
    open(URI.escape(DATA_FEED)) do |f|
      f.each_line do |line|
        puts line
      end
    end
  end

end