import sys
from typing import List, Optional
from log_checker import StatusEnum, Error, ReturnJSON, load_log_file, output_result


def check_network_config_changes(log_data: dict) -> List[Error]:
    errors = []
    network_change_detected = False

    for log_entry in log_data.get('logs', []):
        if log_entry['level'] == 'INFO' and 'Network configuration changed' in log_entry['message']:
            network_change_detected = True

    if network_change_detected:
        errors.append(Error(
            message="Network Configuration Change Detected",
            details="A change in network configuration was detected during the session."
        ))

    return errors


def process_log(log_path: str) -> ReturnJSON:
    try:
        log_data = load_log_file(log_path)

        # Check for network configuration changes
        errors = check_network_config_changes(log_data)

        status = StatusEnum.CHEATING if errors else StatusEnum.SUCCESS

        return ReturnJSON(status=status, errors=errors if errors else None)

    except Exception as e:
        return ReturnJSON(status=StatusEnum.ERROR, errors=[Error(message=str(e))])


def main():
    if len(sys.argv) != 2:
        print("Usage: python network_config_check.py <path_to_log_file>")
        sys.exit(1)

    log_path = sys.argv[1]
    result = process_log(log_path)

    output_result(result)


if __name__ == "__main__":
    main()
