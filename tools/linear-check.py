import matplotlib.pyplot as plt
from matplotlib.ticker import FuncFormatter
from tqdm import tqdm
import numpy as np
from scipy.stats import norm
from numpy.random import default_rng


NUM_REPEATS = 5
RING_SIZE = (50, 200)
NUM_TRADERS = list(range(0, 10 ** 6 + 1, 10 ** 4))[1:]


generator = default_rng().uniform
distribution = norm(loc=0, scale=1)

results = {}
for num_traders in tqdm(NUM_TRADERS):
    total = 0
    for _ in range(NUM_REPEATS):
        consume_prob = distribution.rvs(size=num_traders)
        produce_prob = distribution.rvs(size=num_traders)

        consumers, producers = 0, 0
        for (con_threshould, prod_threshould) in zip(consume_prob, produce_prob):
            con_threshould /= abs(con_threshould) + 1
            prod_threshould /= abs(prod_threshould) + 1
            con_threshould = (con_threshould + 1) / 2
            prod_threshould = (prod_threshould + 1) / 2

            con_prob, prod_prob = generator(0, 1), generator(0, 1)
            if con_prob < con_threshould and prod_prob > prod_threshould:
                consumers += 1
            if prod_prob < prod_threshould and con_prob > con_threshould:
                producers += 1
        
        num_fractals = 0
        num_rings = min(consumers, producers)
        while True:
            ring_size = generator(RING_SIZE[0], RING_SIZE[1])
            if num_rings >= ring_size:
                num_fractals += 1
                num_rings -= ring_size
            else:
                break
        total += num_fractals

    avg_rings = total / NUM_REPEATS
    results[num_traders] = avg_rings

x = list(results.keys())
y = list(results.values())
fitted = np.polyfit(x, y, 1)
y_fit = np.polyval(fitted, x)

plt.plot(x, y, label='Observed Average Fractal Rings')
plt.plot(x, y_fit, label='Fitted Line', linestyle='--')
plt.xlabel('Number of Traders')
plt.ylabel('Average Number of Fractal Rings')
# plt.title('Number of Fractal Rings in one Checkpoint')
plt.legend()

# Format x-axis to show as Ax10^B
def format_func(value, tick_number):
    if value == 0 or value <= 0:
        return '0'
    exponent = int(np.floor(np.log10(value)))
    mantissa = value / (10 ** exponent)
    # Round mantissa to avoid floating point issues
    mantissa = round(mantissa, 2)
    if mantissa == 1.0:
        return r'$10^{{{}}}$'.format(exponent)
    else:
        mantissa_str = str(int(mantissa)) if mantissa == int(mantissa) else str(mantissa)
        return r'${}×10^{{{}}}$'.format(mantissa_str, exponent)

ax = plt.gca()
ax.xaxis.set_major_formatter(FuncFormatter(format_func))

plt.savefig('linear_check.png')