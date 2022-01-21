#!/usr/bin/env python3
import os
import argparse
from datetime import datetime
from typing import Callable, List, Tuple
import argparse
import matplotlib.pyplot as plt
from matplotlib.pyplot import MultipleLocator
import random
import shutil
import zipfile
import csv

# Constants
GFUZZ_LOG_FILE = "fuzzer.log"
GFUZZ_EXEC_FOLDER = "exec"
GFUZZ_EXEC_STDOUT = "stdout"

class EntryStat:
    def __init__(self, entry, num_of_runs, total_duration) -> None:
        self.entry = entry
        self.num_of_runs = num_of_runs
        self.total_duration = total_duration

class ExecStat:
    def __init__(self, exec_id, start):
        self.exec_id = exec_id
        self.start = start
        self.duration = None
        self.bugs = []
        self.timeout = False



def analyze_gfuzz_output_dir(output_dir):
    analyze_gfuzz_log(output_dir)

def find_next_in_same_worker(logs:List[str], from_idx:int, worker_id:str, if_fn: Callable[[str], bool]) -> str:
    for i in range(from_idx, len(logs)):
        line = logs[i]
        parts = line.split(" ")
        curr_wid = parts[3].rstrip("]")
        if curr_wid == worker_id:
            if if_fn(line):
                return line
    return None

def get_log_worker_id(log:str) -> str:
    parts = log.split(" ")
    return parts[3].rstrip("]")

def get_log_time(log:str) -> datetime:
    if log.startswith("20"):
        parts = log.split(" ")
        time_str = parts[0] + " " + parts[1]
        return datetime.strptime(time_str, '%Y/%m/%d %H:%M:%S')
    return None

def get_logs(log_file:str) -> List[str]:
    with open(log_file) as log_f:
        log_lines = log_f.read().splitlines()
    return log_lines

def get_exec_stats_from_logs(logs:List[str]) -> List[ExecStat]:
    exec_stats = {}
    workers_curr_exec = {}
    for _, line in enumerate(logs):
        parts = line.split(" ")
        curr_t = get_log_time(line)
        if curr_t:
            if line.find("] received ") != -1:
                exec_id = parts[5]
                exec_stats[exec_id] = ExecStat(exec_id, curr_t)
                worker_id = get_log_worker_id(line)
                workers_curr_exec[worker_id] = exec_id
            elif line.find("found unique bug: ") != -1:
                bug = parts[-1]
                worker_id = get_log_worker_id(line)
                exec_id = workers_curr_exec[worker_id]
                exec_stats[exec_id].bugs.append(bug)
            elif line.find("] finished ") != -1:
                exec_id = parts[5]
                exec_stats[exec_id].duration = (curr_t - exec_stats[exec_id].start).total_seconds()
            elif line.find("ignore because of timeout") != -1:
                exec_id = parts[4][:-1] # remove semicolon at the end of exec id
                exec_stats[exec_id].timeout = True
    return exec_stats.values()


