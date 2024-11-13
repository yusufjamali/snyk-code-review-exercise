#include "model.h"
#include "util.h"

#include <string>
#include <sstream>
#include <nlohmann/json.hpp>

namespace npm {
  namespace {
    using model::VersionedPackage;
    constexpr std::string_view NPM_REGISTRY_URL = "https://registry.npmjs.org";
  }

  VersionedPackage get_versioned_package(const std::string& name, const std::string& version) {
    std::string url = std::string(NPM_REGISTRY_URL) + "/" + name + "/" + version;

    std::ostringstream response_stream;
    nlohmann::json npm_package = nlohmann::json::parse(util::curl_request(url, &response_stream));
    response_stream.clear();

    std::vector<std::pair<std::string, std::string>> dependencies;
    for (const auto& dep : npm_package["dependencies"].items()) {
        dependencies.emplace_back(dep.key(), dep.value());
    }

    VersionedPackage package{
        npm_package["name"],
        npm_package["version"],
        npm_package["description"],
        std::move(dependencies)
    };

    return package;
  }
}