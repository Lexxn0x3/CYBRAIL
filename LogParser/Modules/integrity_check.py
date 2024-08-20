import json
import sys
from enum import Enum
from typing import List, Optional

class StatusEnum(str, Enum):
    SUCCESS = "success"
    INTEGRITY_FAILURE = "cheating"
    ERROR = "error"

class Error:
    def __init__(self, message: str, details: Optional[str] = None):
        self.message = message
        self.details = details

class ReturnJSON:
    def __init__(self, status: StatusEnum, errors: Optional[List[Error]] = None):
        self.status = status
        self.errors = errors

def check_integrity(log_data: dict) -> List[Error]:
    errors = []
    integrity_start = False
    integrity_success = False

    for log_entry in log_data.get('logs', []):
        if log_entry['level'] == 'INFO' and 'Attempting to verify application integrity' in log_entry['message']:
            integrity_start = True
        if log_entry['level'] == 'INFO' and 'Application integrity successfully verified' in log_entry['message']:
            integrity_success = True

    if integrity_start and not integrity_success:
        errors.append(Error(
            message="Integrity Check Failed",
            details="Integrity verification was attempted but not successfully completed."
        ))
    
    if not integrity_start:
        errors.append(Error(
            message="Integrity Check Not Found",
            details="No attempt to verify application integrity was found in the logs."
        ))

    return errors

def process_log(log_path: str) -> ReturnJSON:
    try:
        with open(log_path, 'r') as file:
            log_data = json.load(file)

        # Check for integrity verification
        errors = check_integrity(log_data)

        status = StatusEnum.INTEGRITY_FAILURE if errors else StatusEnum.SUCCESS

        return ReturnJSON(status=status, errors=errors if errors else None)

    except Exception as e:
        return ReturnJSON(status=StatusEnum.ERROR, errors=[Error(message=str(e))])

def main():
    if len(sys.argv) != 2:
        print("Usage: python integrity_check.py <path_to_log_file>")
        sys.exit(1)

    log_path = sys.argv[1]
    result = process_log(log_path)

    print(json.dumps(result, default=lambda o: o.__dict__, indent=2))

if __name__ == "__main__":
    main()
