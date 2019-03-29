## Estimated Arrival Time for Shuttles

>Main Formula 
- <img width="464" alt="main" src="https://user-images.githubusercontent.com/42976354/55211380-1afffc80-51c2-11e9-93c1-652b068231d7.PNG">

> Step 1: Calculate V-avg
- https://shuttles.rpi.edu/history Endpoint for the data of all the shuttles in the past 30 days. There are 11 dictionaries, which represent 11 shuttles. Since this endpoint will keep change, though slow, we should read it from the web address. However, since this json file is too large, (more than 700 MB), I decide to download it as local file, and read from local. I set to renew this local file every 30 days. 
-  Figure 1, example for data in History json
- <img width="261" alt="history1" src="https://user-images.githubusercontent.com/42976354/55211501-95c91780-51c2-11e9-9eb6-c98ed92c2309.PNG">

> Step 2: Calculate the real time distance between the shuttle and the nearest stop
- https://shuttles.rpi.edu/updates This endpoint shows real time positions and speeds for the shuttles that are in operation, respectively. 
- Figure 2, example for data in Update json
- <img width="272" alt="update" src="https://user-images.githubusercontent.com/42976354/55211543-c6a94c80-51c2-11e9-871a-21738768b5c9.PNG">
- https://shuttles.rpi.edu/routes This endpoint almost never change. It shows the three routes for school: East Route, West Route, Late Night Route. Those three routes all are curves. To make the calculation easier, since it is hard to calculate the length of curves, I divide each routes to many short line segments. Then add up the length of all line segments, then I will get the very approximate length for those 3 curves. As show in Figure 4.
- Figure 3, example for data in Routes json
- <img width="187" alt="routes" src="https://user-images.githubusercontent.com/42976354/55211567-f9534500-51c2-11e9-88af-ff5bc6e889ad.PNG">
- The algorithm to calculate is showed in the below Figure 4. The dark green point "Raw Shuttle Coordinate" is the real time position for a single shuttle. Then, in https://shuttles.rpi.edu/routes, find the nearest point to the dark green point. The nearest point is shown as the light green point "Projected Point on Route". 
- Then find the nearest stop "yellow point" to the light green point "Projected Point on Route" in the driving direction, by calculating the distance between the light green point with all the stops in the "Stops json" https://shuttles.rpi.edu/stops and find the nearest stop in the driving direction. The reason why I did a lot of  matching/sorting is because the points in the routes are not matched in "Routes json" and "History json" and "Update json". Then calculate the length between the light green point and the nearest stop in the driving direction, by add up the line segments between them.
- Figure 4, algorithm for this calculation.
- <img width="423" alt="raw" src="https://user-images.githubusercontent.com/42976354/55211596-1ab43100-51c3-11e9-84ca-cb4b91e08f01.PNG">
> Step 3: Calculate V
- <img width="232" alt="V" src="https://user-images.githubusercontent.com/42976354/55211167-47674900-51c1-11e9-8306-59ac9fc58d18.PNG">
- n represent number of line segments, V-curr is from https://shuttles.rpi.edu/updates 
> Step 4: finally calculate ETA
- <img width="91" alt="t" src="https://user-images.githubusercontent.com/42976354/55211267-9a410080-51c1-11e9-846a-0cea06a83393.PNG">
- Use go module to import the "t"
