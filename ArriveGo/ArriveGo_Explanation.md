## Estimated Arrival Time for Shuttles

>Main Formula 
<img width="317" alt="11" src="https://user-images.githubusercontent.com/42976354/54885140-768f5a80-4e4f-11e9-811d-b71baef45394.PNG">

> Step 1: Calculate <img width="36" alt="12" src="https://user-images.githubusercontent.com/42976354/54885184-eef61b80-4e4f-11e9-9801-d43986881ee6.PNG">
- https://shuttles.rpi.edu/history Endpoint for the data of all the shuttles in the past 30 days. There are 11 dictionaries, which represent 11 shuttles. Since this endpoint will keep change, though slow, we should read it from the web address. However, since this json file is too large, (more than 700 MB), I decide to download it as local file, and read from local. I set to renew this local file every 30 days. 
-  Figure 1, example for data in History json
- <img width="151" alt="13" src="https://user-images.githubusercontent.com/42976354/54885341-67111100-4e51-11e9-9554-63dc9fcc13e6.PNG">

> Step 2: Calculate the real time distance between the shuttle and the nearest stop
- https://shuttles.rpi.edu/updates This endpoint shows real time positions and speeds for the shuttles that are in operation, respectively. 
- Figure 2, example for data in Update json
- <img width="520" alt="update" src="https://user-images.githubusercontent.com/42976354/55192408-c4bb9b00-517a-11e9-894b-991449d11af5.PNG">
- https://shuttles.rpi.edu/routes This endpoint almost never change. It shows the three routes for school: East Route, West Route, Late Night Route. Those three routes all are curves. To make the calculation easier, since it is hard to calculate the length of curves, I divide each routes to many short line segments. Then add up the length of all line segments, then I will get the very approximate length for those 3 curves. As show in Figure 4.
- Figure 3, example for data in Routes json
- <img width="116" alt="22" src="https://user-images.githubusercontent.com/42976354/55191261-14e52e00-5178-11e9-8f34-a82b03a62eec.PNG">
- The algorithm to calculate is showed in the below Figure 4. The dark green point "Raw Shuttle Coordinate" is the real time position for a single shuttle. Then, in https://shuttles.rpi.edu/routes, find the nearest point to the dark green point. The nearest point is shown as the light green point "Projected Point on Route". 
- Then find the nearest stop "yellow point" to the light green point "Projected Point on Route" in the driving direction, by calculating the distance between the light green point with all the stops in the "Stops json" https://shuttles.rpi.edu/stops and find the nearest stop in the driving direction. The reason why I did a lot of  matching/sorting is because the points in the routes are not matched in "Routes json" and "History json" and "Update json". Then calculate the length between the light green point and the nearest stop in the driving direction, by add up the line segments between them.
- Figure 4, algorithm for this calculation.
- <img width="382" alt="33" src="https://user-images.githubusercontent.com/42976354/55192265-6f7f8980-517a-11e9-9c40-ab5ad54dde51.PNG">
> Step 3: Calculate V
- <img width="232" alt="V" src="https://user-images.githubusercontent.com/42976354/55211167-47674900-51c1-11e9-8306-59ac9fc58d18.PNG">
- n represent number of line segments, V-curr is from https://shuttles.rpi.edu/updates 
> Step 4: finally calculate ETA
- <img width="91" alt="t" src="https://user-images.githubusercontent.com/42976354/55211267-9a410080-51c1-11e9-846a-0cea06a83393.PNG">
- Use go module to import the "t"
