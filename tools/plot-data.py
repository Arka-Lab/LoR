import numpy as np
import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D

# Define the percentages
random_behavior_percentages = np.arange(0, 80, 10)
bad_behavior_percentages = np.arange(0, 80, 10)

# Prepare the data (random example data here)
# This should be replaced with your actual data for each (random_behavior, bad_behavior) pair.
# For example, data[i, j] corresponds to the value for random_behavior_percentages[i] and bad_behavior_percentages[j]
data = np.random.rand(len(random_behavior_percentages), len(bad_behavior_percentages))  # Replace this with your actual data
data[0] = -2

# Create a meshgrid for the 3D plot
X, Y = np.meshgrid(random_behavior_percentages, bad_behavior_percentages)

# Flatten the grid to use with bar3d
x_pos = X.flatten()
y_pos = Y.flatten()
z_pos = np.zeros_like(x_pos)  # All bars start at z=0

# TODO: fix this
z_pos = np.ones_like(x_pos) * np.min(data)
data -= np.min(data)

# Bar dimensions
dx = dy = 10  # Each bar will represent 10% intervals
dz = data.flatten()

# Create the figure and 3D axis
fig = plt.figure(figsize=(10, 7))
ax = fig.add_subplot(111, projection='3d')

# Create the 3D bar plot
ax.bar3d(x_pos, y_pos, z_pos, dx, dy, dz, shade=True)

# Set axis labels
ax.set_xlabel('Random Behavior (%)')
ax.set_ylabel('Bad Behavior (%)')
ax.set_zlabel('Value (Metric)')

# Set ticks for better readability
ax.set_xticks(random_behavior_percentages)
ax.set_yticks(bad_behavior_percentages)

# Set title
ax.set_title('3D Bar Plot for Trader Behavior Data')

# Show the plot
plt.show()