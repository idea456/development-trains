Maybe a greedy method first?
Iterate through each package, for each package:

-   Assign a train that is closest to it first
    -   The train must have enough capacity to fit it
    -   Execute the plan
    -   Once the train picks up the train, can the train have more capacity to pick up more packages?
        -   If yes, find the next nearest package which the train can pick up that fits its capacity
            -   But what if its better to just drop the packages first, since the train could be near the drop site already?
            -   We can have another train pick up that package
            -   2 decisions, both depends which route to them has the least travel time:
                -   MinTravelTime(Pick up another package from X station to Z station, Drop off a package from X station to Y station)
        -   If no:
            -   execute the plan to deliver the packages first to their destinations
            -   find the next package that is nearest to a train that can pick it up
-

Train assignment for different packages

-   pickup, or dropoff?
    -   pickup when we have a nearest package
    -   dropoff if the trains are near a drop site (only possible if the trains are carrying a load)
        -   need to keep track of all packages that have been carried, and check their distance from their drop destinations
-   keep track of packages that have not been picked up in a priority queue
    -   or can track nearest trains
-   keep track of packages that have been pciked up in a priroity queue
    -   2 packages with the same weight, different routes, which one to pickup?
        -   does the train have a destination to go?
            -   If yes, pick the route that is going to the same destination, along the shortest path
            -   If no, pick any route (TODO: optimize this)
    -   Group packages with similar destination?
        -
