# Lem-in

### Objective

Lem-in is an algorithm project, where ants have to go the fastest as possible from start-room to end-room. The graph can involve any number of rooms (nodes) with many different connecting paths. Ants must take turns moving from room to room and only one ant can populate a room at any given time. This program solves this problem for any given scenario.

### Example

Run the program like so:

```sh
$> go run . <filename>.txt
```

An example of the input:

```
4
##start
0 0 3
2 2 5
3 4 0
##end
1 8 3
0-2
2-3
3-1
```

With the format being as follows:

```
Number of ants (4)
Rooms ["name" "x-co-ordinate" "y-co-ordinate"]
Tunnels between rooms ("room"-"room")
```

The ants will travel from the room followed by the ```##start``` command to the room followed by the ```##end``` room. Error handling has been implemented to ensure input files are formatted correctly with informative feedback.

An example of the output:

```console
# file contents

L1-2
L1-3 L2-2
L1-1 L2-3 L3-2
L2-1 L3-3 L4-2
L3-1 L4-3
L4-1
```

In this scenario, the path the ants follow is `0 -> 2 -> 3 -> 1`, under the format `Lx-y` where x is each ant's ID and y is the room they've moved to. In this scenario, it takes 6 'turns' to get all 4 ants from room 0 to room 1.

### Implementation

To solve this problem, the program performs the following steps:

- Take and break down a txt file as input, splitting into corresponding variables
- Find every path possible from start-room to end-room
- Of those paths determine the maximum number of paths that could be used (determined by number of rooms linked to the start/end rooms)
- For 1 path to n maximum number of paths*, determine the optimal paths by starting with the shortest available path and searching for the next shortest path which uses no rooms from the preceding paths
- For each set of paths, run the simulation where ants travel down the paths, recording each turn
- Count the number of turns taken for each set of paths and display the results of the set of paths that took the least number of turns

*With unique paths being the key factor, sometimes more paths isn't the most optimal; it depends on the number of ants travelling the graph.

### Instructions

The following rules of the problem are adhered to by the program:

- A room will never start with the letter `L` or with `#` and must have no spaces.
- A tunnel joins only two rooms together never more than that.
- A room can be linked to multiple rooms.
- Two rooms can't have more than one tunnel connecting them.
- Each room can only contain one ant at a time (except at `##start` and `##end` which can contain as many ants as necessary).
- Each tunnel can only be used once per turn.
- To be the first to arrive, ants will need to take the shortest path or paths. They will also need to avoid traffic jams as well as walking all over their fellow ants.
- Display the ants that moved at each turn, and move each ant only once and through a tunnel (the room at the receiving end must be empty).
- The rooms names will not necessarily be numbers, and in order.
- Any unknown command will be ignored.
- The program handles errors carefully.

### Author
Jez Godfrey - As part of the 01 Founders fellowship
