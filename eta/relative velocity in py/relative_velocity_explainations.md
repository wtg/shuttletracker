## Relative Velocity Explanations Documentation

###Defination and Calculation choice###

>Relative velocity defination: the amount of vehicles minus one, then multiply by the average 
velocity, then add the current velocity, finally the sum divided by the amount of vehicle. And the estimated time are drived from the relative velosity. see below.
<img width="380" alt="eta-v" src="https://user-images.githubusercontent.com/42976354/53597218-9fce0b00-3b6f-11e9-8351-7aa138710b16.PNG">

>Different from reading data from json file the last semester, this time I directly read time data from the 
https://shuttles.rpi.edu/history and 
https://shuttles.rpi.edu/updates, which makes the estimated calculations more precise. When calculating with the average 
velocity based on the historical data, there are two choices for calculating average velocity: 1. Average velocity 
for all the data records in the history endpoint. 2. First find average velocities for each shuttle, then find the average 
velocity of these 11 shuttle. I think the second one is better, since later we may need the individual data for each shuttle.

>For the first step, I will roughly collect just ONE relative velocity to apply for all shuttles for all stops. Later, when we are done with that phase, I will find the real time relative velocity for each shuttle, so we can have more precise result.


###Relative Velocity VS Real Time Velocity###

>The reason why we are not using real time velocity from the end point https://shuttles.rpi.edu/updates is because this endpoint updates every 5 minutes, so the velocity show upon the web page actually is not that "real time". So we implement the historical endpoint https://shuttles.rpi.edu/history that collect the real time velocity for the past 30 days, which shows a consistency of the shuttle operation, since the data for overal performance of a single shuttle for past 30 days is more reliable than how it performs in a single day or single hour. Based on the image below, we can see as the value of "n" increase, the velocity become more accurate.
<img width="353" alt="iii" src="https://user-images.githubusercontent.com/42976354/53599249-bdea3a00-3b74-11e9-8641-e24a32753b75.PNG">
