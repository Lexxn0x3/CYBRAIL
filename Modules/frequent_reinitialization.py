import sys
from log_checker import StatusEnum, Error, ReturnJSON, load_log_file, output_result
from typing import List

def check_frequent_reinitialization(log_data: dict) -> List[Error]:
    errors = []
    reinitialization_count = 0

    for log_entry in log_data.get('logs', []):
        if log_entry['level'] == 'INFO' and 'Initializing new session configuration' in log_entry['message']:
            reinitialization_count += 1

    if reinitialization_count > 1:
        errors.append(Error(
            message="Frequent Reinitialization Detected",
            details=f"Multiple session initializations detected: {reinitialization_count} times."
        ))

    return errors

def process_log(log_path: str) -> ReturnJSON:
    try:
        log_data = load_log_file(log_path)

        # Check for frequent reinitialization
        errors = check_frequent_reinitialization(log_data)

        status = StatusEnum.CHEATING if errors else StatusEnum.SUCCESS

        return ReturnJSON(status=status, errors=errors if errors else None)

    except Exception as e:
        return ReturnJSON(status=StatusEnum.ERROR, errors=[Error(message=str(e))])

def main():
    if len(sys.argv) != 2:
        print("Usage: python frequent_reinitialization.py <path_to_log_file>")
        sys.exit(1)

    log_path = sys.argv[1]
    result = process_log(log_path)

    output_result(result)

if __name__ == "__main__":
    main()