def analyze_gfuzz_log(output_dir):
    log_start_time = None
    log_end_time = None
    logs = get_logs(os.path.join(output_dir, GFUZZ_LOG_FILE))
    
    for line in logs:
        try:
            parts = line.split(" ")
            if line.startswith("20"):
                time_str = parts[0] + " " + parts[1]
                cur_time = datetime.strptime(time_str, '%Y/%m/%d %H:%M:%S')

                if log_start_time is None:
                    log_start_time = cur_time
                
                log_end_time = cur_time
        except Exception as ex:
            print(f"failed to parse line {line}: {ex}")
    
    exec_stats_arr = get_exec_stats_from_logs(logs)
    entry_stats_arr = analyze_exec_stats_arr(exec_stats_arr)
    num_of_runs_without_timeout = 0
    total_dur_without_timeout = 0
    for e in exec_stats_arr:
        if not e.timeout and e.duration:
            num_of_runs_without_timeout += 1
            total_dur_without_timeout += e.duration

    num_of_runs = len(exec_stats_arr)
    total_dur = (log_end_time - log_start_time).total_seconds()
    print(f"""
total entries: {len(entry_stats_arr)}
total runs: {num_of_runs}
total duration (sec): {total_dur}
average (run/sec): {num_of_runs/total_dur:.2f}

total runs (without timeout): {num_of_runs_without_timeout}
total duration (without timeout): {total_dur_without_timeout}
average (run/sec): {num_of_runs_without_timeout/total_dur_without_timeout:.2f}
    """)

    print("bug statistics:")
    print("used (hours), buggy primitive location, gfuzz exec")
    start_str = log_start_time.strftime("%Y/%m/%d %H:%M:%S")
    for i, e in enumerate(exec_stats_arr):
        skip = False
        for b in e.bugs:
            if b.find("testmain.go") != -1:
                skip = True 
        if skip:
            continue
        t = e.start.strftime("%Y/%m/%d %H:%M:%S")
        dur = (e.start - log_start_time).total_seconds() / 3600
        exec_base_dir = os.path.join(output_dir, GFUZZ_EXEC_FOLDER, e.exec_id)
        stdout_file = os.path.join(exec_base_dir, GFUZZ_EXEC_STDOUT)
        with open(stdout_file) as stdout_f:
            s = stdout_f.read()
            if bug_filter_cond(s):
                continue
        for b in e.bugs:
            print(f"{dur:.2f} h,{b},{e.exec_id}")

    # Most time-consuming entries
    # print("most time-consuming entries:")
    # for e in top_n_time_consuming_entries(entry_stats_arr, 10):
    #     print(f"{e.entry}, {e.num_of_runs} runs, {e.total_duration} seconds\n")

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
        with open(log_fp) as log_f:
            log_lines = log_f.readlines()
            times, nums = get_times_found_bug_nums(log_lines)
        times_arr.append(times)
        nums_arr.append(nums)
        if output_dir.endswith("/"):
            output_dir = output_dir[:-1]
        legends.append(os.path.basename(output_dir))
    
    x_major_locator=MultipleLocator(1)

    plt.figure()
    ax = plt.subplot()

    for i in range(len(legends)):
        ax.plot(times_arr[i], nums_arr[i])
    
    plt.title("GFuzz", fontsize=20)
    plt.xlabel("Time (h)", fontsize=20)
    plt.ylabel("Num of Unique Bugs", fontsize=20)
    plt.xticks(fontsize=20)
    plt.yticks(fontsize=20)
    leg = plt.legend(legends, fontsize=14, handlelength=3)
    plt.xlim([0, 12])
    ax.xaxis.set_major_locator(x_major_locator)

    plt.grid()

    plt.tight_layout()
    plt.savefig(graph_fp, dpi = 200)
    

def get_times_found_bug_nums(log_lines:List[str])->Tuple[List[int], List[int]]:
    start_time = None
    num_of_unique_bug = 0
    times = []
    nums = []
    prev_duration = -0.2
    for idx, line in enumerate(log_lines):
        parts = line.split(" ")
        if line.startswith("20"):
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

            if line.find("found unique bug:") != -1 and line.find("_testmain.go") == -1:
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



def filter_bug(gfuzz_out_dir:str, bug_filter:Callable[[str], bool]) -> List[ExecStat]:
    """To test bug filter condition will filter out which bugs

    Args:
        gfuzz_out_dir (str): GFuzz's output directory
    """
    log_file = os.path.join(gfuzz_out_dir, GFUZZ_LOG_FILE)
    logs = get_logs(log_file)
    stats = get_exec_stats_from_logs(logs)

    filtered_execs = []

    for e in stats:
        if len(e.bugs) == 0:
            continue
        exec_base_dir = os.path.join(gfuzz_out_dir, GFUZZ_EXEC_FOLDER, e.exec_id)
        stdout_file = os.path.join(exec_base_dir, GFUZZ_EXEC_STDOUT)
        with open(stdout_file) as stdout_f:
            s = stdout_f.read()
            if not bug_filter_cond(s):
                filtered_execs.append(e)

    return filtered_execs


def bug_filter_cond(stdout: str) -> bool:
    return stdout.find("[IO wait]") != -1 or stdout.find("[syscall]") != -1 or stdout.find("[runnable]") != -1


