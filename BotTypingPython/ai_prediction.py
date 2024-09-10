import json
import numpy as np
from tensorflow.keras.models import load_model

# Load the saved model
model = load_model('saved_model/best_model.keras')


def load_json_log(filename):
    """Load typing data from a JSON log and return processed intervals."""
    with open(filename, 'r') as file:
        data = json.load(file)

    all_intervals = []
    for session in data["events"]:
        # Extract intervals directly from the session
        intervals = session["intervals"]

        if len(intervals) > 0:
            all_intervals.append(np.array(intervals))

    return all_intervals


def prepare_data_for_prediction(intervals, n_steps=10):
    """Create sequences of intervals to match the input format of the model."""
    X = []
    for interval_seq in intervals:
        for i in range(len(interval_seq) - n_steps):
            X.append(interval_seq[i:i + n_steps])
    return np.array(X).reshape((len(X), n_steps, 1))  # Reshape for LSTM


def predict_bot_likelihood(model, X):
    """Use the trained model to predict bot-likelihood for each sequence."""
    predictions = model.predict(X)
    return predictions


def overall_bot_score(predictions):
    """Calculate the overall bot-likelihood score."""
    return np.mean(predictions)  # Average prediction score


# Load the JSON log of typing events
log_filename = 'samples/interval/human/typing_session_20240909_141805.json'  # Replace with your log file path
intervals = load_json_log(log_filename)

# Prepare the data for prediction (assuming model was trained with n_steps=10)
X = prepare_data_for_prediction(intervals, n_steps=10)

# Predict bot likelihood for each event
predictions = predict_bot_likelihood(model, X)

# Calculate overall bot score
score = overall_bot_score(predictions)

# Output results
print(f"Overall Bot Likelihood Score: {score:.4f}")

# Optionally, you can print each prediction if needed:
for i, prediction in enumerate(predictions):
    print(f"Event {i + 1}: Bot Likelihood = {prediction[0]:.4f}")
