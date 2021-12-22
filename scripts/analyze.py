#!/usr/bin/env python3
import os
import argparse
from datetime import datetime
from typing import List, Tuple
import path
import argparse
import datetime
import matplotlib.pyplot as plt
from matplotlib.pyplot import MultipleLocator
import random

class EntryStat:
    def __init__(self, entry, num_of_runs, total_duration) -> None:
        self.entry = entry
        self.num_of_runs = num_of_runs
        self.total_duration = total_duration

class ExecStat:
    def __init__(self, exec_id, start, duration):
        self.exec_id = exec_id
        self.start = start
        self.duration = duration

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--log', type=str)
    parser.add_argument('--bug-time-graph', type=str)
    parser.add_argument('--gfuzz-out-dir', type=str, nargs='*')
    args = parser.parse_args()

    if args.log is not None:
        analyze_gfuzz_log(args.log)
    
    if args.bug_time_graph is not None:
        generate_bug_time_graph(args.gfuzz_out_dir, args.bug_time_graph)


def analyze_gfuzz_output_dir(output_dir):
    log_fp = os.path.join(output_dir, "fuzzer.log")
    analyze_gfuzz_log(log_fp)

def analyze_gfuzz_log(log_fp):
    exec_stats = {}
    with open(log_fp, "r") as log_f:
        for line in log_f:
            try:
                parts = line.split(" ")
                if line.startswith("2021"):
                    time_str = parts[0] + " " + parts[1]
                    cur_time = datetime.strptime(time_str, '%Y/%m/%d %H:%M:%S')

                    if line.find("] received ") != -1:
                        exec_id = parts[5]
                        exec_stats[exec_id] = ExecStat(exec_id, cur_time, None)
                    elif line.find("] finished ") != -1:
                        exec_id = parts[5]
                        exec_stats[exec_id].duration = (cur_time - exec_stats[exec_id].start).total_seconds()
            except Exception as ex:
                print(f"failed to parse line {line}: {ex}")
    exec_stats_arr = exec_stats.values()
    entry_stats_arr = analyze_exec_stats_arr(exec_stats_arr)
    
    print(f"""
total entries: {len(entry_stats_arr)}
total runs: {len(exec_stats_arr)}

    """)

    # Most time-consuming entries
    print("most time-consuming entries:")
    for e in top_n_time_consuming_entries(entry_stats_arr, 100):
        print(f"{e.entry}, {e.num_of_runs} runs, {e.total_duration} seconds\n")

def top_n_time_consuming_entries(entry_stats_arr, top):
    sorted_arr = sorted(entry_stats_arr, key= lambda e:e.total_duration, reverse=True)
    return sorted_arr[:top]

def analyze_exec_stats_arr(exec_stats):
    entry_stats = {}
    for es in exec_stats:
        entry = get_entry_from_exec_id(es.exec_id)
        if es.duration is None:
            print(f"ignore {es.exec_id} since duration is None\n")
            continue
        if entry in entry_stats:
            entry_stats[entry].num_of_runs += 1
            entry_stats[entry].total_duration += es.duration
        else:
            entry_stats[entry] = EntryStat(
                entry,
                1,
                es.duration
            )
    return entry_stats.values()

def random_color():
    r = random.random()
    b = random.random()
    g = random.random()
    color = (r, g, b)
    return color

def generate_bug_time_graph(output_dirs:List[str], graph_fp):
    print(output_dirs, graph_fp)
    times_arr = []
    nums_arr = []
    legends = []
    for output_dir in output_dirs:
        log_fp = os.path.join(output_dir, "fuzzer.log")
        times, nums = get_times_found_bug_nums(log_fp)
        times_arr.append(times)
        nums_arr.append(nums)
        legends.append(os.path.basename(output_dir))
    
    x_major_locator=MultipleLocator(0.5)

    plt.figure()
    ax = plt.subplot()

    for i in range(len(legends)):
        ax.plot(times_arr[i], nums_arr[0], c=random_color())
    
    plt.title("GFuzz", fontsize=20)
    plt.xlabel("Time (h)", fontsize=20)
    plt.ylabel("Num of Unique Bugs", fontsize=20)
    plt.xticks(fontsize=20)
    plt.yticks(fontsize=20)
    leg = plt.legend(legends, fontsize=14, handlelength=3)
    plt.xlim([0,3])
    plt.ylim([0, 20])
    ax.xaxis.set_major_locator(x_major_locator)

    plt.grid()

    plt.tight_layout()
    plt.savefig(graph_fp, dpi = 200)
    

def get_times_found_bug_nums(log_lines:List[str])->Tuple[List[int], List[int]]:
    start_time = None
    num_of_unique_bug = 0
    times = []
    nums = []
    for idx, line in log_lines:
        parts = line.split(" ")
        if line.startswith("2021"):
            time_str = parts[0] + " " + parts[1]
            cur_time = datetime.strptime(time_str, '%Y/%m/%d %H:%M:%S')

            if start_time is None:
                start_time = cur_time

            cur_duration = (cur_time - start_time).total_seconds()
            # To hours
            cur_duration = cur_duration / 3600

            if idx == len(log_lines) - 1:
                times.append(cur_duration)
                if len(nums) > 0:
                    nums.append(nums[-1])
                else:
                    nums.append(0)

            if "found unique bug:" in line:
                num_of_unique_bug += 1
                continue
            # Restrict time frame
            if cur_duration > 12.0:
                times.append(cur_duration)
                if len(nums) > 0:
                    nums.append(nums[-1])
                else:
                    nums.append(0)
                break
            
            if cur_duration - prev_duration >= 0.03:
                times.append(cur_duration)
                nums.append(num_of_unique_bug)
                prev_duration = cur_duration
    return times, nums



def get_entry_from_exec_id(exec_id:str):
    parts = exec_id.split("-")
    filtered = parts[2:-1]
    return '-'.join(filtered)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--log', type=str)
    args = parser.parse_args()

    if args.log != None:
        analyze_gfuzz_log(args.log)

if __name__ == "__main__":
    main()