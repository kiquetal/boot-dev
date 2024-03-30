### Install the virtualenvironment

python3 -m venv .venv


### Chapter I
  
- Dijistra's: 

  A graph is an abstract data type that represents a set of vertices and the edges that connect to those vertices.
  
|vertex | relation | v | v | v | v |
| ----- | -------- |---|---|---|---|
|0| connects with | 1 | 4 |   |   |
|1|connects with| 0 | 2 | 3 | 4 |
|2|connects with | 1 | 3 |   |   |
|3|connects with| 1 | 2 | 3 |   |
|4|connects with| 0 | 1 | 3 |   |


- DIJKSTRAS: GET_MIN_DIST_NODE()

GET_MIN_DIST_NODE() FUNCTION
To get the full Dijkstra's algorithm working, we'll need a small helper function called get_min_dist_node(). It looks at a dictionary called distances which is a mapping of
node labels -> distances
And returns the "node label" (which is just a string) with the smallest "distance" value that exists in the unvisited set.
In other words, we're trying to find the node with the smallest distance value that hasn't yet been visited.
ASSIGNMENT
Complete the get_min_dist_node function.


- Full Code
  ASSIGNMENT
  Complete the dijkstra function as described below.

DIJKSTRA(GRAPH, SRC, DEST)
INPUTS
graph: An adjacency list
src: The key of the current node
dest: The key of the destination node
OUTPUTS
The shortest path from src to dest, represented as an array of node keys
ALGORITHM
Create an empty unvisited set to keep track of the nodes that haven't yet been visited
Create an empty predecessors dictionary to keep track of the path we're traversing
Create an empty distances dictionary to keep track of the shortest known distance from the src to each node
Add each node to the unvisited set
Set the distance to the src node to 0
Set the distance to all other nodes to positive infinity
While there are still nodes to visit:
Get the node with the minimum distance from the src that hasn't yet been visited using get_min_dist_node()
Mark that node as visited
If the min_dist_node is the destination:
We're done! Return the full path from src to dest using get_path()
Otherwise:
For each unvisited neighbor of the min_dist_node:
Get the currently known distance to the min_dist_node
Get the distance in the graph between the min_dist_node and the neighbor
Use those 2 values to calculate the total distance from the src to the neighbor
If the total distance to the neighbor is less than the distance we previously had stored for the neighbor in the distances dictionary:
Update the distance we have stored for the neighbor
Update the predecessor of the neighbor to be the min_dist_node
