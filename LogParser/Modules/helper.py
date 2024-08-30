import json
import os

def parse_config(path):
    current_directory = os.path.dirname(os.path.abspath(__file__))

    config_file_path = os.path.join(current_directory, path)
    with open(config_file_path, 'r') as file:
        config_data = json.load(file)

    parsed_config = {}
    for key, value in config_data.items():
        # Extract the type and value
        data_type = value.get('type')
        data_value = value.get('value')
        
        # Convert the value to the correct type
        if data_type == 'string':
            parsed_config[key] = str(data_value)
        elif data_type == 'int':
            parsed_config[key] = int(data_value)
        elif data_type == 'float':
            parsed_config[key] = float(data_value)
        elif data_type == 'bool':
            parsed_config[key] = bool(data_value)
        elif data_type == 'list':
            parsed_config[key] = list(data_value)
        # Add other data types as needed
        else:
            parsed_config[key] = data_value  # default case, keep the original value
    
    return parsed_config

