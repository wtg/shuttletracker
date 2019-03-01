###Explanation for the Distance Calculation Algorithm###

>Point to Point: How I calculate the length of each route is to divide each route to small line segment, point to point. I calculated the length of each line segments then add then up.  By the calculation, the west route is about 2.34 miles, east is about 4.2, the night route is about 7 miles. I keep 15 digits for decimal parts.
<img width="410" alt="distancego" src="https://user-images.githubusercontent.com/42976354/53600051-9ac08a00-3b76-11e9-86b9-8436737848b7.PNG">

>When we dpoing the real time ETA, we will do the same thing. We will divide the distance between the current location and the next stop into many small line segments, and add them up to find the distance. Then I will use relative velocity to do the division, to get the estimate arival time for a shuttle. This can be as precise as to the second.
><img width="484" alt="disgo" src="https://user-images.githubusercontent.com/42976354/53601270-b6795f80-3b79-11e9-9786-089364849c4c.PNG">




