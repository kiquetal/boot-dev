def replace_ending(sentence, old, new):
    # Check if the old string is at the end of the sentence
    arraySentence= sentence.split()
    if (arraySentence[len(arraySentence)-1] == old):
        # Using i as the slicing index, combine the part
        # of the sentence up to the matched string at the
        # end with the new string
        i = sentence.rfind(old)
        new_sentence = sentence[:i] + new +  sentence[i+len(old):]
        return new_sentence

    # Return the original sentence if there is no match
    return sentence

print(replace_ending("It's raining cats and cats", "cats", "dogs"))
# Should display "It's raining cats and dogs"
print(replace_ending("She sells seashells by the seashore", "seashells", "donuts"))
# Should display "She sells seashells by the seashore"
print(replace_ending("The weather is nice in May", "may", "april"))
# Should display "The weather is nice in May"
print(replace_ending("The weather is nice in May", "May", "April"))
# Should display "The weather is nice in April"


filenames = ["program.c", "stdio.hpp", "sample.hpp", "a.out", "math.hpp", "hpp.out"]
# Generate newfilenames as a list containing the new filenames
# using as many lines of code as your chosen method requires.

for index,file in enumerate(filenames):
    if file.endswith(".hpp"):
        filenames[index] = file.replace(".hpp",".h")

print(filenames)


man=["Mike", "Karen", "Jake", "Tasha"]
lt=",".join(man)
print(lt)



def combine_lists(list1, list2):
    # Generate a new list containing the elements of list2
    # Followed by the elements of list1 in reverse order
    list2.reverse()
    list1.extend(list2)
    return list1
Jamies_list = ["Alice", "Cindy", "Bobby", "Jan", "Peter"]


for (index, name) in enumerate(Jamies_list):
    print(index, name)
Drews_list = ["Mike", "Carol", "Greg", "Marcia"]
print(combine_lists(Jamies_list, Drews_list))

colors = ["red", "white", "blue"]
colors.insert(2, "yellow")
print(colors)



Jamies_list = ["Alice", "Cindy", "Bobby", "Jan", "Peter"]

for (index, name) in enumerate(Jamies_list):
    print(index, name)
