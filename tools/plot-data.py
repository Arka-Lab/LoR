import os
import numpy as np
import matplotlib.pyplot as plt

def plot_data(data, title, z_label):
    # Define the percentages
    random_behavior_percentages = np.arange(0, 80, 10)
    bad_behavior_percentages = np.arange(0, 80, 10)

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
    for file in files:
        result = {}
        with open(os.path.join(dir_path, file), 'r') as f:
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
            result['max_adjacency'] = float(lines[10].split(': ')[1])
        data[file.split('.')[0]] = result
    return data

if __name__ == '__main__':
    raw_data = load_data('output/')

    # Plot the percentage of run coins
    data = np.zeros((8, 8))
    for i in range(8):
        for j in range(8):
            if i + j < 8:
                data[i][j] = raw_data[f'{i*10}-{j*10}']['run_coins']
    plot_data(data, 'Percentage of Run Coins', 'Run Coins (%)')

    # Plot percentage of invalid accepted fractal rings
    data = np.zeros((8, 8))
    for i in range(8):
        for j in range(8):
            if i + j < 8:
                data[i][j] = raw_data[f'{i*10}-{j*10}']['invalid_accept_fractal'] / raw_data[f'{i*10}-{j*10}']['fractals']
    plot_data(data, 'Percentage of Invalid Accepted Fractal Rings', 'Invalid Accepted Fractal Rings (%)')

    # Plot percentage of valid rejected fractal rings
    data = np.zeros((8, 8))
    for i in range(8):
        for j in range(8):
            if i + j < 8:
                data[i][j] = raw_data[f'{i*10}-{j*10}']['valid_reject_fractal'] / raw_data[f'{i*10}-{j*10}']['fractals']
    plot_data(data, 'Percentage of Valid Rejected Fractal Rings', 'Valid Rejected Fractal Rings (%)')

    # Average adjacency per trader
    data = np.zeros((8, 8))
    for i in range(8):
        for j in range(8):
            if i + j < 8:
                data[i][j] = raw_data[f'{i*10}-{j*10}']['average_adjacency']
    plot_data(data, 'Average Number of Communication Channels', 'Number of Communication Channels')

    # Maximum adjacency per trader
    data = np.zeros((8, 8))
    for i in range(8):
        for j in range(8):
            if i + j < 8:
                data[i][j] = raw_data[f'{i*10}-{j*10}']['max_adjacency']
    plot_data(data, 'Maximum Number of Communication Channels', 'Number of Communication Channels')

    # Average fractal ring acceptance rate per trader
    data = np.zeros((8, 8))
    for i in range(8):
        for j in range(8):
            if i + j < 8:
                data[i][j] = raw_data[f'{i*10}-{j*10}']['accept_fractal']
    plot_data(data, 'Average Fractal Ring Acceptance Rate', 'Fractal Ring Acceptance Rate (%)')

    # Average satisfaction per coin
    data = np.zeros((8, 8))
    for i in range(8):
        for j in range(8):
            if i + j < 8:
                data[i][j] = raw_data[f'{i*10}-{j*10}']['coin_satisfaction']
    plot_data(data, 'Average Coin Satisfaction', 'Coin Satisfaction (%)')

    # Average satisfaction per trader
    data = np.zeros((8, 8))
    for i in range(8):
        for j in range(8):
            if i + j < 8:
                data[i][j] = raw_data[f'{i*10}-{j*10}']['trader_satisfaction']
    plot_data(data, 'Average Trader Satisfaction', 'Trader Satisfaction (%)')
