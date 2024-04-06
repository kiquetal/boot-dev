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
