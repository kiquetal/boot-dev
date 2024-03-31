class PriorityQueue:
    def __init__(self):
        self.elements = []

    def empty(self):
        return len(self.elements) == 0

    def push(self, priority, item):
        self.elements.append((priority,item))

    def pop(self):
        item = None
        if self.empty():
            return None
        min_index_so_far = 0
        min_priority = float("inf")
        for i,(priority,item) in enumerate(self.elements):
            if priority < min_priority:
                min_priority = priority
                min_index_so_far = i
        if min_index_so_far is not float("inf"):
            item = self.elements[min_index_so_far]
            del self.elements[min_index_so_far]
        return item[1]



