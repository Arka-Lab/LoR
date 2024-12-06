import os
import sys
import numpy as np
from matplotlib import cm
from typing import DefaultDict
import matplotlib.pyplot as plt

def plot_3d_data(data, title, z_label):
    # Define the percentages
    random_behavior_percentages = np.arange(0, 105, 5)
    bad_behavior_percentages = np.arange(0, 105, 5)

    # Create a meshgrid for the 3D plot
    X, Y = np.meshgrid(random_behavior_percentages, bad_behavior_percentages)

    # Flatten the grid to use with bar3d
    x_pos = X.flatten()                                                                   
    y_pos = Y.flatten()
    z_pos = np.zeros_like(x_pos)  # All bars start at z=0

    # Bar dimensions
    dx = dy = 5  # Each bar will represent 5% intervals
    dz = np.transpose(data).flatten()

    # Gradient color based on the alpha
    gradient_values = np.where(x_pos + y_pos > 100, 0, x_pos / 10 + y_pos)
    gradient_values -= gradient_values.min()
    gradient_values = gradient_values / gradient_values.max()

    # Define the colormap and set the colors
    colors = cm.coolwarm(gradient_values)
    colors[np.where(x_pos + y_pos > 100)] = 0

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

def plot_2d_data(data, title, y_label):
    # Process the data and find different alpha percentages
    alpha_percentages = np.sort([alpha for alpha in data.keys() if len(data[alpha]) > 0])
    data = [np.mean(data[alpha]) for alpha in alpha_percentages]

    # Create the figure and axis
    fig = plt.figure(figsize=(10, 7))
    ax = fig.add_subplot(111)

    # Plot the data
    ax.plot(alpha_percentages, data)

    # Set axis labels
    ax.set_xlabel('Alpha (%)')
    ax.set_ylabel(y_label)

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
                if result['valid_reject_fractal'] > 2:
                    print(file_name, 'has huge amount of valid rejected fractals:', result['valid_reject_fractal'])
                if result['coins'] < 11000:
                    print(file_name, 'has small amount of coins:', result['coins'])
            except:
                print(f'Error reading file {file_name}')
    return data

if __name__ == '__main__':
    # Check the number of arguments
    if len(sys.argv) != 2:
        print('Usage: python plot-data.py <directory_path>')
        sys.exit(1)

    # Load the data
    file_path = sys.argv[1]
    raw_data = load_data(file_path)

    # Plot the percentage of run coins
    data, data_2d = DefaultDict(list), np.zeros((21, 21))
    for i in range(21):
        for j in range(21):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                data_2d[i][j] = raw_data[file_name]['run_coins']
                data[i / 2 + j * 5] = np.append(data[i + 2 * j], raw_data[file_name]['run_coins'])
    plot_3d_data(data_2d, 'Percentage of Run Coins', 'Run Coins (%)')
    plot_2d_data(data, 'Percentage of Run Coins', 'Run Coins (%)')

    # Plot percentage of invalid accepted fractal rings
    data, data_2d = DefaultDict(list), np.zeros((21, 21))
    for i in range(21):
        for j in range(21):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                data_2d[i][j] = raw_data[file_name]['invalid_accept_fractal'] / raw_data[file_name]['fractals'] * 100
                data[i / 2 + j * 5] = np.append(data[i + 2 * j], raw_data[file_name]['invalid_accept_fractal'] / raw_data[file_name]['fractals'] * 100)
    plot_3d_data(data_2d, 'Percentage of Invalid Accepted Fractal Rings', 'Invalid Accepted Fractal Rings (%)')
    plot_2d_data(data, 'Percentage of Invalid Accepted Fractal Rings', 'Invalid Accepted Fractal Rings (%)')

    # Plot percentage of valid rejected fractal rings
    data, data_2d = DefaultDict(list), np.zeros((21, 21))
    for i in range(21):
        for j in range(21):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                data_2d[i][j] = raw_data[file_name]['valid_reject_fractal'] / raw_data[file_name]['fractals'] * 100
                data[i / 2 + j * 5] = np.append(data[i + 2 * j], raw_data[file_name]['valid_reject_fractal'] / raw_data[file_name]['fractals'] * 100)
    plot_3d_data(data_2d, 'Percentage of Valid Rejected Fractal Rings', 'Valid Rejected Fractal Rings (%)')
    plot_2d_data(data, 'Percentage of Valid Rejected Fractal Rings', 'Valid Rejected Fractal Rings (%)')

    # Average adjacency per trader
    data, data_2d = DefaultDict(list), np.zeros((21, 21))
    for i in range(21):
        for j in range(21):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                data_2d[i][j] = raw_data[file_name]['average_adjacency']
                data[i / 2 + j * 5] = np.append(data[i + 2 * j], raw_data[file_name]['average_adjacency'])
    plot_3d_data(data_2d, 'Average Number of Communication Complexity', 'Number of Communications')
    plot_2d_data(data, 'Average Number of Communication Complexity', 'Number of Communications')

    # Maximum adjacency per trader
    data, data_2d = DefaultDict(list), np.zeros((21, 21))
    for i in range(21):
        for j in range(21):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                data_2d[i][j] = raw_data[file_name]['max_adjacency']
                data[i / 2 + j * 5] = np.append(data[i + 2 * j], raw_data[file_name]['max_adjacency'])
    plot_3d_data(data_2d, 'Maximum Number of Communication Complexity', 'Number of Communications')
    plot_2d_data(data, 'Maximum Number of Communication Complexity', 'Number of Communications')

    # Average fractal ring acceptance rate per trader
    data, data_2d = DefaultDict(list), np.zeros((21, 21))
    for i in range(21):
        for j in range(21):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                data_2d[i][j] = raw_data[file_name]['accept_fractal']
                data[i / 2 + j * 5] = np.append(data[i + 2 * j], raw_data[file_name]['accept_fractal'])
    plot_3d_data(data_2d, 'Average Fractal Ring Acceptance Rate', 'Fractal Ring Acceptance Rate (%)')
    plot_2d_data(data, 'Average Fractal Ring Acceptance Rate', 'Fractal Ring Acceptance Rate (%)')

    # Average satisfaction per coin
    data, data_2d = DefaultDict(list), np.zeros((21, 21))
    for i in range(21):
        for j in range(21):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                data_2d[i][j] = raw_data[file_name]['coin_satisfaction']
                data[i / 2 + j * 5] = np.append(data[i + 2 * j], raw_data[file_name]['coin_satisfaction'])
    plot_3d_data(data_2d, 'Average Coin Satisfaction', 'Coin Satisfaction (%)')
    plot_2d_data(data, 'Average Coin Satisfaction', 'Coin Satisfaction (%)')

    # Average satisfaction per trader
    data, data_2d = DefaultDict(list), np.zeros((21, 21))
    for i in range(21):
        for j in range(21):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                data_2d[i][j] = raw_data[file_name]['trader_satisfaction']
                data[i / 2 + j * 5] = np.append(data[i + 2 * j], raw_data[file_name]['trader_satisfaction'])
    plot_3d_data(data_2d, 'Average Trader Satisfaction', 'Trader Satisfaction (%)')
    plot_2d_data(data, 'Average Trader Satisfaction', 'Trader Satisfaction (%)')