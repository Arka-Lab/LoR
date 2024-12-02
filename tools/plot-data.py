import numpy as np
import matplotlib.pyplot as plt

# Define the percentages
random_behavior_percentages = np.arange(0, 80, 10)
bad_behavior_percentages = np.arange(0, 80, 10)

# Prepare the data
data = np.random.rand(len(random_behavior_percentages), len(bad_behavior_percentages)) # TODO: replace with actual data

# Create a meshgrid for the 3D plot
X, Y = np.meshgrid(random_behavior_percentages, bad_behavior_percentages)

# Flatten the grid to use with bar3d
x_pos = X.flatten()
y_pos = Y.flatten()
z_pos = np.zeros_like(x_pos)  # All bars start at z=0

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
ax.set_zlabel('Value (Metric)') # TODO: replace with actual metric

# Set ticks for better readability
ax.set_xticks(random_behavior_percentages)
ax.set_yticks(bad_behavior_percentages)

# Set title
ax.set_title('3D Bar Plot for Trader Behavior Data')

# Show the plot
plt.show()