import os
import sys
import json
import numpy as np

# Suppress TensorFlow and CUDA log messages
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3'  # Suppress INFO and WARNING messages from TensorFlow
os.environ['CUDA_VISIBLE_DEVICES'] = '-1'  # Disable GPU entirely
os.environ['TF_ENABLE_ONEDNN_OPTS'] = '0'  # Disable oneDNN custom operations to avoid those warnings

import tensorflow as tf
tf.get_logger().setLevel('ERROR')  # Set TensorFlow logger to show only errors

from tensorflow.keras.models import load_model
from enum import Enum
from typing import Optional, List
import configparser

from helper import parse_config
# Enums and Classes for Status
class StatusEnum(str, Enum):
    SUCCESS = "success"
    CHEATING = "cheating"
    ERROR = "error"


class CustomField:
    def __init__(self, name: str, property: str):
        self.name = name
        self.property = property


class Metric:
    def __init__(self, short_info: str, long_info: str, name: str):
        self.short_info = short_info
        self.long_info = long_info
        self.name = name


class Error:
    def __init__(self, message: str, details: Optional[str] = None):
        self.message = message
        self.details = details


class ReturnJSON:
    def __init__(self, status: StatusEnum, results: Optional[List[Metric]] = None,
                 errors: Optional[List[Error]] = None):
        self.status = status
        self.results = results
        self.errors = errors


# Load the JSON log data
def load_json_log(filename: str):
    try:
        with open(filename, 'r') as file:
            data = json.load(file)

        all_intervals = []
        for session in data["events"]:
            intervals = session["intervals"]
            if len(intervals) > 0:
                all_intervals.append(np.array(intervals))

        return all_intervals
    except Exception as e:
        raise ValueError(f"Error loading JSON log: {str(e)}")


# Prepare the intervals data for model input
def prepare_data_for_prediction(intervals, n_steps=10):
    try:
        X = []
        for interval_seq in intervals:
            for i in range(len(interval_seq) - n_steps):
                X.append(interval_seq[i:i + n_steps])
        return np.array(X).reshape((len(X), n_steps, 1))  # Reshape for LSTM
    except Exception as e:
        raise ValueError(f"Error preparing data for prediction: {str(e)}")


# Predict the likelihood of bot typing
def predict_bot_likelihood(model, X):
    try:
        predictions = model.predict(X, verbose=0)
        return predictions
    except Exception as e:
        raise ValueError(f"Error predicting bot likelihood: {str(e)}")


# Check for suspicious activity based on bot likelihood score
def check_for_suspicious_activity(predictions, threshold=0.8):
    errors = []
    avg_score = np.mean(predictions)
    if avg_score > threshold:
        errors.append(
            Error(message="Suspicious bot-like typing detected", details=f"Average bot likelihood: {avg_score*100:.2f}%"))
    return errors


# Main function to process the log and return the result
def process_log(log_path: str):
    try:
        # Parse the configuration
        config_path = os.path.join(os.path.dirname(__file__), 'bot_typing.config')
        parsed_config = parse_config(config_path)

        # Load the model
        model_path = os.path.join(os.path.dirname(__file__), parsed_config["model_path"])
        model = load_model(model_path)

        # Load log data
        log_data = load_json_log(log_path)

        # Prepare the data for prediction
        X = prepare_data_for_prediction(log_data, n_steps=10)

        # Predict bot likelihood
        predictions = predict_bot_likelihood(model, X)

        # Check for suspicious activity
        errors = check_for_suspicious_activity(predictions, threshold=parsed_config["threshold"])

        # Determine status
        status = StatusEnum.CHEATING if errors else StatusEnum.SUCCESS

        # Build results
        results = [Metric(
            short_info="Bot Detection Score",
            long_info=f"Average likelihood of bot typing: {np.mean(predictions):.4f}",
            name="bot_detection_score"
        )]

        return ReturnJSON(status=status, results=results if not errors else None, errors=errors if errors else None)

    except Exception as e:
        return ReturnJSON(status=StatusEnum.ERROR, errors=[Error(message=str(e))])


# Main entry point for the script
if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python script.py <path_to_log_file>")
        sys.exit(1)

    log_path = sys.argv[1]
    result = process_log(log_path)

    print(json.dumps(result, default=lambda o: o.__dict__, indent=2))
