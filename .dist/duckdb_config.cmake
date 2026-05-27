# 1. Statically load the requested extensions
duckdb_extension_load(autocomplete)
duckdb_extension_load(json)
duckdb_extension_load(parquet)
duckdb_extension_load(icu)
duckdb_extension_load(inet)
duckdb_extension_load(fts)

# 2. Equivalent CMake variables (set globally with CACHE FORCE to override defaults)
set(BUILD_UNITTESTS 0 CACHE BOOL "Disable unit tests" FORCE)
set(BUILD_SHELL 0 CACHE BOOL "Disable CLI shell" FORCE)
set(ENABLE_EXTENSION_AUTOLOADING 1 CACHE BOOL "Enable extension autoloading" FORCE)
set(ENABLE_EXTENSION_AUTOINSTALL 1 CACHE BOOL "Enable extension autoinstall" FORCE)
