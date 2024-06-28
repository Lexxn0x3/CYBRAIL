import json
from enum import Enum
from typing import List, Optional


class StatusEnum(str, Enum):
    SUCCESS = "success"
    CHEATING = "cheating"
    ERROR = "error"


class Error:
    def __init__(self, message: str, details: Optional[str] = None):
        self.message = message
        self.details = details


class ReturnJSON:
    def __init__(self, status: StatusEnum, errors: Optional[List[Error]] = None):
        self.status = status
        self.errors = errors


def load_log_file(log_path: str) -> dict:
    with open(log_path, 'r') as file:
        return json.load(file)


def output_result(result: ReturnJSON):
    print(json.dumps(result, default=lambda o: o.__dict__, indent=2))
