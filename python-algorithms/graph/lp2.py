import matplotlib.pyplot as plt
import numpy as np

# define the constraints
x_cakes = np.linspace(0, 300, 400)  # x_cakes represents the number of cakes
y_cookies1 = 200 - x_cakes*0
y_cookies2 = 300 - x_cakes

# plot the constraints
plt.plot(x_cakes, y_cookies1, label='num_cookies <= 200')
plt.plot(x_cakes, y_cookies2, label='num_cakes + num_cookies <= 300')

# add the vertical line at num_cakes = 250
plt.axvline(x=250, color='r', linestyle='--', label='num_cakes <= 250')

# fill the feasible region

plt.fill_between(x_cakes, 0, np.minimum(y_cookies1, y_cookies2), where=(x_cakes <=250), color='gray', alpha=0.5)


plt.xlim((0, 300))
plt.ylim((0, 300))
plt.xlabel('num_cakes')
plt.ylabel('num_cookies')

# Add a legend
plt.legend(loc='upper right')

plt.show()

