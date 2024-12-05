import os
import sys
import numpy as np
from matplotlib import cm
import matplotlib.pyplot as plt

def plot_data(data, title, z_label):
    # Define the percentages
    random_behavior_percentages = np.arange(0, 75, 5)
    bad_behavior_percentages = np.arange(0, 75, 5)

    # Create a meshgrid for the 3D plot
    X, Y = np.meshgrid(random_behavior_percentages, bad_behavior_percentages)

    # Flatten the grid to use with bar3d
    x_pos = X.flatten()                                                                   
    y_pos = Y.flatten()
    z_pos = np.zeros_like(x_pos)  # All bars start at z=0

    # Bar dimensions
    dx = dy = 5  # Each bar will represent 5% intervals
    dz = np.transpose(data).flatten()

    # Define the colormap
    gradient_values = x_pos / 10 + y_pos
    gradient_values -= gradient_values.min()
    gradient_values = gradient_values / gradient_values.max()
    colors = cm.coolwarm(gradient_values)

    # Create the figure and 3D axis
    fig = plt.figure(figsize=(10, 7))
    ax = fig.add_subplot(111, projection='3d')

    # Create the 3D bar plot
    ax.bar3d(x_pos, y_pos, z_pos, dx, dy, dz, color=colors, shade=True)

    # Set axis labels
    ax.set_xlabel('Random Behavior (%)')
    ax.set_ylabel('Bad Behavior (%)')
    ax.set_zlabel(z_label)

    # Set ticks for better readability
    ax.set_xticks(random_behavior_percentages)
    ax.set_yticks(bad_behavior_percentages)

    # Set title
    ax.set_title(title)

    # Show the plot
    plt.show()

def load_data(dir_path):
    data = {}
    files = os.listdir(dir_path)
    for file_name in files:
        if file_name.endswith('.result'):
            result = {}
            try:
                with open(os.path.join(dir_path, file_name), 'r') as f:
                    lines = f.readlines()
                    result['coins'] = int(lines[0].split(': ')[1])
                    result['fractals'] = int(lines[1].split(': ')[1])
                    result['run_coins'] = float(lines[2].split(': ')[1].replace('%', ''))
                    result['submit_fractal'] = float(lines[3].split(': ')[1])
                    result['accept_fractal'] = float(lines[4].split(': ')[1].replace('%', ''))
                    result['invalid_accept_fractal'] = int(lines[5].split(': ')[1])
                    result['valid_reject_fractal'] = int(lines[6].split(': ')[1])
                    result['coin_satisfaction'] = float(lines[7].split(': ')[1].replace('%', ''))
                    result['trader_satisfaction'] = float(lines[8].split(': ')[1].replace('%', ''))
                    result['average_adjacency'] = float(lines[9].split(': ')[1])
                    result['max_adjacency'] = int(lines[10].split(': ')[1])
                data[file_name.split('.')[0]] = result
            except:
                print(f'Error reading file {file_name}')
    return data

if __name__ == '__main__':
    file_path = 'output/'
    if len(sys.argv) > 1:
        file_path = sys.argv[1]
    raw_data = load_data(file_path)

    # Plot the percentage of run coins
    data = np.zeros((15, 15))
    for i in range(15):
        for j in range(15):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                data[i][j] = raw_data[file_name]['run_coins']
    plot_data(data, 'Percentage of Run Coins', 'Run Coins (%)')

    # Plot percentage of invalid accepted fractal rings
    data = np.zeros((15, 15))
    for i in range(15):
        for j in range(15):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                data[i][j] = raw_data[file_name]['invalid_accept_fractal'] / raw_data[file_name]['fractals'] * 100
    plot_data(data, 'Percentage of Invalid Accepted Fractal Rings', 'Invalid Accepted Fractal Rings (%)')

    # Plot percentage of valid rejected fractal rings
    data = np.zeros((15, 15))
    for i in range(15):
        for j in range(15):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                data[i][j] = raw_data[file_name]['valid_reject_fractal'] / raw_data[file_name]['fractals'] * 100
    plot_data(data, 'Percentage of Valid Rejected Fractal Rings', 'Valid Rejected Fractal Rings (%)')

    # Average adjacency per trader
    data = np.zeros((15, 15))
    for i in range(15):
        for j in range(15):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                data[i][j] = raw_data[file_name]['average_adjacency']
    plot_data(data, 'Average Number of Communication Complexity', 'Number of Communications')

    # Maximum adjacency per trader
    data = np.zeros((15, 15))
    for i in range(15):
        for j in range(15):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                data[i][j] = raw_data[file_name]['max_adjacency']
    plot_data(data, 'Maximum Number of Communication Complexity', 'Number of Communications')

    # Average fractal ring acceptance rate per trader
    data = np.zeros((15, 15))
    for i in range(15):
        for j in range(15):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                data[i][j] = raw_data[file_name]['accept_fractal']
    plot_data(data, 'Average Fractal Ring Acceptance Rate', 'Fractal Ring Acceptance Rate (%)')

    # Average satisfaction per coin
    data = np.zeros((15, 15))
    for i in range(15):
        for j in range(15):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                data[i][j] = raw_data[file_name]['coin_satisfaction']
    plot_data(data, 'Average Coin Satisfaction', 'Coin Satisfaction (%)')

    # Average satisfaction per trader
    data = np.zeros((15, 15))
    for i in range(15):
        for j in range(15):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                data[i][j] = raw_data[file_name]['trader_satisfaction']
    plot_data(data, 'Average Trader Satisfaction', 'Trader Satisfaction (%)')
