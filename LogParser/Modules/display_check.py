import sys
from log_checker import StatusEnum, Error, ReturnJSON, load_log_file, output_result

def check_multiple_displays(log_data: dict) -> list:
    errors = []
    active_displays = 0
    allowed_displays = 1

    for log_entry in log_data.get('logs', []):
        if log_entry['level'] == 'INFO':
            if 'Detected active, external display' in log_entry['message']:
                active_displays += 1
            elif 'Detected' in log_entry['message'] and 'active displays' in log_entry['message']:
                try:
                    parts = log_entry['message'].split()
                    active_displays = int(parts[2])
                    allowed_displays = int(parts[5])  # Correctly extract the number of allowed displays
                except (ValueError, IndexError) as e:
                    errors.append(Error(
                        message="Failed to parse display information",
                        details=f"Log entry: {log_entry['message']} - Error: {str(e)}"
                    ))

    if active_displays > allowed_displays:
        errors.append(Error(
            message="Multiple Displays Detected",
            details=f"{active_displays} active displays detected, but only {allowed_displays} are allowed."
        ))

    return errors

def process_log(log_path: str) -> ReturnJSON:
    try:
        log_data = load_log_file(log_path)

        # Check for multiple displays
        errors = check_multiple_displays(log_data)

        status = StatusEnum.CHEATING if errors else StatusEnum.SUCCESS

        return ReturnJSON(status=status, errors=errors if errors else None)

    except Exception as e:
        return ReturnJSON(status=StatusEnum.ERROR, errors=[Error(message="Exception occurred", details=str(e))])

def main():
    if len(sys.argv) != 2:
        print("Usage: python multiple_displays_check.py <path_to_log_file>")
        sys.exit(1)

    log_path = sys.argv[1]
    result = process_log(log_path)

    output_result(result)

if __name__ == "__main__":
    main()
