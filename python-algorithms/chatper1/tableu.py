class SimplexSolver:
    def get_solution_from_tableau(self):
        columns = []
        num_columns = self.rows[0]
        for c in range(len(num_columns)):
            column = [0] * len(self.rows)
            #print(self.rows)
            for r in range(len(self.rows)):
                column[r] = self.rows[r][c]
            columns.append(column)
        results = []
        for c in range(len(columns)-1):
            sum = 0
            val_row = columns[c]
            if val_row.count(1) == 1 and val_row.count(0) == len(val_row) -1:
                results.append(columns[-1][val_row.index(1)])
            else:
                results.append(0)
        return (results,self.objective[-1])

    # don't touch below this line

    def __init__(self, func_coefficients):
        self.objective = []
        for func_coefficient in func_coefficients:
            self.objective.append(func_coefficient)
        self.rows = []
        self.constraints = []

    def add_constraint(self, coefficients, value):
        row = []
        for coefficient in coefficients:
            row.append(coefficient)
        self.rows.append(row)
        self.constraints.append(value)

    def add_slack_variables(self):
        for i in range(len(self.rows)):
            self.objective.append(0)
            basic_cols = [0] * len(self.rows)
            basic_cols[i] = 1
            basic_cols.append(self.constraints[i])
            self.rows[i] += basic_cols
        self.objective.append(0)
