print "Hello, World!"

x = 42
y = 3.14
name = "GoPy"

print "x =", x
print "y =", y
print "name =", name

for i in range(5):
    print "i =", i

if x > 40:
    print "x is greater than 40"
    
numbers = [1, 2, 3, 4, 5]
print "numbers =", numbers
print "length =", len(numbers)

def greet(name):
    print "Hello,", name
    return "Greeted " + name

result = greet("Python")
print "Result:", result