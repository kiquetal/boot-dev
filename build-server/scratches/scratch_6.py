import random


def guessing_number():
    guess = random.randint(1, 100)
    number = int(input("Enter a number: \n"))
    attempts = 1
    while number != guess and attempts < 3:
        if number > guess:
            print("Too high")
        else:
            print("Too low")
        number = int(input("Enter a number:\n"))
        attempts += 1
    if number == guess:
        print("You guessed it")
    else:
        print("You lost")


def mysum(*args):
    sum = 0
    for x in args:
        sum += x
    print(sum)

def hex_output():
    decnum = 0
    hexnum = input("Enter a hexadecimal number: ")
    for power, digit in enumerate(reversed(hexnum)):
        # we use the base16 only to recongnize the hex number.
        decnum += int(digit,16) * (16 ** power)
    print(f"The decimal value of {hexnum} is {decnum}")

def pig_latin(word):
    if word[0] in 'aeiou':
        return f'{word}way'
    return f'{word[1:]}{word[0]}ay'

def pig_latin_sentence(sentence):
    words = sentence.split()
    new_words = []
    for word in words:
        new_words.append(pig_latin(word))
    return " ".join(new_words)

if __name__ == "__main__":
    print(pig_latin_sentence("hello world this is the pig latin translator"))
    # mysum()
    # guessing_number()

