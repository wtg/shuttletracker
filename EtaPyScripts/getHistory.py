import urllib.request
import datetime as dt
import json

MAX_TIME_DIFFERENCE_MIN = 10

def response(url):
    return urllib.request.urlopen(url)

def loadJSON(response):
    return json.loads(response.read())

def time_in_range(start, end, x):
    """Return true if x is in the range [start, end]"""
    if start <= end:
        return start <= x <= end
    else:
        return start <= x or x <= end

def getAvgVelocity(data, route_id, current_time, weekday):
    totalVelocity = 0
    count = 0
    for i in data:
        for j in i:
            dataArrayTime = j["time"].split(":")
            dataHour = int(dataArrayTime[0].split("T")[1])
            dataMin = int(dataArrayTime[1])
            dataTime = dt.time(dataHour, dataMin, 0)

            dataArrayDay = dataArrayTime[0].split('-')
            dataYear = int(dataArrayDay[0])
            dataMonth = int(dataArrayDay[1])
            dataDay = int(dataArrayDay[2].split('T')[0])

            day = dt.date(dataYear, dataMonth, dataDay)
            dataWeekday = day.weekday()


            start = dt.time(current_time.hour, current_time.minute, current_time.second)
            tmp_startDate = dt.datetime.combine(dt.date(1,1,1), start)

            start = tmp_startDate - dt.timedelta(minutes=MAX_TIME_DIFFERENCE_MIN)
            start = start.time()


            end = tmp_startDate + dt.timedelta(minutes=MAX_TIME_DIFFERENCE_MIN)
            end = end.time()


            if j["route_id"] == route_id and dataWeekday == weekday and time_in_range(start, end, dataTime):
                totalVelocity += j["speed"]
                count += 1
            else:
                continue
    return totalVelocity/count


if __name__ == '__main__':
    # Get what day of the week it is today
    targetWeekday = dt.datetime.today().weekday()
    targetWeekday = 2
    # Get what the current time is now
    targetTime = dt.datetime.now().time()
    targetTime = dt.time(22, 45, 50)


    # Specify which route you want to calculate the average velocity for
    targetRoute = 20


    # Currently on localhost
    url = "http://localhost:8080/history"

    response = urllib.request.urlopen(url)
    data = json.loads(response.read())

    # print(len(data))
    # for i in data:
    #     for j in i:
    #         print(j["time"].split(":"))
        # print(i)

    print(getAvgVelocity(data, 1, targetTime, targetWeekday))

