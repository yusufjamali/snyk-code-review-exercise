#include <gtest/gtest.h>
#include <nlohmann/json.hpp>
#include "model.h"
#include "npm.h"

using nlohmann::json;
using model::VersionedPackage;

TEST(VersionedPackageTest, TestGetVersionedPackage) {
    // Define an instance of VersionedPackage
    VersionedPackage package = npm::get_versioned_package("express", "2.0.0");

    // Serialize the instance to JSON
    json generated_json;
    model::to_json(generated_json, package);

    std::string expected_json_str = R"({
        "dependencies": [
          ["connect", ">= 1.1.0 < 2.0.0"],
          ["mime", ">= 0.0.1"], 
          ["qs", ">= 0.0.6"]
        ],
        "name": "express",
        "version": "2.0.0",
        "description": "Sinatra inspired web development framework"
    })";
    json expected_json = json::parse(expected_json_str);

    ASSERT_EQ(generated_json, expected_json);
}