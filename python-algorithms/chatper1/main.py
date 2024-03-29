from typing import Dict, List


class Graph:
    def __init__(self):
        self.graph = Dict[int,List[int]]
    def add_edge(self, u, v):
       ## check if u resides in graphlib
       if u in self.graph.keys():
           l = self.graph[u]
           l.append(v)


