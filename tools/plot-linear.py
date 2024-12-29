import os
import sys
import numpy as np
from typing import DefaultDict
import matplotlib.pyplot as plt

def plot_data(data, title, y_label):
    # Process the data and find different alpha percentages
    alpha_percentages = np.sort([alpha for alpha in data.keys()])
    data = [data[alpha] for alpha in alpha_percentages]

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
    data = {}
    for i in range(0, 101, 1):
        file_name = f'{i}'
        if file_name in raw_data:
            data[i] = raw_data[file_name]['run_coins'] / raw_data[file_name]['coins'] * 100
    plot_data(data, 'Percentage of Run Coins', 'Run Coins (%)')

    # Plot percentage of invalid accepted fractal rings
    data = {}
    for i in range(0, 101, 1):
        file_name = f'{i}'
        if file_name in raw_data:
            data[i] = raw_data[file_name]['invalid_accept_fractal'] / raw_data[file_name]['fractals'] * 100
    plot_data(data, 'Percentage of Invalid Accepted Fractal Rings', 'Invalid Accepted Fractal Rings (%)')

    # Plot percentage of valid rejected fractal rings
    data = {}
    for i in range(0, 101, 1):
        file_name = f'{i}'
        if file_name in raw_data:
            data[i] = raw_data[file_name]['valid_reject_fractal'] / raw_data[file_name]['fractals'] * 100
    plot_data(data, 'Percentage of Valid Rejected Fractal Rings', 'Valid Rejected Fractal Rings (%)')

    # Average adjacency per trader
    data = {}
    for i in range(0, 101, 1):
        file_name = f'{i}'
        if file_name in raw_data:
            data[i] = raw_data[file_name]['average_adjacency']
    plot_data(data, 'Average Number of Communication Complexity', 'Number of Communications')

    # Maximum adjacency per trader
    data = {}
    for i in range(0, 101, 1):
        file_name = f'{i}'
        if file_name in raw_data:
            data[i] = raw_data[file_name]['max_adjacency']
    plot_data(data, 'Maximum Number of Communication Complexity', 'Number of Communications')

    # Average fractal ring acceptance rate per trader
    data = {}
    for i in range(0, 101, 1):
        file_name = f'{i}'
        if file_name in raw_data:
            data[i] = raw_data[file_name]['accept_fractal']
    plot_data(data, 'Average Fractal Ring Acceptance Rate', 'Fractal Ring Acceptance Rate (%)')

    # Average satisfaction per coin
    data = {}
    for i in range(0, 101, 1):
        file_name = f'{i}'
        if file_name in raw_data:
            data[i] = raw_data[file_name]['coin_satisfaction']
    plot_data(data, 'Average Coin Satisfaction', 'Coin Satisfaction (%)')

    # Average satisfaction per trader
    data = {}
    for i in range(0, 101, 1):
        file_name = f'{i}'
        if file_name in raw_data:
            data[i] = raw_data[file_name]['trader_satisfaction']
    plot_data(data, 'Average Trader Satisfaction', 'Trader Satisfaction (%)')