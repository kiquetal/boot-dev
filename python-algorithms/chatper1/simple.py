class SimplexSolver:
    def solve(self):
        self.add_slack_variables()
        while self.should_pivot():
            pivot_colum = self.get_pivot_col()
            pivot_row = self.get_pivot_row(pivot_colum)
            self.pivot(pivot_row,pivot_colum)


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

    def get_pivot_col(self):
        low = 0
        pivot_idx = 0
        for i in range(len(self.objective) - 1):
            if self.objective[i] < low:
                low = self.objective[i]
                pivot_idx = i
        return pivot_idx

    def get_pivot_row(self, col_idx):
        last_col = [self.rows[i][-1] for i in range(len(self.rows))]
        pivot_col = [self.rows[i][col_idx] for i in range(len(self.rows))]
        min_ratio = float("inf")
        min_ratio_idx = -1
        for i in range(len(last_col)):
            ratio = float("inf")
            if pivot_col[i] == 0:
                ratio = 99999999
            else:
                ratio = last_col[i] / pivot_col[i]
            if ratio < 0:
                continue
            if ratio < min_ratio:
                min_ratio = ratio
                min_ratio_idx = i
        if min_ratio_idx == -1:
            raise Exception("no non-negative ratios, problem doesn't have a solution")
        return min_ratio_idx

    def pivot(self, pivot_row_idx, pivot_col_idx):
        pivot_val = self.rows[pivot_row_idx][pivot_col_idx]
        for i in range(len(self.rows[pivot_row_idx])):
            self.rows[pivot_row_idx][i] = self.rows[pivot_row_idx][i] / pivot_val
        for i in range(len(self.rows)):
            if i == pivot_row_idx:
                continue
            mul = self.rows[i][pivot_col_idx]
            for j in range(len(self.rows[i])):
                self.rows[i][j] = self.rows[i][j] - mul * self.rows[pivot_row_idx][j]
        mul = self.objective[pivot_col_idx]
        for i in range(len(self.objective)):
            self.objective[i] = self.objective[i] - mul * self.rows[pivot_row_idx][i]

    def should_pivot(self):
        return min(self.objective[:-1]) < 0

    def add_slack_variables(self):
        for i in range(len(self.rows)):
            self.objective.append(0)
            basic_cols = [0] * len(self.rows)
            basic_cols[i] = 1
            basic_cols.append(self.constraints[i])
            self.rows[i] += basic_cols
        self.objective.append(0)

    def get_solution_from_tableau(self):
        cols = []
        for colI in range(len(self.rows[0])):
            col = [0] * len(self.rows)
            for rowI in range(len(self.rows)):
                col[rowI] = self.rows[rowI][colI]
            cols.append(col)

        results = []
        for i in range(len(cols) - 1):
            if cols[i].count(0) == len(cols[i]) - 1 and 1 in cols[i]:
                results.append(cols[-1][cols[i].index(1)])
            else:
                results.append(0)
        return results, self.objective[-1]
