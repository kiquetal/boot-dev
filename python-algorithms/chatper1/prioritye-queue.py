class PriorityQueue:
    def __init__(self):
        self.elements = []

    def empty(self):
        return len(self.elements) == 0

    def push(self, priority, item):
        self.elements.append((priority, item))

    def pop(self):
        item = None
        if self.empty():
            return None
        min_index_so_far = 0
        min_priority = float("inf")
        for i, (priority, item) in enumerate(self.elements):
            if priority < min_priority:
                min_priority = priority
                min_index_so_far = i
        if min_index_so_far is not float("inf"):
            item = self.elements[min_index_so_far]
            del self.elements[min_index_so_far]
        return item[1]


class MinHeap:
    def push(self, priority, value):
        index_to_add = -1
        for i, (p, val) in enumerate(self.elements):
            if (priority < p):
                index_to_add = i

        if index_to_add == -1:
            print("non priority")
            self.elements.append((priority, value))
            index_to_add = len(self.elements) - 1
        else:
            self.elements.insert(index_to_add, (priority, value))
        self.bubble_up(index_to_add)

    def bubble_up(self, index):
        if index == 0:
            return None
        parent_index = int((index - 1) / 2)
        if self.elements[parent_index][0] > self.elements[index][0]:
            self.elements[parent_index] = self.elements[index]
            self.bubble_up(parent_index)


    def pop(self):
        print(self.elements)
        if len(self.elements) == 0:
            return None
        el = self.elements[0]
        if len(self.elements) == 1:
            del self.elements[0]
            return el
        last = self.elements[-1]

        self.elements[0] = last
        del self.elements[-1]
        print(self.elements)
        self.bubble_down(0)
        print(self.elements)
        return el


    def bubble_down(self, index):
        left = (2 * index) + 1
        right = 2 * (index + 1)
        smallest = index
        print("actual smallest", smallest)
        if (left < len(self.elements) and self.elements[left][0] < self.elements[smallest][0]):
            smallest = left
        if (right < len(self.elements) and self.elements[right][0] < self.elements[smallest][0]):
            smallest = right
        if (smallest != index):
            print("swapping")
            prev = self.elements[index]
            self.elements[index] = self.elements[smallest]
            self.elements[smallest] = prev
            self.bubble_down(smallest)

    # don't touch below this line

    def __init__(self):
        self.elements = []

    def peek(self):
        if len(self.elements) == 0:
            return None
        return self.elements[0][1]
