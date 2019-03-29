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
- Then
距离算法如下图，深绿色的Raw Shuttle Coordinate就是校车的实时位置。 然后在https://shuttles.rpi.edu/routes 里找出与校车实时位置最近的点。会比较耗时，在loop里算很多次。需要把实时位置（提供经纬度）的和页面上的每个点算距离，然后选出最短的。和实时位置最近的点，就是我们用来算距离的点，如下图的浅绿色点。然后我们就要算浅绿色点和最近的行驶方向的stop（停靠站）的距离。 东线，西线，深夜线，都是环线，都是有固定行驶方向的，不会逆向行驶。行驶的方向如上面的点的顺序在https://shuttles.rpi.edu/routes。把浅绿色点和最近行驶方向的停靠站之间线段的长度加起来，就是这条曲线的大致距离（我们所求的x）, 有11个。还没讲怎么找停靠站的位置，下面会讲stop（停靠站）的位置。
- Figure 4, algorithm for this
- <img width="382" alt="33" src="https://user-images.githubusercontent.com/42976354/55192265-6f7f8980-517a-11e9-9c40-ab5ad54dde51.PNG">




> Step 3

> Step 4
