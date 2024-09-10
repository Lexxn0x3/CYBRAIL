import keyboard
import json
from datetime import datetime, timedelta
import os

class TypingSessionCollector:
    def __init__(self, timeout_seconds):
        self.timeout_seconds = timeout_seconds
        self.sessions = []
        self.current_session = []
        self.last_event_time = datetime.now()

    def key_event(self, event):
        """Capture key events and organize them into sessions."""
        current_time = datetime.now()
        if (current_time - self.last_event_time).total_seconds() > self.timeout_seconds:
            # Start a new session if the timeout has been exceeded
            if self.current_session:
                self.sessions.append({'events': self.current_session})
                self.current_session = []
        self.last_event_time = current_time
        self.current_session.append({'key': event.name, 'timestamp': current_time.isoformat()})

    def save_sessions(self, filename):
        """Save the captured sessions to a JSON file."""
        with open(filename, 'w') as f:
            json.dump(self.sessions, f, indent=4)

def main():
    timeout_seconds = 2  # Pause threshold to start a new session
    collector = TypingSessionCollector(timeout_seconds)

    print("Start typing... (press ESC to stop and save the session)")
    keyboard.hook(collector.key_event)
    keyboard.wait('esc')  # Wait until ESC is pressed

    # Save the last session if it exists
    if collector.current_session:
        collector.sessions.append({'events': collector.current_session})

    # Save sessions to a file
    filename = os.path.join('samples', 'human', f'typing_session_{datetime.now().strftime("%Y%m%d_%H%M%S")}.json')
    os.makedirs(os.path.dirname(filename), exist_ok=True)
    collector.save_sessions(filename)
    print(f"Saved sessions to {filename}")

if __name__ == "__main__":
    main()
