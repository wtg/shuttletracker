## Relative Velocity Explanations

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


