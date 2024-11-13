#include "util.h"
#include <boost/python.hpp>

#include <stdexcept>
#include <string>
#include <sstream>

#include <curlpp/cURLpp.hpp>
#include <curlpp/Options.hpp>
#include <curlpp/Easy.hpp>

namespace util {
  using namespace boost::python;

  std::string curl_request(const std::string& url, std::ostringstream* response_stream) {
    curlpp::Cleanup cleanup;
    curlpp::Easy request;

    request.setOpt(new curlpp::options::Url(url));
    request.setOpt(new curlpp::options::WriteStream(response_stream));

    request.perform();
    return response_stream->str();
  }
  
  // semver isn't available in C++ ecosystem. Using python binding with package from poetry.
  std::string max_satisfying(const std::vector<std::string>& versions, const std::string& range) {
    try {
      Py_Initialize();
      object sys = import("sys");
      sys.attr("path").attr("insert")(0, POETRY_ENV_PATH);

      list versions_py;
      for (const auto& version : versions) {
          versions_py.append(version);
      }

      object result = import("semver").attr("max_satisfying")(versions_py, range);
      if (result.is_none()) {
        throw std::runtime_error("No satisfying version found for range: " + range);
      }
      return extract<std::string>(result);
    } catch (const error_already_set&) {
      PyErr_Print();
      throw std::runtime_error("Error calling Python function");
    }
  }
}