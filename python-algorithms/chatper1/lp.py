class SimplexSolver:
    def __init__(self, func_coefficients):
        self.objective = [ x for x in func_coefficients]
        self.rows = []
        self.constraints = []


    def add_constraint(self, coefficients, value):
        self.rows.append(coefficients)
        self.constraints.append(value)
