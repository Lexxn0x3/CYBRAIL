import json
import sys
from enum import Enum
from typing import List, Optional

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

def check_for_suspicious_activity(log_data: dict) -> List[Error]:
    errors = []

    # Check for suspicious application info
    app_info = log_data.get('application', {})
    if app_info.get('name') == "Safe Exam Browser":
        version = app_info.get('version')
        build = app_info.get('build')
        if version != "3.7.0" or build != "3.7.0.682":
            errors.append(Error(
                message="Suspicious Application Version or Build",
                details=f"Expected version 3.7.0 and build 3.7.0.682 but found version {version} and build {build}"
            ))

    # Check for suspicious system info
    system_info = log_data.get('system', {})
    suspicious_keywords = suspicious_keywords = [
        "vm", "vmware", "virtualbox", "qemu", "xen", "hyper-v", "parallels",
        "vbox", "kvm", "virt", "virtual", "cloud", "unknown", "generic",
        "standard", "emulated", "bochs", "innotek",
        "guest"
    ]
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
