# Compliments of @bamnet 
# https://github.com/wtg/shuttle_tracking/blob/master/config/auto_updater.god
RAILS_ROOT = "/home/gbprz/programming/WTG/shuttle_tracking_rails4"

God.pid_file_directory = File.join(RAILS_ROOT, "log")

God.watch do |w|
  w.name = "auto_updater"
  w.interval = 30.seconds
  w.start = "bundle exec rake update_vehicles:auto_update"
  w.dir = RAILS_ROOT
  w.log = File.join(RAILS_ROOT, "log/auto_update.log")

  # clean pid files before start if necessary
  w.behavior(:clean_pid_file)
  
  # determine the state on startup
  w.transition(:init, { true => :up, false => :start }) do |on|
    on.condition(:process_running) do |c|
      c.running = true
    end
  end

  # determine when process has finished starting
  w.transition([:start, :restart], :up) do |on|
    on.condition(:process_running) do |c|
      c.running = true
      c.interval = 5.seconds
    end

    # failsafe
    on.condition(:tries) do |c|
      c.times = 5
      c.transition = :start
      c.interval = 5.seconds
    end
  end

  # start if process is not running
  w.transition(:up, :start) do |on|
    on.condition(:process_running) do |c|
      c.running = false
    end
  end
  
  # retart if memory gets too high
  w.transition(:up, :restart) do |on|
    on.condition(:memory_usage) do |c|
      c.above = 300.megabytes
      c.times = 2
    end
  end

end