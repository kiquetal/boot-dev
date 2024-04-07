import matplotlib.pyplot as plt
import numpy as np

x = np.linspace(-10, 10, 400)
y1 = 5*x
y2 = 5*x + 2
y3 = 5*x - 2

plt.figure()
plt.plot(x, y1, label='5x')
plt.plot(x, y2, label='5x + 2')
plt.plot(x, y3, label='5x - 2')
plt.xlabel('x')
plt.ylabel('f(x)')
plt.legend()
plt.grid(True)
plt.show()

