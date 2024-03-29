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
