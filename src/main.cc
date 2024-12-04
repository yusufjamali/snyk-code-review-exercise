#include "npm.h"
#include "util.h"
#include "model.h"
#include <crow.h>
#include <iostream>
#include <string>
#include <vector>
#include <nlohmann/json.hpp>

int main() {
    crow::SimpleApp app;

    CROW_ROUTE(app, "/package/<string>/<string>")
        .methods(crow::HTTPMethod::GET)([](const crow::request&, crow::response& res, std::string package_name, std::string package_version) {
            try {
                model::VersionedPackage package = npm::get_versioned_package(package_name, package_version);

                nlohmann::json j;
                model::to_json(j, package);

                res.set_header("Content-Type", "application/json");
                res.write(j.dump(2));
                res.end();
            } catch (const std::exception& e) {
                res.code = 500;
                res.write(std::string("Server Error: ") + e.what());
                res.end();
            }
        });

    app.port(3000).multithreaded().run();

    return 0;
}