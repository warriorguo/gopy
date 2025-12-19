# For Loop Examples for GoPy

print "=== Basic For Loop Examples ==="

# Simple range
print "Counting to 5:"
for i in range(5):
    print "Count:", i

# Range with start and end  
print "\nNumbers 2 to 4:"
for i in range(2, 5):
    print i

# Iterating over a list
print "\nList iteration:"
fruits = ["apple", "banana", "cherry"]
for fruit in fruits:
    print "Fruit:", fruit

# Iterating over numbers
print "\nSquaring numbers:"
numbers = [1, 2, 3, 4]
for num in numbers:
    square = num * num
    print num, "squared is", square

# Accumulator pattern
print "\nSum of first 10 numbers:"
total = 0
for i in range(1, 11):
    total = total + i
print "Total:", total

# Nested for loops
print "\nMultiplication table (2x2):"
for i in range(1, 3):
    for j in range(1, 3):
        result = i * j
        print i, "x", j, "=", result

# For loop with conditionals
print "\nEven and odd numbers:"
for i in range(6):
    if i % 2 == 0:
        print i, "is even"
    else:
        print i, "is odd"

# Building a list with for loop
print "\nBuilding a list of squares:"
squares = []
for x in range(1, 6):
    square = x * x
    squares = squares + [square]

print "Squares:", squares

# Using for loop to find maximum
print "\nFinding maximum:"
numbers = [3, 7, 2, 9, 1, 8]
max_val = numbers[0]
for num in numbers:
    if num > max_val:
        max_val = num
print "Maximum value:", max_val

# String iteration
print "\nIterating over string:"
word = "hello"
for char in word:
    print "Character:", char

print "\n=== For Loop Tests Complete ==="