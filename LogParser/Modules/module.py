import json
import sys
from enum import Enum
from typing import List, Optional
from helper import parse_config

parsed_config = parse_config("module.config")

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
    def __init__(self, status: StatusEnum, results: Optional[List[Metric]] = None, errors: Optional[List[Error]] = None):
        self.status = status
        self.results = results
        self.errors = errors
def compare_version_or_build(current, expected):
    current_parts = [int(part) for part in current.split('.')]
    expected_parts = [int(part) for part in expected.split('.')]
    
# Compare each part
    for curr, exp in zip(current_parts, expected_parts):
        if curr < exp:
            return False
        elif curr > exp:
            return True
    
    # If all parts are equal but current has more parts
    return len(current_parts) >= len(expected_parts)
    
def check_for_suspicious_activity(log_data: dict) -> List[Error]:
    errors = []

    # Check for suspicious application info
    app_info = log_data.get('application', {})
    if app_info.get('name') == "Safe Exam Browser":
        version = app_info.get('version')
        build = app_info.get('build')

        expected_version = parsed_config["min_version"]
        expected_build = parsed_config["min_build"]
        # Compare the current version and build with the expected version and build
        
        if not compare_version_or_build(version, expected_version) or not compare_version_or_build(build, expected_build):
            errors.append(Error(
                message="Suspicious Application Version or Build",
                details=f"Expected version >= {expected_version} and build >= {expected_build} but found version {version} and build {build}"
            ))


    # Check for suspicious system info
    system_info = log_data.get('system', {})
    suspicious_keywords = parsed_config["suspicious_keywords"]

    os_name = system_info.get('os', '').lower()
    model = system_info.get('model', '').lower()
    manufacturer = system_info.get('manufacturer', '').lower()

    if any(keyword in os_name for keyword in suspicious_keywords):
        errors.append(Error(
            message="Suspicious Operating System",
            details=f"Operating system contains suspicious keywords: {os_name}"
        ))

    if any(keyword in model for keyword in suspicious_keywords):
        errors.append(Error(
            message="Suspicious Motherboard Model",
            details=f"Motherboard model contains suspicious keywords: {model}"
        ))

    if any(keyword in manufacturer for keyword in suspicious_keywords):
        errors.append(Error(
            message="Suspicious Manufacturer",
            details=f"Manufacturer contains suspicious keywords: {manufacturer}"
        ))

    return errors

def process_log(log_path: str) -> ReturnJSON:
    try:
        with open(log_path, 'r') as file:
            log_data = json.load(file)

        # Check for suspicious activity
        errors = check_for_suspicious_activity(log_data)

        status = StatusEnum.CHEATING if errors else StatusEnum.SUCCESS

        return ReturnJSON(status=status, errors=errors if errors else None)

    except Exception as e:
        return ReturnJSON(status=StatusEnum.ERROR, errors=[Error(message=str(e))])

def main():
    if len(sys.argv) != 2:
        print("Usage: python script.py <path_to_log_file>")
        sys.exit(1)

    log_path = sys.argv[1]
    result = process_log(log_path)

    print(json.dumps(result, default=lambda o: o.__dict__, indent=2))

if __name__ == "__main__":
    main()
