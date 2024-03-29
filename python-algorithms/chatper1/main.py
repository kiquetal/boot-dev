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
