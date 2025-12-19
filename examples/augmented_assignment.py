print "=== Augmented Assignment Demo ==="

# Basic arithmetic operations
x = 10
print "x =", x

x += 5
print "x += 5 =>", x

x -= 3
print "x -= 3 =>", x

# String concatenation
message = "Hello"
message += ", World!"
print "String concatenation:", message

# With expressions
count = 1
count += 2 * 3
print "count += 2 * 3 =>", count

# In a loop
total = 0
print "Summing 1 to 5:"
for i in range(1, 6):
    total += i
    print "total +=", i, "=>", total

print "Final total:", total

# With functions
def get_increment():
    return 10

value = 100
value += get_increment()
print "value += get_increment() =>", value