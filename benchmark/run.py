#!/usr/bin/env python3
import shutil
import subprocess
import pathlib
import os
import argparse
from shutil import copytree
from time import time
from typing import List
from glob import glob

PROJ_ROOT_DIR = os.path.dirname(os.path.dirname(os.path.realpath(__file__)))
print(f"project root: {PROJ_ROOT_DIR}")

INST_SCRIPT = os.path.join(PROJ_ROOT_DIR, "goFuzz/scripts/instrument.py")
PATCH_RUNTIME_SCRIPT = os.path.join(PROJ_ROOT_DIR, "scripts/patch-go-runtime.sh")

BENCHMARK_DIR = os.path.join(PROJ_ROOT_DIR, "benchmark")
STD_INPUT_FILE = os.path.join(BENCHMARK_DIR, "std-input")

TEST_PKG = os.path.join(BENCHMARK_DIR, "tests")


TMP_FOLDER = os.path.join(PROJ_ROOT_DIR, "benchmark/tmp")
TEST_PKG_INST_TMP = os.path.join(TMP_FOLDER, "instrument")
TEST_BIN_NATIVE = os.path.join(TMP_FOLDER, "native.test")
TEST_BIN_INST = os.path.join(TMP_FOLDER, "inst.test")

STD_RECORD_FILE = os.path.join(TMP_FOLDER, "record")
STD_OUTPUT_FILE = os.path.join(TMP_FOLDER, "output")

class BinTest:
    def __init__(self, bin:str, func:str) -> None:
        self.bin = bin
        self.func = func

def benchmark_simple(run_mode:str):
    """evaluate performance by using tests from `benchmark/tests` folder

    Args:
        run_mode (str): `native`, `inst`
            native: run without patching golang runtime and insturmenting code
            inst: run with patched golang runtime & instrumented code
    """
    
    # prepare test bin
    if run_mode == "inst":
        patch_go_runtime()
        copytree(TEST_PKG, TEST_PKG_INST_TMP)
        inst_dir(TEST_PKG_INST_TMP)
        compile_test_bin(TEST_PKG_INST_TMP, TEST_BIN_INST)
        target_bin = TEST_BIN_INST
    else:
        compile_test_bin(TEST_PKG, TEST_BIN_NATIVE)
        target_bin = TEST_BIN_NATIVE
    
    tests = get_tests_from_test_bin(target_bin)
    run_benchmark_with_tests(tests, run_mode)

def benchmark_custom(dir:str, mode:str, selected_bins: List[str] = None):
    if not selected_bins:
        selected_bins = glob(f"{dir}/*")

    tests = []

    for bin in selected_bins:
        ts = get_tests_from_test_bin(bin)
        tests.extend(ts)
    
    run_benchmark_with_tests(tests, mode)

def benchmark_custom_native_parallel(dir:str, selected_bins: List[str] = None):
    if not selected_bins:
        selected_bins = glob(f"{dir}/*")

    total_num_of_tests = 0
    total_duration = 0
    for bin in selected_bins:
        ts = get_tests_from_test_bin(bin)
        num_of_tests = len(ts)
        dur = benchmark(1, lambda: subprocess.run([bin, "-test.timeout", "10s", "-test.parallel","5"]))
        total_num_of_tests += num_of_tests
        total_duration += dur
    
    print(f"total {total_num_of_tests} tests, total {total_duration} seconds")

    
    
def run_benchmark_with_tests(tests: List[BinTest], mode:str):
    if mode == "inst":
        inst_run_env = os.environ.copy()
        inst_run_env["GF_BENCHMARK"] = "1"
    total_dur = 0
    tests_dur = {}
    for t in tests:
        if mode == "inst":
            dur = benchmark(10, lambda: subprocess.run(
                [t.bin, "-test.timeout", "10s","-test.run", t.func], 
                env=inst_run_env, timeout=10
                ))
        elif mode == "native":
            dur = benchmark(10, lambda: subprocess.run([t.bin, "-test.timeout", "10s", "-test.run", t.func]))
        full_name = f"{t.bin}->{t.func}"
        if dur == -1:
            print(f"{full_name}: timeout")
            continue
        tests_dur[full_name] = dur
        total_dur += dur
    print(f"total {len(tests)} tests")
    for k, v in tests_dur.items():
        print(f"{k}:{v:.04f} seconds")
    print(f"total average {total_dur/len(tests):.04f} seconds / test")


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('testsuite', choices=["simple", "custom"])
    parser.add_argument('--mode', choices=["native", "inst", "parallel-native"], required=True)
    parser.add_argument('--dir', type=str)
    parser.add_argument('--bins',  nargs='*')
    parser.add_argument('--bins-list-file', type=str)
    args = parser.parse_args()

    print(args)

    bins_dir = args.dir
    selected_bins = args.bins
    mode = args.mode
    testsuite = args.testsuite
    bins_list_file = args.bins_list_file

    if testsuite == "custom" and not bins_dir:
        raise Exception("--dir is required if testsuite is custom")

    if mode == "parallel-native":
        return benchmark_custom_native_parallel(bins_dir, selected_bins)
    
    if testsuite == "simple":
        benchmark_simple(mode)
    else:
        if bins_list_file:
            bins = get_bins_from_file(bins_list_file)
            selected_bins = [os.path.join(bins_dir, bin) for bin in bins]
        benchmark_custom(bins_dir, mode, selected_bins)

def get_bins_from_file(file: str)->List[str]:
    with open(file, "r") as f:
        bins = f.read().splitlines()
    return bins

def inst_dir(dir: str):
    subprocess.run([INST_SCRIPT, dir]).check_returncode()

def compile_test_bin(pkg_dir:str, dest:str):
    subprocess.run(
        ["go","test","-c", "-o", dest, pkg_dir], 
        cwd=BENCHMARK_DIR
    ).check_returncode()

def patch_go_runtime():
    subprocess.run([PATCH_RUNTIME_SCRIPT]).check_returncode()

def restore_inst_run(std_input_content:str):
    if os.path.exists(STD_RECORD_FILE):
        os.remove(STD_RECORD_FILE)
    
    if os.path.exists(STD_OUTPUT_FILE):
        os.remove(STD_OUTPUT_FILE)

    with open(STD_INPUT_FILE, 'w') as f:
        f.write(std_input_content)

def get_tests_from_test_bin(bin: str) -> List[BinTest]:
    p = subprocess.run([bin, '-test.list', '.*'], stdout=subprocess.PIPE)
    p.check_returncode()
    out = p.stdout.decode('utf-8')
    funcs = out.splitlines()
    print("found tests", funcs)
    return [BinTest(bin, func) for func in funcs]

def benchmark(reps, func, prefunc=None):
    dur = 0
    cnt = 0
    for _ in range(0, reps):
        try:
            if prefunc:
                prefunc()
            start = time()
            func()
            end = time()
            dur += (end-start)
            cnt += 1
        except:
            pass
    if cnt == 0:
        return -1
    return dur / reps

if __name__ == "__main__":
    main()