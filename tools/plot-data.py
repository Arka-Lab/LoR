import os
import sys
import numpy as np
from PIL import Image
from matplotlib import cm
from typing import DefaultDict
import matplotlib.pyplot as plt

def trim_image(file_name):
    # Load image and convert to numpy array
    image = Image.open(file_name)
    image_data = np.asarray(image)

    # Find the non-white pixels
    non_white = np.where(image_data != 255)
    x_min, x_max = non_white[0].min(), non_white[0].max() + 1
    y_min, y_max = non_white[1].min(), non_white[1].max() + 1

    # Crop the image and save it
    image_data = image_data[x_min:x_max, y_min:y_max]
    image = Image.fromarray(image_data)
    image.save(file_name)

def plot_3d_data(data, title, z_label, save_as=None):
    # Create a meshgrid for the 3D plot
    X, Y = np.meshgrid(np.arange(0, 101, 5), np.arange(0, 101, 5))

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
    ax.set_xlabel(r'$\beta$ (%)', fontsize=12)
    ax.set_ylabel(r'$\alpha$ (%)', fontsize=12)
    ax.set_zlabel(z_label)

    # Set ticks for better readability
    ax.set_xticks(np.arange(0, 101, 10))
    ax.set_yticks(np.arange(0, 101, 10))

    # Set initial view
    ax.view_init(elev=30, azim=-20)

    # Save the plot if a file name is provided
    if save_as:
        plt.savefig(save_as)
        trim_image(save_as)

    # Set title
    ax.set_title(title)

    # Show the plot
    plt.show()

def plot_2d_data(data, title, y_label, fit_degree=3, save_as=None):
    # Process the data and find different alpha percentages
    alpha_percentages = np.sort([alpha for alpha in data.keys() if len(data[alpha]) > 0])
    data = [np.mean(data[alpha]) for alpha in alpha_percentages]

    # Fit a polynomial to the data
    z = np.polyfit(alpha_percentages, data, fit_degree)
    p = np.poly1d(z)
    fitted_data = p(alpha_percentages)

    # Gradient color based on the alpha
    gradient_values = np.where(alpha_percentages > 100, 0, alpha_percentages)
    gradient_values -= gradient_values.min()
    gradient_values = gradient_values / gradient_values.max()

    # Define the colormap and set the colors
    colors = cm.coolwarm(gradient_values)

    # Create the figure and axis
    fig = plt.figure(figsize=(10, 7))
    ax = fig.add_subplot(111)

    # Plot the data with a gradient color
    for i in range(len(alpha_percentages) - 1):
        ax.plot(alpha_percentages[i:i+2], data[i:i+2], '-', color=colors[i])

    # Plot the fitted data
    ax.plot(alpha_percentages, fitted_data, '--', color='purple')

    # Set axis labels
    ax.set_xlabel(r'$\gamma$ (%)', fontsize=12)
    ax.set_ylabel(y_label)

    # Save the plot if a file name is provided
    if save_as:
        plt.savefig(save_as)
        trim_image(save_as)

    # Set title
    ax.set_title(title)

    # Show the plot
    plt.show()

