package flatpakmods

import (
	"github.com/egdaemon/eg/runtime/x/wasi/egflatpak"
)

func Libduckdb() egflatpak.Module {
	// Duckdb is a pita to build in general and esp on flatpak due to its extensions.
	// they download them during the build process which is disallowed in flatpak.
	// their src archive doesnt include the build tooling for cmake nor does it include the extensions.
	// as a result we pull the prebuilt binaries for the library.

	// [ 11%] Creating directories for 'inet_extension_fc-populate'
	// [ 22%] Performing download step (git clone) for 'inet_extension_fc-populate'
	// Cloning into 'inet_extension_fc-src'...
	// fatal: unable to access 'https://github.com/duckdb/duckdb-inet/': Could not resolve host: github.com
	// Cloning into 'inet_extension_fc-src'...
	// fatal: unable to access 'https://github.com/duckdb/duckdb-inet/': Could not resolve host: github.com
	// Cloning into 'inet_extension_fc-src'...
	// fatal: unable to access 'https://github.com/duckdb/duckdb-inet/': Could not resolve host: github.com
	// Had to git clone more than once: 3 times.
	// CMake Error at inet_extension_fc-subbuild/inet_extension_fc-populate-prefix/tmp/inet_extension_fc-populate-gitclone.cmake:50 (message):
	//   Failed to clone repository: 'https://github.com/duckdb/duckdb-inet'

	return egflatpak.NewModule("duckdb", "simple", egflatpak.ModuleOptions().Commands(
		"cp -r . /app/lib",
	).Sources(
		egflatpak.SourceTarball(
			"https://github.com/duckdb/duckdb/releases/download/v1.4.1/libduckdb-linux-amd64.zip",
			"ce859962fe96ca952d53571964dafa0168644283a76a2c669b3f73120d710edb",
			egflatpak.SourceOptions().Arch("x86_64").Destination("duckdb.zip")...,
		),
		egflatpak.SourceTarball(
			"https://github.com/duckdb/duckdb/releases/download/v1.4.1/libduckdb-linux-arm64.zip",
			"4aa05a74956b1d57f05a139623943231caeee36286c9f42ff10dc278b6df0b6e",
			egflatpak.SourceOptions().Arch("aarch64").Destination("duckdb.zip")...,
		),
	)...)
}

func Libduckdbsrc() egflatpak.Module {
	return egflatpak.NewModule("duckdb", "cmake", egflatpak.ModuleOptions().ConfigOptions(
		"-DEXTENSION_STATIC_BUILD=1",
		"-DBUILD_EXTENSIONS=autocomplete;json;parquet;icu;inet;fts",
		"-DENABLE_EXTENSION_AUTOLOADING=1",
		"-DENABLE_EXTENSION_AUTOINSTALL=1",
		"-DCMAKE_VERBOSE_MAKEFILE=on",
		"-DBUILD_UNITTESTS=0",
		"-DBUILD_SHELL=0",
		"-DCMAKE_BUILD_TYPE=Release",
	).Sources(
		egflatpak.SourceGit(
			"https://github.com/duckdb/duckdb.git",
			"d1dc88f950d456d72493df452dabdcd13aa413dd", // v1.4.3
		),
	)...)
}

// pulled from: https://github.com/flathub/io.mpv.Mpv/blob/d895bc41c09a17d0bdca40cd57f77340e44fdca5/io.mpv.Mpv.yml
func Libx264() egflatpak.Module {
	return egflatpak.NewModule("libx264", "autotools", egflatpak.ModuleOptions().ConfigOptions(
		"--disable-cli",
		"--enable-shared",
	).Cleanup(
		"/include",
		"/lib/pkgconfig",
		"/share/man",
	).Sources(
		egflatpak.SourceGit(
			"https://github.com/mirror/x264.git",
			"31e19f92f00c7003fa115047ce50978bc98c3a0d",
		),
	)...)
}

