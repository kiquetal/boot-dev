class Graph:
    def __init__(self):
        self.graph = {}

    def add_edge(self, u, v):
        ## check if u resides in graphlib
        # print(f"Lets add {u} and {v} to the graph")
        # print(f"Graph: {self.graph}")
        if u in self.graph.keys():
            # call add only if is not none
            self.graph[u].add(v)
        else:
            if v is not None:
                self.graph[u] = {v}
        if v in self.graph.keys():
            self.graph[v].add(u)
        else:
            if u is not None:
                self.graph[v] = {u}


def get_path(dest, predecessors):
    path = []
    node = dest
    path.insert(0, node)
    while node in predecessors.keys():
        node = predecessors[node]
        path.insert(0, node)
    return path


def get_min_dist_node(distances, unvisited):
    min_dist = 0
    min_node = None
    for n in distances.keys():
        if n in unvisited:
            if min_dist == 0 or distances[n] < min_dist:
                min_dist = distances[n]
                min_node = n
    return min_node
