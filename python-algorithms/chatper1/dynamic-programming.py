def fibonacci(n, precomputed={}):
    if n == 0:
        return 0
    if n == 1:
        return 1
    if n not in precomputed:
        precomputed[n] = fibonacci(n - 1, precomputed) + fibonacci(n - 2, precomputed)
    return precomputed[n]


print(fibonacci(9))

def edit_distance(str1, str2):
    if not str1:
        return len(str2)
    if not str2:
        return len(str1)
    if (str1[-1]==str2[-1]):
        return edit_distance(str1[:-1],str2[:-1])
    insert = edit_distance(str1,str2[:-1])
    remove = edit_distance(str1[:-1],str2)
    substitute = edit_distance(str1[:-1],str2[:-1])
    d=min([insert,remove,substitute])
    return 1+d

def edit_distance_2(str1,str2):
    col = len(str2) + 1
    print(str1,col)
    rows = len(str1) + 1
    print(str2,rows)
    m = [ [ 0 for _ in range(col)] for _ in range(rows)  ]
    for i in range(rows):
        for j in range(col):
            if i == 0:
                m[i][j]= j
            else:
                if j==0:
                    m[i][j]=i
                else:
                    if str1[i-1]==str2[j-1]:
                        m[i][j]=m[i-1][j-1]
                    else:
                        val=[m[i][j-1],m[i-1][j],m[i-1][j-1]]
                        minV = min(val)
                        m[i][j]=1+ minV

    return m[rows-1][col-1]
