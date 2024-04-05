import random


class Title:

    def __init__(self, x, y):
        self.x = x
        self.y = y

    def cost(self):
        random.seed(self.__hash__())
        return random.randint(1, 25)

    def __hash__(self):
        return (self.x * 1000) + self.y

    def __repr__(self):
        return f"({self.x},{self.y})"


