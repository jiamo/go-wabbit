

import pathlib
import pytest
import random
from wabbit.parse import parse_file
right_program_path = pathlib.Path(__file__).resolve().parent / "Parser"

right_files = list(right_program_path.glob("*.wb"))
random.shuffle(right_files)


def test_right_program():
    # print(right_program_path)
    # print(list(right_files))
    for right_file in right_files:
        print(f"handing ", right_file)
        parse_file(right_file)
