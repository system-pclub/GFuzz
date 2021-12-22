#!/usr/bin/env python3
import os
import argparse
from datetime import datetime

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