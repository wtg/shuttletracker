## Estimated Arrival Time for Shuttles

>Main Formula 
<img width="317" alt="11" src="https://user-images.githubusercontent.com/42976354/54885140-768f5a80-4e4f-11e9-811d-b71baef45394.PNG">

> Step 1: Calculate <img width="36" alt="12" src="https://user-images.githubusercontent.com/42976354/54885184-eef61b80-4e4f-11e9-9801-d43986881ee6.PNG">
- https://shuttles.rpi.edu/history Endpoint for the data of all the shuttles in the past 30 days. There are 11 dictionaries, which represent 11 shuttles. Since this endpoint will keep change, though slow, we should read it from the web address. However, since this json file is too large, (more than 700 MB), I decide to download it as local file, and read from local. I set to renew this local file every 30 days. 
-  Example for data in History json
- <img width="151" alt="13" src="https://user-images.githubusercontent.com/42976354/54885341-67111100-4e51-11e9-9554-63dc9fcc13e6.PNG">

> Step 2: Calculate the real time distance between the shuttle and the nearest stop
- https://shuttles.rpi.edu/updates 网页是一个json file， 会变化的，所以需要实时网页读取不能下载json file到本地读本地的文件。周一到周五只有晚上七八点到第二天中午有数据，其他时间是空的。周末晚上九点十点之后到第二天早上五六点才有数据，其他时间是空的。写的时候要考虑到空的情况。在有数据的情况下，长得是一下这个样子，有很多条这样的数据。实时有几辆车在运行就有几条这样的数据。每辆车都要算单独的实时距离。这个网页提供了每辆车的实时位置。


> Step 3

> Step 4
