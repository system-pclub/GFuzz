#!/usr/bin/env python3
import sys

def main():
    b1fp = sys.argv[1]
    b2fp = sys.argv[2]

    with open(b1fp) as b1f:
        logs = b1f.read().splitlines()
        b1rec = parse_benchmark(logs)
    
    with open(b2fp) as b2f:
        logs = b2f.read().splitlines()
        b2rec = parse_benchmark(logs)

    b1cnt = 0
    b1dur = 0
    b2cnt = 0
    b2dur = 0
    for k, v in b1rec.items():
        
        if v < 0.001:
            continue
        if v > 10:
            continue
        if k in b2rec:
            b1cnt += 1
            b2cnt += 1
            b1dur += v
            b2dur += b2rec[k]
    
    print(f"common tests: {b2cnt}")
    print(f"first average {b1dur/b1cnt:0.4f}")
    print(f"second average {b2dur/b2cnt:0.4f}")


def parse_benchmark(logs:str):
    rec = {}
    for log in logs:
        parts = log.split(":")
        id = parts[0]
        if not id.split("->")[1].startswith("Test"):
            # we should ignore any other methods that does not start with Test
            continue
        id = id.split("/")[-1]
        time = float(parts[1].split(" ")[0])
        rec[id] = time
    return rec

if __name__ == "__main__":
    main()