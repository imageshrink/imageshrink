cmake_minimum_required(VERSION 3.0)
project(imageshrink)

set(CMAKE_CXX_STANDARD 17)

find_package(Boost COMPONENTS REQUIRED filesystem system)

include_directories(${Boost_INCLUDE_DIRS})

add_executable(imageshrink
  main.cpp
)
target_link_libraries(imageshrink ${Boost_LIBRARIES})
target_link_libraries(imageshrink tbb)
