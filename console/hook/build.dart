import 'dart:io';
import 'package:code_assets/code_assets.dart';
import 'package:hooks/hooks.dart';
import 'package:path/path.dart' as path;

void main(List<String> args) async {
  await build(args, (input, output) async {
    final pattern =
        Platform.environment['NIX_RETROVIBED_SHARED_NATIVE_LIBS'] ??
        path.normalize(path.absolute('../.eg.cache/dev.native.libs/example.so'));

    final ext = path.extension(pattern);
    if (ext.isEmpty) {
      throw Exception('NIX_RETROVIBED_SHARED_NATIVE_LIBS must be of the form \'{dir}/.{ext}\'. is \'${pattern}\'.');
    }

    final dirPath = path.dirname(pattern);
    final dir = Directory(dirPath);

    if (!dir.existsSync()) {
      throw Exception('Directory does not exist: $dirPath');
    }

    // android's override of NIX_RETROVIBED_SHARED_NATIVE_LIBS points at
    // android.JNILibRoot() (.eg/android/android.go), which nests one
    // subdirectory per architecture since the build hook is invoked once
    // per target architecture during a multi-abi release build. Recurse
    // and prefer whichever candidates' path names this invocation's
    // targetArchitecture, so each per-arch invocation only picks up its
    // own .so instead of every arch's. Desktop dev builds are single-arch
    // with libs sitting flat (no arch-named path segment), so nothing
    // there ever matches and the plain fallback (every matching file)
    // applies unchanged.
    final candidates = dir
        .listSync(recursive: true)
        .whereType<File>()
        .where((file) => file.path.endsWith(ext))
        .toList();

    final architecture = input.config.buildCodeAssets
        ? input.config.code.targetArchitecture.name
        : Architecture.x64.name;
    final archMatches = candidates.where((file) => file.path.contains(architecture)).toList();
    final files = archMatches.isNotEmpty ? archMatches : candidates;

    for (final libFile in files) {
      final filename = path.basename(libFile.path);

      output.assets.code.add(
        CodeAsset(
          package: input.packageName,
          name: filename,
          linkMode: DynamicLoadingBundled(),
          file: libFile.uri,
        ),
      );

      output.dependencies.add(libFile.uri);
    }

    output.dependencies.add(dir.uri);
  });
}
