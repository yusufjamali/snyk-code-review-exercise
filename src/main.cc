#include "npm.h"
#include "util.h"
#include "model.h"

#include <iostream>
#include <string>
#include <vector>

int main(int argc, char* argv[]) {
  if(argc < 3) {
    throw std::runtime_error("Error: Package name and version must be provided.");
  }

  model::VersionedPackage package = npm::get_versioned_package(argv[1], argv[2]);

  nlohmann::json j;
  model::to_json(j, package);
  std::cout << j.dump(2) << std::endl;

  return 0;
}