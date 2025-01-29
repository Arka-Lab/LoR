import os
import sys
import numpy as np
from PIL import Image
from matplotlib import cm
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

def plot_data(data, title, y_label, fit_degree=3, save_as=None):
    # Get the keys and values of the data
    keys = np.sort(data.keys())
    values = np.array([data[key] for key in keys])

    # Fit a polynomial to the data
    z = np.polyfit(keys, values, fit_degree)
    p = np.poly1d(z)
    fitted_data = p(keys)

    # Gradient color based on the alpha
    gradient_values = keys
    gradient_values -= gradient_values.min()
    gradient_values = gradient_values / gradient_values.max()

    # Define the colormap and set the colors
    colors = cm.coolwarm(gradient_values)

    # Create the figure and axis
    fig = plt.figure(figsize=(10, 7))
    ax = fig.add_subplot(111)

    # Plot the data with a gradient color
    for i in range(len(keys) - 1):
        ax.plot(keys[i:i+2], data[i:i+2], '-', color=colors[i])

    # Plot the fitted data
    ax.plot(keys, fitted_data, '--', color='purple')

    # Set axis labels
    ax.set_xlabel(r'$p$ (%)', fontsize=12)
    ax.set_ylabel(y_label)

    # Save the plot if a file name is provided
    if save_as:
        plt.savefig(save_as)
        trim_image(save_as)

    # Set title
    ax.set_title(title)

    # Show the plot
    plt.show()

    # Return the fitted data for further analysis
    return keys, fitted_data

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

                data[int(file_name.split('.')[0])] = result
            except:
                print(f'Error reading file {file_name}')
    return data

if __name__ == '__main__':
    # Check the number of arguments
    if len(sys.argv) != 2:
        print('Usage: python plot-linear-data.py <directory_path>')
        sys.exit(1)

    # Load the data
    file_path = sys.argv[1]
    raw_data = load_data(file_path)

    # Plot percentage of invalid accepted fractal rings
    data = {}
    for file_name in raw_data:
        data[file_name] = raw_data[file_name]['invalid_accept_fractal'] / raw_data[file_name]['fractals'] * 100
    plot_data(data, 'Percentage of Invalid Accepted Fractal Rings', 'Invalid Accepted Fractal Rings (%)', 12, 'images/invalid-accepted-linear.png')

    # Plot percentage of valid rejected fractal rings
    data = {}
    for file_name in raw_data:
        data[file_name] = raw_data[file_name]['valid_reject_fractal'] / raw_data[file_name]['fractals'] * 100
    plot_data(data, 'Percentage of Valid Rejected Fractal Rings', 'Valid Rejected Fractal Rings (%)', 12, 'images/valid-rejected-linear.png')

    # Average adjacency per trader
    data = {}
    for file_name in raw_data:
        data[file_name] = raw_data[file_name]['average_adjacency']
    plot_data(data, 'Average Number of Communication Complexity', 'Number of Communications', 5, 'images/average-communication-linear.png')

    # Maximum adjacency per trader
    data = {}
    for file_name in raw_data:
        data[file_name] = raw_data[file_name]['max_adjacency']
    plot_data(data, 'Maximum Number of Communication Complexity', 'Number of Communications', 5, save_as='images/maximum-communication-linear.png')

    # Average fractal ring acceptance rate per trader
    data = {}
    for file_name in raw_data:
        data[file_name] = raw_data[file_name]['accept_fractal']
    plot_data(data, 'Average Fractal Ring Acceptance Rate', 'Fractal Ring Acceptance Rate (%)', 7, 'images/fractal-acceptance-linear.png')

    # Average satisfaction per coin
    data = {}
    for file_name in raw_data:
        data[file_name] = raw_data[file_name]['coin_satisfaction']
    plot_data(data, 'Average Coin Satisfaction', 'Coin Satisfaction (%)', 12, 'images/coin-satisfaction-linear.png')

    # Average satisfaction per trader
    data = {}
    for file_name in raw_data:
        data[file_name] = raw_data[file_name]['trader_satisfaction']
    plot_data(data, 'Average Trader Satisfaction', 'Trader Satisfaction (%)', 12, 'images/trader-satisfaction-linear.png')