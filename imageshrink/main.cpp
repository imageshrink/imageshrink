/**
 * @author Huahang Liu (huahang@xiaomi.com)
 * @date 2019-03-05
 */

#include <tbb/task.h>
#include <tbb/task_group.h>
#include <tbb/task_scheduler_init.h>

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
    std::scoped_lock<std::mutex> lock(consoleMutex);
    std::cout << "\33[2K\r";
    std::cout << "[Processing] " << filename << std::flush;
  }
  std::string command = "convert -resize 4000x4000 -quality 90 ";
  command += filename;
  command += " ";
  command += filename;
  std::system(command.c_str());
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
  tbb::task_scheduler_init init(tbb::task_scheduler_init::automatic);
  tbb::task_group taskGroup;
  while (!queue.empty()) {
    path p(queue.front());
    queue.pop_front();
    directory_iterator end_itr;
    for (directory_iterator itr(p); itr != end_itr; ++itr) {
      if (is_directory(itr->path())) {
        queue.push_back(itr->path());
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
      taskGroup.run([=] { shrinkTask(filename); });
    }
  }
  taskGroup.wait();
  std::cout << std::endl;
  return EXIT_SUCCESS;
}
