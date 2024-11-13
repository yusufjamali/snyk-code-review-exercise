#include "model.h"
#include <string>

namespace util {
  std::string curl_request(const std::string& url, std::ostringstream* response_stream);

  std::string max_satisfying(const std::vector<std::string>& versions, const std::string& range);
}