import os
import json
import numpy as np
from datetime import datetime


def load_data_and_generate_intervals(directory, output_directory, interval_threshold=2.0):
    """Load data from JSON files, compute intervals, and group them into events based on typing pauses."""
    if not os.path.exists(output_directory):
        os.makedirs(output_directory)

    for filename in os.listdir(directory):
        filepath = os.path.join(directory, filename)
        with open(filepath, 'r') as file:
            data = json.load(file)

            all_grouped_events = []  # This will hold all grouped events from multiple sessions

            for session in data:
                events = session['events']

                # Extract timestamps and compute intervals
                timestamps = [datetime.fromisoformat(event['timestamp']) for event in events]
                intervals = np.diff(timestamps).astype('timedelta64[ms]').astype(np.float64) / 1000.0  # Convert to seconds

                current_event = []
                for i in range(len(intervals)):
                    interval = intervals[i]

                    # Start a new event if interval exceeds the threshold
                    if interval > interval_threshold and current_event:
                        all_grouped_events.append({"intervals": current_event})
                        current_event = []

                    current_event.append(interval)

                # Add the last event if any intervals remain
                if current_event:
                    all_grouped_events.append({"intervals": current_event})

            # Create the final structure as a single object with "events"
            output_data = {"events": all_grouped_events}

            # Save the modified session back to a new JSON file in the output directory
            output_filepath = os.path.join(output_directory, filename)
            with open(output_filepath, 'w') as output_file:
                json.dump(output_data, output_file, indent=4)


# Example usage:
# Load data from 'samples/human' and 'samples/bot', and generate grouped interval-only JSON files
load_data_and_generate_intervals('samples/human', 'samples/interval/human')
load_data_and_generate_intervals('samples/bot', 'samples/interval/bot')
