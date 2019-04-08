import json
import pandas
import requests
import time
import math
import numpy as np


# Step1: Get V-average
def getAveV():
    r = requests.get('https://shuttles.rpi.edu/history')
    while r.status_code != 200:
        time.sleep(2)
        r = requests.get('https://shuttles.rpi.edu/history')
    text = r.json() # get data from history endpoint
    car_ave = {}  # Store all vehicles and their corresponding average V

    car = {}   # Store all vehicles and their corresponding history speed data, to calculate for V-ave
    velo_list = []   # list of speed data for a single vehicle

    car_id = []
    car_vel = []
    for item in text:
        # only need car ID and speed
        car_id.append(item['tracker_id'])
        car_vel.append(item['speed'])
    for i in len(text):          # again and again
        if car_id[i] not in car.keys():
            velo_list.append(car_vel[i])
            car[car_id[i]] = velo_list
            velo_list = []
        else:
            velo_list = car[car_id[i]]
            velo_list.append(car_vel[i])
            car[car_id[i]] = velo_list
            velo_list = []
    # For V-ave
    for id, vel in car.items():
        # sum of V
        p = 0
        for v in vel:
            p += int(v)
        ave = int(p)/len(vel)
        car_ave[id] = ave
    return car_ave


# Step 2: distance between 2 points
def distance(self, a, b):
    dist = math.sqrt((a.x-b.x)**2 + (a.y-b.y)**2)
    return dist


# real time position for shuttle  car_real=[{id1:{x:x1,y:y1,route_id:id}}, {id2:{x:x2,y:y2,route_id:id}}]
def getReal():
    r = requests.get('https://shuttles.rpi.edu/updates')
    while r.status_code != 200:
        time.sleep(2)
        r = requests.get('https://shuttles.rpi.edu/updates')
    text = r.json()
    car_real = []
    car_position = {}
    position = {}  # longtitude, latitude
    for item in text:
        position["x"] = item.get("longitude")
        position["y"] = item.get("latitude")
        position["route_id"] = item.get("route_id")
        position["speed"] = item.get("speed")
        car_position[item['tracker_id']] = position
        car_real.append(car_position)
        position = {}
        car_position = {}
    return car_real

# get all info for all routes including routes, stops, longtitude and latitudeid
def getRoutes():
    r = requests.get('https://shuttles.rpi.edu/routes')
    while r.status_code!=200:
        time.sleep(2)
        r = requests.get('https://shuttles.rpi.edu/routes')
    text = r.json()
    all_routes = []
    routes = {}
    info = {}
    for item in text:
        info["stop_ids"] = item.get("stop_ids")
        info["points"] = item.get("points")
        routes[item.get("id")] = info
        all_routes.append(routes)
        info = {}
        routes = {}
    return all_routes


def getStop():
    r = requests.get('https://shuttles.rpi.edu/stops')
    print(type(r.status_code))
    while r.status_code != 200:
        time.sleep(1)
        r = requests.get('https://shuttles.rpi.edu/stops')
    text = r.json()
    all_stops = []
    stops = {}
    info = {}
    for item in text:
        info["longitude"] = item.get("longitude")
        info["latitude"] = item.get("latitude")
        stops[item.get("id")] = info
        all_stops.append(stops)
        info = {}
        stops = {}
    return all_stops

# Still step 2：find the nearest point to the real time position of shuttle
def findClosestPoint():
    # Find the nearest stop in the driving direction
    # find all routes
    real =  getReal()
    print(real)
    print("=====real time route, success======")
    route = getRoutes()
    print(route)
    print("=====all routes, success======")
    stop = getStop()
    print(stop)
    print("=====all stops, succes======")
    # sorting
    all_point_list = []  # [shortest distance between shuttle and all points in its driving routes, route ID]
    for car_item in real:
        print("=====into loop======")
        # current route ID and real time position
        car_id = car_item.keys()
        car_info = car_item.values()
        car_speed = car_info.get("speed")
        car_route_id = car_info.get("route_id")
        car_real_position = {'x':car_info.get("x"),'y':car_info.get("y")}
        # get info for this shuttle
        for l_route in route:
            if int(l_route.keys()) == int(car_route_id):
                route_info = l_route.values()
        # get all stops for this shuttle
        stop_ids = route_info.get("stop_ids")
        route_points = route_info.get("points")   # []
        route_len = len(route_points)
        # loop to find the nearest stop
        point_list1 = []  # use to store all the distances between stops and this shuttle
        print("=====now calculating the light green point 【Projected point on routes】 ======".format(car_id))
        for point in route_points:
            point = {'x': point.get("longitude"), 'y': point.get("latitude")}
            di_point = distance(car_real_position,point)
            point_list1.append(di_point)
        # Step3 sorting on all current operating shuttles, find the nearest point
        point_list1 = np.array(point_list1)
        index_list = np.argsort(point_list1)
        index = index_list[0]  # find the shortest
        print("=====得出车{car_id}该条行驶线路上最近距离点======".format(car_id))
        # 得出所有行驶车辆以及对应线路上的最近点位置
        print(route_points[index])

        # 循环此条线路上所有停车点与此最近点的距离
        point_list2 = [] # 存储预排序的所有距离值(上面程序求出的离车最近的路线点与该条线路的所有停车站点)
        stop_points = []   # 此条线路上一共几个停车点 列表中存储几个点的坐标
        for stop_id in stop_ids:
            for l_stop in stop:
                if int(l_stop.get("id")) == int(stop_id):
                    stop_point = {'x':l_stop.get("longitude"),'y':l_stop.get("latitude")}
                    stop_points.append(stop_point)
        for point in stop_points:
            di_point = distance(route_points[0],point)
            point_list2.append(di_point)
        point_list2 = np.array(point_list2)
        index_list = np.argsort(point_list2)
        index = index_list[0]  # 排序值最小的下标
        print("=====得出车{car_id}该条行驶线路上最近距离点最近的停车站点位{stop_id}======".format(car_id=car_id,stop_id=stop_ids[index]))
        print(stop_ids[index])
        all_point_list.append({'car_id': car_id, 'car_route_id': car_route_id, 'route_point': route_points[index],'stop_point':stop_points[index],'distance':point_list2[index],'car_speed':car_speed,'route_len':route_len})
    return all_point_list


# 第三步:得到实时速度
def getV(V_avg,all_point_list):
    car_number = len(all_point_list)
    V_list = []
    distance_list = []
    for car in all_point_list:
        car_id = car.get("car_id")
        car_route_id = car.get("car_route_id")
        route_point = car.get("route_point")  # 线路上最近的点
        distance = car.get("distance")   # 线路上最近的点与最近的停车点的距离
        car_speed = car.get("car_speed")
        n = car.get("route_len")
        V_avg = V_avg
        v = ((int(n)-1)/V_avg + car_speed)/int(n)
        V_list.append(v)
        distance_list.append(distance)
    return V_list, distance_list

# 第四步
def getPredictTime(V_list, distance_list):
    T_list = []
    for i in range(len(V_list)):
        T = distance_list[i]/V_list[i]
        T_list.append(T)
    return T_list




if __name__ == '__main__':
    # 一
    Vave = getAveV()
    # 二
    all_point_list = findClosestPoint()
    # 三
    V_list, distance_list = getV(Vave,all_point_list)
    # 四
    Final_list = getPredictTime(V_list, distance_list)
    print(Final_list)

