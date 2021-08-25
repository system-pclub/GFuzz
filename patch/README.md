`patch` directory contains files expected to be patched to golang environment before fuzzing or any compilation related to fuzzing (compiling target application/test into binary in advance).

|directory|target directory|
| --- | --- |
| runtime | $GOROOT/src/runtime|
| time | $GOROOT/src/time |
| sync | $GOROOT/src/sync |
| reflect | $GOROOT/src/reflect |