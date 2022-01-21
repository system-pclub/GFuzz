#!/usr/bin/env python3
import shutil
import subprocess
import pathlib
import os
import argparse
from shutil import copytree
from time import time
from typing import List, Tuple
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

REPEAT = 10
class BinTest:
    def __init__(self, bin:str, func:str) -> None:
        self.bin = bin
        self.func = func


def benchmark_custom(out:str, dir:str, mode:str, selected_bins: List[str] = None):
    if not selected_bins:
        selected_bins = glob(f"{dir}/*")
    tests = []

    for bin in selected_bins:
        ts = get_tests_from_test_bin(bin)
        tests.extend(ts)
    
    run_benchmark_with_tests(out, tests, mode)

def benchmark_custom_native_parallel(dir:str, selected_bins: List[str] = None) -> Tuple[int, int]:
    if not selected_bins:
        selected_bins = glob(f"{dir}/*")

    total_num_of_tests = 0
    total_duration = 0
    for bin in selected_bins:
        ts = get_tests_from_test_bin(bin)
        num_of_tests = len(ts)
        dur = benchmark(1, lambda: subprocess.run([bin, "-test.timeout=5m", "-test.parallel","5"]))
        total_num_of_tests += num_of_tests
        total_duration += dur
    
    return (total_num_of_tests, total_duration)

def benchmark_all_native_parallel(dirs:List[str], repeat=3) -> Tuple[int, int]:
    # 1. find common tests
    #    in each dirs: find common both in inst/* and native/*
    common_bins = []
    for d in dirs:
        native_bins = glob(f"{d}/native/*")
        for nb in native_bins:
            filename = os.path.basename(nb)
            inst_version = os.path.join(d, "inst", filename)
            if os.path.exists(inst_version):
                common_bins.append(nb)

    print(f"total {len(common_bins)} bins")
    
    # 2. loop three times by executing with parallel=5
    total_num_of_tests = 0
    total_duration = 0
    for _ in range(repeat):
        num_of_tests, dur = benchmark_custom_native_parallel(None, common_bins)
        total_duration += dur
        total_num_of_tests += num_of_tests
    
    return (total_num_of_tests, total_duration) 


    
    
def run_benchmark_with_tests(out:str, tests: List[BinTest], mode:str):
    if mode == "inst":
        inst_run_env = os.environ.copy()
        inst_run_env["GF_BENCHMARK"] = "1"
    total_dur = 0
    tests_dur = {}
    global REPEAT
    for t in tests:
        strict_func_name = "^"+t.func+"$"
        if mode == "inst":
            dur = benchmark(REPEAT, lambda: subprocess.run(
                [t.bin, "-test.timeout", "10s","-test.run", strict_func_name], 
                env=inst_run_env, timeout=10
                ))
        elif mode == "native":
            dur = benchmark(REPEAT, lambda: subprocess.run([t.bin, "-test.timeout", "10s", "-test.run", strict_func_name], timeout=10))
        full_name = f"{t.bin}->{t.func}"
        if dur == -1:
            print(f"{full_name}: timeout")
            continue
        tests_dur[full_name] = dur
        total_dur += dur
    print(f"total {len(tests)} tests")

    if not os.path.exists(os.path.dirname(out)):
        os.makedirs(os.path.dirname(out))

    with open(out, "w") as outf:
        for k, v in tests_dur.items():
            outf.write(f"{k}:{v:.2f} seconds\n")

    if len(tests) == 0:
        return print(f"no test found/run")

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('action', choices=["count-tests","benchmark"])
    #parser.add_argument('testsuite', choices=["simple", "custom"])
    parser.add_argument('--mode', choices=["native", "inst", "parallel-native"])
    parser.add_argument('--dir', type=str)
    parser.add_argument('--out', type=str)
    parser.add_argument('--bins',  nargs='*')
    parser.add_argument('--bins-list-file', type=str)
    parser.add_argument('--repeat', type=int, default=10)
    args = parser.parse_args()

    print(args)

    bins_dir = args.dir
    selected_bins = args.bins
    mode = args.mode
    out = args.out
    #testsuite = args.testsuite
    testsuite = "custom"
    bins_list_file = args.bins_list_file
    global REPEAT
    if args.repeat:
        REPEAT = args.repeat

    if mode == "parallel-native":
        dirs_to_test = [
            "tmp/builder/etcd",
            "tmp/builder/go-ethereum",
            "tmp/builder/grpc-go",
            "tmp/builder/kubernetes",
            "tmp/builder/moby",
            "tmp/builder/prometheus",
            "tmp/builder/tidb"
        ]
        for d in dirs_to_test:
            total_num_of_tests, total_duration = benchmark_all_native_parallel([
                os.path.join(PROJ_ROOT_DIR, d)
            ])
            with open("result", 'a+') as f:
                f.write(d+":\n")
                f.write(f"total {total_num_of_tests} tests, total {total_duration} seconds\n")
                f.write(f"average {total_num_of_tests/total_duration:.2f} test / sec\n")
        return

    if args.action == "count-tests":
        cnt = count_tests_from_bins_dir(bins_dir)
        print(f"{bins_dir} contains {cnt} tests")
        return


    if testsuite == "custom" and not bins_dir:
        raise Exception("--dir is required if testsuite is custom")

    if bins_list_file:
        bins = get_bins_from_file(bins_list_file)
        selected_bins = [os.path.join(bins_dir, bin) for bin in bins]
    if out is None:
        raise Exception("--out is required")
    benchmark_custom(out, bins_dir, mode, selected_bins)

def count_tests_from_bins_dir(bins_dir:str) -> int:
    bins = glob(f"{bins_dir}/*")
    tests_cnt = 0
    for bin in bins:
        try:
            tests = get_tests_from_test_bin(bin)
            tests_cnt += len(tests)
        except Exception as err:
            print("ignore error: {err}")
    return tests_cnt

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
    try:
        p = subprocess.run([bin, '-test.list', '.*'], stdout=subprocess.PIPE)
        p.check_returncode()
        out = p.stdout.decode('utf-8')
        funcs = out.splitlines()
        print("found tests", funcs)
        return [BinTest(bin, func) for func in funcs]
    except Exception as err:
        print(err)
        return []

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