// pulled from: https://github.com/flathub/io.mpv.Mpv/blob/d895bc41c09a17d0bdca40cd57f77340e44fdca5/io.mpv.Mpv.yml
func Libx265() egflatpak.Module {
	return egflatpak.NewModule("libx265", "cmake", egflatpak.ModuleOptions().SubDirectory("source").ConfigOptions(
		"-DCMAKE_BUILD_TYPE=Release",
		"-DBUILD_STATIC=0",
		"-DCMAKE_POLICY_VERSION_MINIMUM=3.5",
	).Cleanup(
		"/share/man",
	).Sources(
		egflatpak.SourceGit(
			"https://bitbucket.org/multicoreware/x265_git.git",
			"1d117bed4747758b51bd2c124d738527e30392cb",
			egflatpak.SourceOptions().Tag("4.1")...,
		),
		egflatpak.SourceFile(
			"https://raw.githubusercontent.com/retrovibed/retrovibed/9c126143b6b98586810169f6322a49cd39d9ef03/.patches/libx265-cmake.patch",
			egflatpak.SourceOptions().SHA256("7d8a3e1d4862588f95d9fcd26f3bb78ea9ce01346044326d9b9949b28d6f04e9").Destination("libx265-cmake.patch")...,
		),
		egflatpak.SourceShell(
			egflatpak.SourceOptions().Commands(
				// cmake had deprecated/removed build options that caused cmake to fail this patches them out.
				// fixed up stream but not yet tagged in a release.
				"patch -p1 < libx265-cmake.patch",
			)...,
		),
	)...)
}

// pulled from: https://github.com/flathub/io.mpv.Mpv/blob/d895bc41c09a17d0bdca40cd57f77340e44fdca5/io.mpv.Mpv.yml
func Libass() egflatpak.Module {
	return egflatpak.NewModule("libass", "autotools", egflatpak.ModuleOptions().ConfigOptions(
		"--disable-static",
		"--enable-asm",
		"--enable-harfbuzz",
		"--enable-fontconfig",
	).Sources(
		egflatpak.SourceGit(
			"https://github.com/libass/libass.git",
			"bbb3c7f1570a4a021e52683f3fbdf74fe492ae84",
			egflatpak.SourceOptions().Tag(
				"0.17.4",
			)...,
		),
	)...)
}

// pulled from: https://github.com/flathub/io.mpv.Mpv/blob/d895bc41c09a17d0bdca40cd57f77340e44fdca5/io.mpv.Mpv.yml
func Libbs2b() egflatpak.Module {
	return egflatpak.NewModule("libbs2b", "autotools", egflatpak.ModuleOptions().ConfigOptions(
		"--disable-static",
	).Cleanup(
		"/include",
		"/lib/pkgconfig",
	).Sources(
		egflatpak.SourceTarball(
			"https://downloads.sourceforge.net/sourceforge/bs2b/libbs2b-3.1.0.tar.gz",
			"6aaafd81aae3898ee40148dd1349aab348db9bfae9767d0e66e0b07ddd4b2528",
		),
		egflatpak.SourceShell(
			egflatpak.SourceOptions().Commands(
				"sed -i -e 's/lzma/xz/g' configure.ac",
				"autoreconf -vif",
			)...,
		),
	)...)
}

