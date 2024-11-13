# npm dependency cpp

A cpp binary that provides a basic functions for querying the dependency
tree of a [npm](https://npmjs.org) package.

## Prerequisites

- Python 3.10 or later
- gcc
- cmake
- poetry
- conan
- numpy < 2.0.0 (Otherwise `conan install` fails for boost. e.g. `pip3 install "numpy<2.0" --break-system-package`)

## Getting Started

To install dependencies (choose the correct --profile:build under `conan_profiles/`):

```sh
poetry install
conan install . --build=missing --profile:host=conan_profiles/gcc_linux_x86_64 --profile:build=conan_profiles/gcc_linux_x86_64 --output-folder=.build
# Mac:
# conan install . --build=missing --profile:host=conan_profiles/gcc_mac_arm --profile:build=conan_profiles/gcc_mac_arm --output-folder=.build
cmake -DCMAKE_TOOLCHAIN_FILE=.build/conan_toolchain.cmake -DCMAKE_BUILD_TYPE=Release -B .build .
```

Now run the main program with

```sh
cmake --build .build
./.build/main react 16.13.0
```

You can run the tests with:

```sh
cmake --build .build
./.build/tests
```

Occasionally you might want to consider cleaning up:

```sh
rm -R ./.build
```