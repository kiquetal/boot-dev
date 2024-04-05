import random

class TrafficGrid:
    def __init__(self, width, height):
        self.width = width
        self.height = height


    def in_bounds(self, tile):
        if (tile.x >= 0 and tile.x < self.width ) and (tile.y >= 0 and tile.y < self.height):
            return True
        return False

    def neighbors(self, tile):
        candidate = []

        candidate.append(Tile(tile.x+1,tile.y))
        candidate.append(Tile(tile.x-1,tile.y))
        candidate.append(Tile(tile.x,tile.y-1))
        candidate.append(Tile(tile.x,tile.y+1))

        valid = [can for can in candidate if self.in_bounds(can)]
        return valid
    # don't touch below this line

    def __repr__(self):
        s = ""
        for y in range(self.height - 1, -1, -1):
            for x in range(self.width):
                t = Tile(x, y)
                s += f"[{t.cost():02d}]"
            s += "\n"
        return s


class Tile:
    def __init__(self, x, y):
        self.x = x
        self.y = y

    def cost(self):
        random.seed(hash(self))
        cost = random.randint(1, 25)
        return cost

    def __eq__(self, other):
        if other is None:
            return False
        else:
            return self.x == other.x and self.y == other.y

    def __hash__(self):
        return (self.x * 1000) + self.y

    def __repr__(self):
        return f"({self.x}, {self.y})"
