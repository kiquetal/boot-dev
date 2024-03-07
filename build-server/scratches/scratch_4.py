def remove_punctuations(m,n):
    punctuations = '''!()-[]{};:'"\,<>./?@#$%^&*_~'''
    array_word = [x for x in punctuations]
    word = "example!"
    str = ""
    for x in word:
        if x in array_word:
            continue
        str+=x
    print(str)

remove_punctuations("bla","ble")
