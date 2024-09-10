import json
import numpy as np
from datetime import datetime, timedelta

def load_text(filename):
    """Load text from a file."""
    with open(filename, 'r', encoding='utf-8') as file:
        return file.read()

def generate_typing_data_from_text(text, n_sessions, n_keys=100, interval_type='constant', base_interval=0.5, variability=0):
    """Generate typing sessions from actual text with specified interval types and speeds."""
    sessions = []
    text_length = len(text)
    for _ in range(n_sessions):
        start_index = np.random.randint(0, text_length - n_keys)
        selected_text = text[start_index:start_index + n_keys]
        start_time = datetime.now()
        events = []
        for i, char in enumerate(selected_text):
            if interval_type == 'constant':
                interval = base_interval
            else:  # random interval
                interval = np.random.uniform(base_interval - variability, base_interval + variability)
            timestamp = start_time + timedelta(seconds=i * interval)
            events.append({'key': char, 'timestamp': timestamp.isoformat()})
        sessions.append({'events': events})
    return sessions

def save_data(sessions, filename):
    """Save sessions to a JSON file."""
    with open(filename, 'w') as file:
        json.dump(sessions, file, indent=4)

# Load real text (ensure the text file is in the correct path)
text = load_text('sample_text.txt')

# Typing speed settings (in seconds per keypress)
speed_settings = [
    {'type': 'constant', 'base_interval': 0.1},   # Very fast typing (bots)
    {'type': 'constant', 'base_interval': 0.2},   # Fast human typing
    {'type': 'constant', 'base_interval': 0.4},   # Average human typing
    {'type': 'constant', 'base_interval': 0.6},   # Slow human typing
    {'type': 'random', 'base_interval': 0.2, 'variability': 0.1},   # Variable fast typing
    {'type': 'random', 'base_interval': 0.4, 'variability': 0.2},   # Variable average typing
    {'type': 'random', 'base_interval': 0.6, 'variability': 0.3}    # Variable slow typing
]

# Generate and save data for each speed setting
for setting in speed_settings:
    data = generate_typing_data_from_text(text, n_sessions=50, n_keys=100, interval_type=setting['type'],
                                          base_interval=setting['base_interval'], variability=setting.get('variability', 0))
    filename = f'samples/bot/{setting["type"]}_bot_data_speed_{setting["base_interval"]}{"_var" if "variability" in setting else ""}.json'
    save_data(data, filename)
