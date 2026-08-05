duckdb_extension_load(vss
    GIT_URL https://github.com/duckdb/duckdb-vss
    GIT_TAG b833341c8737fd3f3558c7720cc575ae8fc82598
    APPLY_PATCHES
)

set(BUILD_UNITTESTS 0 CACHE BOOL "Disable unit tests" FORCE)
set(BUILD_SHELL 0 CACHE BOOL "Disable CLI shell" FORCE)
set(ENABLE_EXTENSION_AUTOLOADING 1 CACHE BOOL "Enable extension autoloading" FORCE)
set(ENABLE_EXTENSION_AUTOINSTALL 1 CACHE BOOL "Enable extension autoinstall" FORCE)

message(STATUS "------------------------------------------- ${CMAKE_CURRENT_LIST_FILE} loaded ------------------------------------------")