// pulled from: https://github.com/flathub/io.mpv.Mpv/blob/d895bc41c09a17d0bdca40cd57f77340e44fdca5/io.mpv.Mpv.yml
func Libffmpeg() egflatpak.Module {
	return egflatpak.NewModule("libffmpeg", "autotools", egflatpak.ModuleOptions().ConfigOptions(
		"--env=CFLAGS=\"-Wno-incompatible-pointer-types\"",
		"--disable-debug",
		"--disable-doc",
		"--disable-static",
		"--enable-encoder=png",
		"--enable-gnutls",
		"--enable-gpl",
		"--enable-shared",
		"--enable-version3",
		"--enable-libaom",
		"--enable-libass",
		"--enable-libbs2b",
		"--enable-libdav1d",
		"--enable-libdrm",
		"--enable-libfreetype",
		"--enable-libmp3lame",
		"--enable-libopus",
		"--enable-libjxl",
		"--enable-libtheora",
		"--enable-libv4l2",
		"--enable-libvorbis",
		"--enable-libvpx",
		"--enable-libx264",
		"--enable-libx265",
		"--enable-libwebp",
		"--enable-libxml2",
		"--enable-vulkan",
		// "--enable-libmysofa",
	).Cleanup(
		"/share/ffmpeg/examples",
		"/include",
		"/lib/pkgconfig",
	).Sources(
		egflatpak.SourceGit(
			"https://github.com/FFmpeg/FFmpeg.git",
			"f893221c8d89cb798b829bebe71d55e1a3f242fd",
			egflatpak.SourceOptions().Tag(
				"n7.1.2",
			)...,
		),
	)...)
}

// pulled from: https://github.com/flathub/io.mpv.Mpv/blob/d895bc41c09a17d0bdca40cd57f77340e44fdca5/io.mpv.Mpv.yml
func Libplacebo() egflatpak.Module {
	return egflatpak.NewModule("libplacebo", "meson", egflatpak.ModuleOptions().ConfigOptions(
		"-Dvulkan=enabled",
		"-Dshaderc=enabled",
		"--libdir=/app/lib", // fixes build issue with flatpak 1.4.2 - https://github.com/flatpak/flatpak-builder/commit/8c036e00630e35423c03388aacc06cd00dda74ea
	).Sources(
		egflatpak.SourceGit(
			"https://github.com/haasn/libplacebo.git",
			"3188549fba13bbdf3a5a98de2a38c2e71f04e21e",
			egflatpak.SourceOptions().Tag(
				"v7.351.0",
			)...,
		),
		egflatpak.SourceFile(
			"https://raw.githubusercontent.com/retrovibed/retrovibed/9c126143b6b98586810169f6322a49cd39d9ef03/.patches/libplacebo-vulkan-utils_gen-python.patch",
			egflatpak.SourceOptions().Destination("libplacebo-vulkan-utils_gen-python.patch").SHA256("f183bf98c3bd1591cd93a12d4d0372f77e3009bdc0a19a3e0f9edb2dfee0b657")...,
		),
		egflatpak.SourceShell(
			egflatpak.SourceOptions().Commands(
				// patching libplacebo due to a breaking change in python fixed upstream but not yet tagged in a release.
				"patch -p1 < libplacebo-vulkan-utils_gen-python.patch",
			)...,
		),
	)...)
}

// pulled from: https://github.com/flathub/io.mpv.Mpv/blob/d895bc41c09a17d0bdca40cd57f77340e44fdca5/io.mpv.Mpv.yml
func Libmpv() egflatpak.Module {
	// https://github.com/mpv-player/mpv/releases/tag/v0.39.0
	return egflatpak.NewModule("mpv", "meson", egflatpak.ModuleOptions().ConfigOptions(
		"-Dlibmpv=true",
		// "-Dcdda=enabled",
		// "-Ddvbin=enabled",
		// "-Ddvdnav=enabled",
		// "-Dlibarchive=enabled",
		// "-Dsdl2=enabled",
		"-Dvulkan=enabled",
		"-Dmanpage-build=disabled",
		"-Dbuild-date=false",
		"--libdir=/app/lib", // fixes build issue with flatpak 1.4.2 - https://github.com/flatpak/flatpak-builder/commit/8c036e00630e35423c03388aacc06cd00dda74ea
	).Sources(egflatpak.SourceTarball("https://github.com/mpv-player/mpv/archive/refs/tags/v0.39.0.tar.gz", "2ca92437affb62c2b559b4419ea4785c70d023590500e8a52e95ea3ab4554683"))...)
}
