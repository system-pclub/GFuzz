package metrics

// "github.com/edsrzf/mmap-go" TODO: use mmap to increase performance

// type Metrics struct {
// 	// Map from bug ID to stdout file
// 	Bugs                map[string]*BugMetrics
// 	NumOfBugsFound      uint64
// 	NumOfRuns           uint64
// 	NumOfFuzzQueryEntry uint64
// 	// How many test cases/binary need to be fuzzed
// 	NumOfTotalTargets uint64
// 	// How many test cases/binary triggered
// 	NumOfExecutedTargets uint64
// 	// When are they reach different stages
// 	ExecutedTargets map[string]*TargetMetrics
// 	TimeoutTargets  map[string]uint32
// 	StartAt         time.Time
// 	// Seconds
// 	Duration uint64
// }

// type BugMetrics struct {
// 	FoundAt time.Time
// 	Stdout  string
// }

// type TargetMetrics struct {
// 	At         map[fuzz.FuzzStage]time.Time
// 	MaxCaseCov float32
// }

// func StreamMetrics(filePath string, interval time.Duration) {
// 	go func() {
// 		defer func() {
// 			if err := recover(); err != nil {
// 				log.Printf("StreamMetrics: recovered from panic: %s", err)
// 			}
// 		}()

// 		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
// 		if err != nil {
// 			log.Fatalf("failed to open metrics file: %v", err)
// 		}
// 		log.Printf("metrics file: %s", filePath)
// 		ticker := time.NewTicker(interval * time.Second)

// 		for {
// 			<-ticker.C
// 			b, err := GetFuzzerMetricsJsonBytes(fuzzerContext)
// 			if err != nil {
// 				log.Printf("failed to serialize metrics: %v", err)
// 				continue
// 			}
// 			if err := f.Truncate(0); err != nil {
// 				log.Printf("failed to truncate file: %v", err)
// 				continue
// 			}
// 			if _, err := f.Seek(0, 0); err != nil {
// 				log.Printf("failed to seek file: %v", err)
// 				continue
// 			}
// 			n, err := f.Write(b)
// 			if err != nil {
// 				log.Printf("failed to write to file: %v", err)
// 				continue
// 			}
// 			if n != len(b) {
// 				log.Printf("failed to write all metrics to file, epected %d, actial: %d", len(b), n)
// 				continue
// 			}
// 		}

// 	}()

// }
