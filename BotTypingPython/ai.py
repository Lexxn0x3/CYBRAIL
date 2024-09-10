import os
import json
import numpy as np
from datetime import datetime
from tensorflow.keras.models import Sequential
from tensorflow.keras.layers import LSTM, Dense
from tensorflow.keras.callbacks import EarlyStopping, ModelCheckpoint
from sklearn.model_selection import train_test_split
import matplotlib.pyplot as plt
from tensorflow.keras.models import load_model
from tensorflow.keras.layers import Input
import tensorflow as tf
import shap


def load_data(directory, label):
    """Load and label data from JSON files in a specified directory."""
    sessions = []
    for filename in os.listdir(directory):
        filepath = os.path.join(directory, filename)
        with open(filepath, 'r') as file:
            data = json.load(file)
            for session in data:
                timestamps = [datetime.fromisoformat(event['timestamp']) for event in session['events']]
                if len(timestamps) > 1:
                    intervals = [t2 - t1 for t1, t2 in zip(timestamps[:-1], timestamps[1:])]
                    intervals = [interval.total_seconds() for interval in intervals]
                    sessions.append((intervals, label))
    return sessions


def create_sequences(sessions, n_steps=10):
    """Create sequences and labels from interval data."""
    X, y = [], []
    for intervals, label in sessions:
        for i in range(len(intervals) - n_steps):
            X.append(intervals[i:i + n_steps])
            y.append(label)
    return np.array(X), np.array(y)


def build_model(input_shape):
    """Build an LSTM model for classification."""
    model = Sequential([
        Input(shape=input_shape),
        LSTM(50, activation='relu'),
        Dense(1, activation='sigmoid')
    ])
    model.compile(optimizer='adam', loss='binary_crossentropy', metrics=['accuracy'])
    return model


def plot_learning_curves(history):
    """Plot training and validation accuracy and loss."""
    plt.figure(figsize=(12, 4))

    # Accuracy plot
    plt.subplot(1, 2, 1)
    plt.plot(history.history['accuracy'], label='Training Accuracy')
    plt.plot(history.history['val_accuracy'], label='Validation Accuracy')
    plt.title('Model Accuracy')
    plt.ylabel('Accuracy')
    plt.xlabel('Epoch')
    plt.legend(loc='upper left')

    # Loss plot
    plt.subplot(1, 2, 2)
    plt.plot(history.history['loss'], label='Training Loss')
    plt.plot(history.history['val_loss'], label='Validation Loss')
    plt.title('Model Loss')
    plt.ylabel('Loss')
    plt.xlabel('Epoch')
    plt.legend(loc='upper left')

    plt.show()


# Load data from both human and bot directories
human_data = load_data('samples/human', 0)
bot_data = load_data('samples/bot', 1)

# Combine human and bot data
all_data = human_data + bot_data

# Prepare sequences
X, y = create_sequences(all_data)

# Split the data into training and testing sets
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

# Reshape input to be [samples, time steps, features] for LSTM
X_train = X_train.reshape((X_train.shape[0], X_train.shape[1], 1))
X_test = X_test.reshape((X_test.shape[0], X_test.shape[1], 1))

# Build the model
model = build_model((X_train.shape[1], 1))

# Set up callbacks
early_stopping = EarlyStopping(monitor='val_loss', patience=3, restore_best_weights=True)
model_checkpoint = ModelCheckpoint('saved_model/best_model.keras', monitor='val_loss', save_best_only=True)

if False:
    # Train the model with early stopping and model checkpointing
    history = model.fit(X_train, y_train, epochs=3, batch_size=32, validation_data=(X_test, y_test),
                        callbacks=[early_stopping, model_checkpoint])

    # Plot the learning curves
    plot_learning_curves(history)

# Load the best model saved during training
best_model = load_model('saved_model/best_model.keras')

# Evaluate the model on the test data
evaluation = best_model.evaluate(X_test, y_test)
print(f'Test Loss: {evaluation[0]}, Test Accuracy: {evaluation[1]}')

# Reshape X_train and X_test to 2D arrays (samples, flattened time steps * features)
X_train_reshaped = X_train.reshape((X_train.shape[0], -1))  # Flatten each sample's time steps into a 2D array
X_test_reshaped = X_test.reshape((X_test.shape[0], -1))     # Same for test data

# Use KernelExplainer as a general solution for SHAP explanations
explainer = shap.KernelExplainer(best_model.predict, X_train_reshaped[:100])

# Explain predictions on the test set (you can increase the sample size if needed)
shap_values = explainer.shap_values(X_test_reshaped[:1])

# Reshape SHAP values to match the test data
shap_values_flattened = np.reshape(shap_values[0], (1, -1))

# Ensure the shapes match before plotting
print(f"Flattened SHAP values shape: {shap_values_flattened.shape}")
print(f"Test data shape: {X_test_reshaped[:1].shape}")

# Plot SHAP summary with flattened SHAP values
shap.summary_plot(shap_values_flattened, X_test_reshaped[:1])