def analyze_exec_with_goleak(gfuzz_out_dir:str):
    log_file = os.path.join(gfuzz_out_dir, GFUZZ_LOG_FILE)
    logs = get_logs(log_file)
    stats = get_exec_stats_from_logs(logs)
    buggy_entries_cnt = {
        # entry -> cnt
    }
    buggy_entries_time = {
        # entry -> first buggy time
    }
    
    for e in stats:
        exec_base_dir = os.path.join(gfuzz_out_dir, GFUZZ_EXEC_FOLDER, e.exec_id)
        stdout_file = os.path.join(exec_base_dir, GFUZZ_EXEC_STDOUT)
        entry = get_entry_from_exec_id(e.exec_id)
        with open(stdout_file) as stdout_f:
            s = stdout_f.read()
            if s.find("Leaked goroutine:") != -1:
                if entry not in buggy_entries_cnt:
                    buggy_entries_cnt[entry] = 0
                buggy_entries_cnt[entry] += 1

                if entry not in buggy_entries_time:
                    buggy_entries_time[entry] = e.start
    
    for e, c in buggy_entries_cnt.items(): 
        t = buggy_entries_time[e].strftime("%Y/%m/%d %H:%M:%S")
        print(f"{t} {e} {c}")

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--bug-analyze', default=False, action='store_true')

    # Bug Time Graph
    parser.add_argument('--btg', type=str)

    parser.add_argument('--gfuzz-out-dir', type=str, nargs='*')

    # Copy only buggy exec to given folder
    parser.add_argument('--buggy-dir', type=str)

    # Zipping only buggy exec to zip file
    parser.add_argument('--buggy-zip', type=str)

    parser.add_argument('--test-bug-filter', default=False, action='store_true')

    parser.add_argument('--analyze-goleak', default=False, action='store_true')

    parser.add_argument('--graph-from-csv', type=str)

    args = parser.parse_args()

    

    if args.bug_analyze is not None:
        if len(args.gfuzz_out_dir) != 1:
            return print("expect --gfuzz-out-dir has exact one argument")
        analyze_gfuzz_log(args.gfuzz_out_dir[0])
        return
    
    if args.test_bug_filter:
        if len(args.gfuzz_out_dir) != 1:
            return print("expect --gfuzz-out-dir has exact one argument")
        filter_bug(args.gfuzz_out_dir[0], bug_filter_cond)
        return
    
    if args.analyze_goleak:
        if len(args.gfuzz_out_dir) != 1:
            return print("expect --gfuzz-out-dir has exact one argument")
        analyze_exec_with_goleak(args.gfuzz_out_dir[0])
        return
    
    if args.btg is not None:
        if args.graph_from_csv is not None:
            draw_btg_from_csv(args.graph_from_csv, args.btg)
            return
        if len(args.gfuzz_out_dir) == 0:
            return print("expect --gfuzz-out-dir has at least one argument")
        generate_bug_time_graph(args.gfuzz_out_dir, args.btg)
        return

    if args.buggy_dir is not None:
        if len(args.gfuzz_out_dir) != 1:
            return print("expect --gfuzz-out-dir has exact one argument")
        extract_buggy_to_dir(args.gfuzz_out_dir[0], args.buggy_dir)
        return

    if args.buggy_zip is not None:
        if len(args.gfuzz_out_dir) != 1:
            return print("expect --gfuzz-out-dir has exact one argument")
        extract_buggy_to_zip(args.gfuzz_out_dir[0], args.buggy_zip)
        return


def extract_buggy_to_dir(gfuzz_out_dir:str, buggy_dst_dir:str):
    log_file = os.path.join(gfuzz_out_dir, GFUZZ_LOG_FILE)
    logs = get_logs(log_file)
    exec_stats = get_exec_stats_from_logs(logs)
    buggy_execs = []
    for e in exec_stats:
        if len(e.bugs) > 0:
            buggy_execs.append(e.exec_id)
    
    os.makedirs(os.path.join(buggy_dst_dir, GFUZZ_EXEC_FOLDER), exist_ok=True)
    for be in buggy_execs:
        src_exec_dir = os.path.join(gfuzz_out_dir, GFUZZ_EXEC_FOLDER, be)
        dst_exec_dir = os.path.join(buggy_dst_dir, GFUZZ_EXEC_FOLDER, be)
        if not os.path.exists(dst_exec_dir):
            shutil.copytree(src_exec_dir, dst_exec_dir)
        else:
            print(f"skip {src_exec_dir} since destination already has one\n")
    
    dst_log_file = os.path.join(buggy_dst_dir, GFUZZ_LOG_FILE)
    if os.path.exists(dst_log_file):
        os.remove(dst_log_file)
    shutil.copy(log_file, dst_log_file)
    
    print(f"{len(buggy_execs)} buggy execs and fuzzer.log are copied")


