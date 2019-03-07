/**
 * @author Huahang Liu (huahang@xiaomi.com)
 * @date 2019-03-05
 */

#ifdef USE_ASIO
#include <boost/asio.hpp>
#include <boost/asio/thread_pool.hpp>
#elif defined(USE_TBB)
#include <tbb/task.h>
#include <tbb/task_group.h>
#include <tbb/task_scheduler_init.h>
#endif

#include <boost/filesystem.hpp>

#include <algorithm>
#include <cstdlib>
#include <cstring>
#include <deque>
#include <functional>
#include <iostream>
#include <mutex>
#include <string>
#include <thread>

void shrinkTask(const std::string& filename) {
  static std::mutex consoleMutex;
  {
    std::lock_guard<std::mutex> lock(consoleMutex);
    std::cout << "\33[2K\r";
    std::cout << "[Processing] " << filename << std::flush;
  }
  std::string command = "convert -resize 4000x4000 -quality 90 ";
  command += filename;
  command += " ";
  command += filename;
  int rc = std::system(command.c_str());
  if (rc != 0) {
    std::lock_guard<std::mutex> lock(consoleMutex);
    std::cerr << "Command failed: " << command << std::endl;
  }
}

int main(int argc, char** argv) {
  using boost::filesystem::directory_iterator;
  using boost::filesystem::path;
  if (argc != 2) {
    std::cerr << "Usage: imageshrink [path to scan]" << std::endl;
    return EXIT_FAILURE;
  }
  std::deque<path> queue;
  queue.emplace_back(path(argv[1]));
#ifdef USE_ASIO
  boost::asio::thread_pool threadPool(std::thread::hardware_concurrency());
#elif defined(USE_TBB)
  tbb::task_scheduler_init init(std::thread::hardware_concurrency());
  tbb::task_group threadPool;
#endif
  while (!queue.empty()) {
    path p(queue.front());
    queue.pop_front();
    directory_iterator end_itr;
    for (directory_iterator itr(p); itr != end_itr; ++itr) {
      if (is_symlink(itr->path())) {
        continue;
      }
      if (is_directory(itr->path())) {
        queue.push_back(itr->path());
        continue;
      }
      if (!is_regular_file(itr->path())) {
        continue;
      }
      std::string filename = itr->path().string();
      std::string extension = itr->path().extension().string();
      std::transform(
        extension.begin(),
        extension.end(),
        extension.begin(),
        tolower
      );
      if (extension != ".jpg" && extension != ".jpeg") {
        continue;
      }
      auto task = [=] { shrinkTask(filename); };
#ifdef USE_ASIO
      boost::asio::post(threadPool, task);
#elif defined(USE_TBB)
      threadPool.run(task);
#endif
    }
  }
#ifdef USE_ASIO
  threadPool.join();
#elif defined(USE_TBB)
  threadPool.wait();
#endif
  std::cout << std::endl;
  return EXIT_SUCCESS;
}
