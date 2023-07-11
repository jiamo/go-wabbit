

import pathlib
import pytest
from wabbit.parse import parse_file
import random
from wabbit.llvm import main
right_program_path = pathlib.Path(__file__).resolve().parent / "Programs"
right_files = list(right_program_path.glob("*.wb"))
random.shuffle(right_files)

def test_right_program():
    # print(right_program_path)
    # print(list(right_files))
    for right_file in right_files:
        print(f"handing ", right_file)
        main(right_file)
