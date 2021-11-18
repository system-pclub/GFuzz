# This script will automatically generate the corresponding test binary files(with and without instrumented go runtime) to the ./tmp/builder

docker build -f docker/builder/Dockerfile -t gfuzzbuilder:latest .


docker run --rm -it \
-v $(pwd)/tmp/builder:/builder \
-v $(pwd)/tmp/pkgmod:/go/pkg/mod \
gfuzzbuilder:latest