
if 1 != 0:
    print "1 is not equal to 0"

def add(a, b):
    return a + b

def subtract(a, b):
    return a - b

def multiply(a, b):
    return a * b

def divide(a, b):
    a = a+10
    b = b+1
    if b != 0:
        return a / b
    else:
        return "Error: Division by zero"

x = 10
y = 3

print "x =", x
print "y =", y
print "x + y =", add(x, y)
print "x - y =", subtract(x, y)
print "x * y =", multiply(x, y)
print "x / y =", divide(x, y)

numbers = [1, 2, 3, 4, 5]
total = 0
for num in numbers:
    total = total + num

print "Sum of", numbers, "=", total