def extract_buggy_to_zip(gfuzz_out_dir:str, dst_zip:str):
    zipf = zipfile.ZipFile(dst_zip, "w")    
    log_file = os.path.join(gfuzz_out_dir, GFUZZ_LOG_FILE)
    logs = get_logs(log_file)
    exec_stats = get_exec_stats_from_logs(logs)
    buggy_execs = []
    for e in exec_stats:
        if len(e.bugs) > 0:
            buggy_execs.append(e.exec_id)
    
    base_dir = os.path.basename(gfuzz_out_dir)
    for be in buggy_execs:
        src_exec_dir = os.path.join(gfuzz_out_dir, GFUZZ_EXEC_FOLDER, be)
        for _, _, children in os.walk(src_exec_dir):
            for child in children:
                src = os.path.join(src_exec_dir, child)
                dst = os.path.join(base_dir, GFUZZ_EXEC_FOLDER, be, child)
                zipf.write(src, dst)

    zipf.write(log_file, os.path.join(base_dir, GFUZZ_LOG_FILE))
    
    zipf.close()
    print(f"{len(buggy_execs)} buggy execs and fuzzer.log are zipped")

def smooth_bugs_durs(bugs, durs):
    oi = 0
    new_bugs = [0]
    new_durs = [0]
    cumulative_dur = 0
    for i in range(0, 1202, 2):
        cumulative_dur += i/100
        while oi < len(bugs) and durs[oi] < cumulative_dur:
            new_bugs.append(bugs[oi])
            new_durs.append(durs[oi])
            oi += 1
        
        new_bugs.append(new_bugs[-1])
        new_durs.append(cumulative_dur)
    
    while oi < len(bugs):
        new_bugs.append(bugs[oi])
        new_durs.append(durs[oi])
        oi += 1

    return (new_bugs[1:], new_durs[1:])

def draw_btg_from_csv(csv_file:str, out:str):
    """Draw bug-time-graph from csv file (google sheet)

    Args:
        csv_file (str): CSV file path
        out (str): Output png file path
    """
    fb_bugs = []
    fb_durs = []

    nfb_bugs = []
    nfb_durs = []

    nose_bugs = []
    nose_durs = []

    nooracle_bug = []
    nooracle_durs = []
    with open(csv_file) as csv_f:
        csv_reader = csv.DictReader(csv_f)
        for num, line in enumerate(csv_reader):
            if 1 <= num <= 27:
                if line["Bug"] == "1":
                    fb_bugs.append(len(fb_bugs)+1)
                    fb_durs.append(float(line["time spent"]))
            elif 30 <= num <= 44:
                if line["Bug"] == "1":
                    nfb_bugs.append(len(nfb_bugs)+1)
                    nfb_durs.append(float(line["time spent"]))
            elif 47 <= num <= 59:
                if line["Bug"] == "1":
                    nose_bugs.append(len(nose_bugs)+1)
                    nose_durs.append(float(line["time spent"]))
            elif 67 <= num <= 69:
                if line["Bug"] == "1":
                    nooracle_bug.append(len(nooracle_bug)+1)
                    nooracle_durs.append(float(line["time spent"]))

    
    plt.figure()
    ax = plt.subplot()

    smoothed_fb_bugs, smoothed_fb_durs = smooth_bugs_durs(fb_bugs, fb_durs)
    smoothed_nfb_bugs, smoothed_nfb_durs = smooth_bugs_durs(nfb_bugs, nfb_durs)
    smoothed_nose_bugs, smoothed_nose_durs = smooth_bugs_durs(nose_bugs, nose_durs)
    smoothed_nooracle_bugs, smoothed_nooracle_durs = smooth_bugs_durs(nooracle_bug, nooracle_durs)

    ax.plot(smoothed_fb_durs, smoothed_fb_bugs, linestyle = 'dashed', marker = 'p', markevery=6, linewidth=2.0, markersize=7)
    ax.plot(smoothed_nfb_durs, smoothed_nfb_bugs, linestyle = (0, (3, 1, 1, 1, 1, 1)), marker = 'o', markevery=6, linewidth=2.0, markersize=7)
    ax.plot(smoothed_nose_durs, smoothed_nose_bugs, linestyle = 'dashdot', marker = 's', markevery=6, linewidth=2.0, markersize=7)
    ax.plot(smoothed_nooracle_durs, smoothed_nooracle_bugs, linestyle = 'dotted', marker = '*', markevery=6, linewidth=2.0, markersize=7)

    plt.xlabel("Time (h)", fontsize=20)
    plt.ylabel("Num of Unique Bugs", fontsize=20)
    plt.xticks(fontsize=20)
    plt.yticks(fontsize=20)
    plt.grid()
    plt.xlim([0, 12])
    plt.legend(["GFuzz","no_feedbacks", "no_mutations", "no_oracle"], fontsize=14, handlelength=2, bbox_to_anchor=(0.57,0.4))
    plt.tight_layout()
    plt.savefig(out, dpi = 200)


if __name__ == "__main__":
    main()