def plot_scenarios(data, title, y_label, save_as=None):
    # Create the figure and axis
    fig = plt.figure(figsize=(10, 7))
    ax = fig.add_subplot(111)

    # Plot the data
    ax.plot(np.arange(0, 101, 5), data[0], '-', color='red', label='Bad Behavior')
    ax.plot(np.arange(0, 101, 5), np.array(data).T[0], '-', color='orange', label='Random Behavior')

    # Set axis labels
    ax.set_xlabel('Amount (%)')
    ax.set_ylabel(y_label)
    ax.legend()

    # Save the plot if a file name is provided
    if save_as:
        plt.savefig(save_as)
        trim_image(save_as)
    
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
                    result['run_coins'] = int(lines[2].split(': ')[1])
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
    # Check the number of arguments
    if len(sys.argv) != 2:
        print('Usage: python plot-data.py <directory_path>')
        sys.exit(1)

    # Load the data
    file_path = sys.argv[1]
    raw_data = load_data(file_path)

    # Plot percentage of invalid accepted fractal rings
    data, data_2d = DefaultDict(list), np.zeros((21, 21))
    for i in range(21):
        for j in range(21):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                value = raw_data[file_name]['invalid_accept_fractal'] / raw_data[file_name]['fractals'] * 100
                data_2d[i][j], data[i / 2 + j * 5] = value, np.append(data[i / 2 + j * 5], value)
    plot_3d_data(data_2d, 'Percentage of Invalid Accepted Fractal Rings', 'Invalid Accepted Fractal Rings (%)', save_as='images/invalid-accepted.png')
    plot_2d_data(data, 'Percentage of Invalid Accepted Fractal Rings', 'Invalid Accepted Fractal Rings (%)', 12, 'images/invalid-accepted-2d.png')
    plot_scenarios(data_2d, 'Percentage of Invalid Accepted Fractal Rings', 'Invalid Accepted Fractal Rings (%)', 'images/invalid-accepted-scenario.png')

    # Plot percentage of valid rejected fractal rings
    data, data_2d = DefaultDict(list), np.zeros((21, 21))
    for i in range(21):
        for j in range(21):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                value = raw_data[file_name]['valid_reject_fractal'] / raw_data[file_name]['fractals'] * 100
                data_2d[i][j], data[i / 2 + j * 5] = value, np.append(data[i / 2 + j * 5], value)
    plot_3d_data(data_2d, 'Percentage of Valid Rejected Fractal Rings', 'Valid Rejected Fractal Rings (%)', save_as='images/valid-rejected.png')
    plot_2d_data(data, 'Percentage of Valid Rejected Fractal Rings', 'Valid Rejected Fractal Rings (%)', 12, 'images/valid-rejected-2d.png')
    plot_scenarios(data_2d, 'Percentage of Valid Rejected Fractal Rings', 'Valid Rejected Fractal Rings (%)', 'images/valid-rejected-scenario.png')

    # Average adjacency per trader
    data, data_2d = DefaultDict(list), np.zeros((21, 21))
    for i in range(21):
        for j in range(21):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                value = raw_data[file_name]['average_adjacency']
                data_2d[i][j], data[i / 2 + j * 5] = value, np.append(data[i / 2 + j * 5], value)
    plot_3d_data(data_2d, 'Average Number of Communication Complexity', 'Number of Communications', save_as='images/average-communication.png')
    plot_2d_data(data, 'Average Number of Communication Complexity', 'Number of Communications', 5, 'images/average-communication-2d.png')
    plot_scenarios(data_2d, 'Average Number of Communication Complexity', 'Number of Communications', 'images/average-communication-scenario.png')

    # Maximum adjacency per trader
    data, data_2d = DefaultDict(list), np.zeros((21, 21))
    for i in range(21):
        for j in range(21):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                value = raw_data[file_name]['max_adjacency']
                data_2d[i][j], data[i / 2 + j * 5] = value, np.append(data[i / 2 + j * 5], value)
    plot_3d_data(data_2d, 'Maximum Number of Communication Complexity', 'Number of Communications', save_as='images/maximum-communication.png')
    plot_2d_data(data, 'Maximum Number of Communication Complexity', 'Number of Communications', 5, save_as='images/maximum-communication-2d.png')
    plot_scenarios(data_2d, 'Maximum Number of Communication Complexity', 'Number of Communications', 'images/maximum-communication-scenario.png')

    # Average fractal ring acceptance rate per trader
    data, data_2d = DefaultDict(list), np.zeros((21, 21))
    for i in range(21):
        for j in range(21):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                value = raw_data[file_name]['accept_fractal']
                data_2d[i][j], data[i / 2 + j * 5] = value, np.append(data[i / 2 + j * 5], value)
    plot_3d_data(data_2d, 'Average Fractal Ring Acceptance Rate', 'Fractal Ring Acceptance Rate (%)', save_as='images/fractal-acceptance.png')
    plot_2d_data(data, 'Average Fractal Ring Acceptance Rate', 'Fractal Ring Acceptance Rate (%)', 7, 'images/fractal-acceptance-2d.png')
    plot_scenarios(data_2d, 'Average Fractal Ring Acceptance Rate', 'Fractal Ring Acceptance Rate (%)', 'images/fractal-acceptance-scenario.png')

    # Average satisfaction per coin
    data, data_2d = DefaultDict(list), np.zeros((21, 21))
    for i in range(21):
        for j in range(21):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                value = raw_data[file_name]['coin_satisfaction']
                data_2d[i][j], data[i / 2 + j * 5] = value, np.append(data[i / 2 + j * 5], value)
    plot_3d_data(data_2d, 'Average Coin Satisfaction', 'Coin Satisfaction (%)', save_as='images/coin-satisfaction.png')
    plot_2d_data(data, 'Average Coin Satisfaction', 'Coin Satisfaction (%)', 12, 'images/coin-satisfaction-2d.png')
    plot_scenarios(data_2d, 'Average Coin Satisfaction', 'Coin Satisfaction (%)', 'images/coin-satisfaction-scenario.png')

    # Average satisfaction per trader
    data, data_2d = DefaultDict(list), np.zeros((21, 21))
    for i in range(21):
        for j in range(21):
            file_name = f'{i*5}-{j*5}'
            if file_name in raw_data:
                value = raw_data[file_name]['trader_satisfaction']
                data_2d[i][j], data[i / 2 + j * 5] = value, np.append(data[i / 2 + j * 5], value)
    plot_3d_data(data_2d, 'Average Trader Satisfaction', 'Trader Satisfaction (%)', save_as='images/trader-satisfaction.png')
    plot_2d_data(data, 'Average Trader Satisfaction', 'Trader Satisfaction (%)', 12, 'images/trader-satisfaction-2d.png')
    plot_scenarios(data_2d, 'Average Trader Satisfaction', 'Trader Satisfaction (%)', 'images/trader-satisfaction-scenario.png')