cmake_minimum_required(VERSION 3.0)
project(imageshrink)

set(CMAKE_CXX_STANDARD 17)

find_package(Boost REQUIRED COMPONENTS filesystem system)

include_directories(${Boost_INCLUDE_DIRS})

if (${CMAKE_SYSTEM_NAME} MATCHES "Darwin")
  add_definitions(-DUSE_ASIO)
endif ()

if (${CMAKE_SYSTEM_NAME} MATCHES "Linux")
  add_definitions(-DUSE_TBB)
endif ()

add_executable(imageshrink main.cpp)
target_link_libraries(imageshrink ${Boost_LIBRARIES})
if (${CMAKE_SYSTEM_NAME} MATCHES "Linux")
  target_link_libraries(imageshrink tbb)
endif ()

install(TARGETS imageshrink RUNTIME DESTINATION bin)
