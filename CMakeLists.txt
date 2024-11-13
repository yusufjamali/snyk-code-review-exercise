cmake_minimum_required(VERSION 3.10)
project(snyk-code-review-exercise)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(curlpp REQUIRED)
find_package(nlohmann_json REQUIRED)
find_package(Python3 COMPONENTS Interpreter Development REQUIRED)
find_package(Boost COMPONENTS python REQUIRED)

set(CMAKE_SHARED_MODULE_PREFIX "")

# Retrieve Poetry virtual environment path
execute_process(
    COMMAND poetry env info --path
    OUTPUT_VARIABLE POETRY_ENV_PATH
    OUTPUT_STRIP_TRAILING_WHITESPACE
    ERROR_QUIET
)

if(NOT POETRY_ENV_PATH)
    message(FATAL_ERROR "Could not retrieve Poetry environment path. Ensure Poetry is installed and the environment is set up.")
endif()

# Set PYTHONPATH to Poetry's site-packages for Python dependencies
set(PYTHONPATH ${POETRY_ENV_PATH}/lib/python${Python3_VERSION_MAJOR}.${Python3_VERSION_MINOR}/site-packages)
message(STATUS "Using PYTHONPATH from Poetry environment: ${PYTHONPATH}")

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/src)
file(GLOB SOURCES ${CMAKE_CURRENT_SOURCE_DIR}/src/*.cc)
list(REMOVE_ITEM SOURCES ${CMAKE_CURRENT_SOURCE_DIR}/src/main.cc)

add_library(my_libs ${SOURCES})

target_link_libraries(my_libs PRIVATE Boost::python curlpp::curlpp nlohmann_json::nlohmann_json Python3::Python)
# Pass PYTHONPATH to targets to ensure they can use Python dependencies from Poetry
target_compile_definitions(my_libs PRIVATE "PYTHONPATH=${PYTHONPATH}")
target_compile_definitions(my_libs PRIVATE "POETRY_ENV_PATH=\"${PYTHONPATH}\"")

add_executable(main src/main.cc)
target_link_libraries(main PRIVATE my_libs  nlohmann_json::nlohmann_json)

# Test framework
include(FetchContent)
FetchContent_Declare(
    googletest
    URL https://github.com/google/googletest/releases/download/v1.15.2/googletest-1.15.2.tar.gz
)
FetchContent_MakeAvailable(googletest)
enable_testing()
add_executable(tests tests/test_versioned_package.cc)
target_link_libraries(tests PRIVATE my_libs gtest_main nlohmann_json::nlohmann_json)
add_test(NAME VersionedPackageTest COMMAND tests)