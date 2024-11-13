#pragma once
#include <string>
#include <vector>
#include <utility>
#include <nlohmann/json.hpp>

namespace model{
  class VersionedPackage {
  public:
      VersionedPackage(const std::string& name,
                      const std::string& version,
                      const std::string& description,
                      const std::vector<std::pair<std::string, std::string>>& dependencies)
          : name_(name), version_(version), description_(description), dependencies_(dependencies) {}

      const std::string& get_name() const { return name_; }
      const std::string& get_version() const { return version_; }
      const std::string& get_description() const { return description_; }
      const std::vector<std::pair<std::string, std::string>>& get_dependencies() const { return dependencies_; }

      friend void to_json(nlohmann::json& j, const VersionedPackage& vp);

  private:
      std::string name_;
      std::string version_;
      std::string description_;
      std::vector<std::pair<std::string, std::string>> dependencies_;
  };

  void to_json(nlohmann::json& j, const VersionedPackage& vp);
}