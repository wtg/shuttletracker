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
- <img width="170" alt="11" src="https://user-images.githubusercontent.com/42976354/55187494-984e5180-516f-11e9-84f3-f94ecb2714a9.PNG">
- https://shuttles.rpi.edu/routes This endpoint almost never change. It shows the three routes for school: East Route, West Route, Late Night Route. Those three routes all are curves. To make the calculation easier, since it is hard to calculate the length of curves, I divide each routes to many short line segments. Then add up the length of all line segments, then I will get the very approximate length for those 3 curves. As show in Figure 4.




> Step 3

> Step 4
