import 'dart:io';
import 'package:code_assets/code_assets.dart';
import 'package:hooks/hooks.dart';
import 'package:path/path.dart' as path;

void main(List<String> args) async {
  await build(args, (input, output) async {
    // Nothing to emit when this invocation isn't building code assets:
    // adding a CodeAsset anyway makes the validator reject it as an
    // unsupported asset type.
    if (!input.config.buildCodeAssets) {
      return;
    }

    final pattern =
        Platform.environment['NIX_RETROVIBED_SHARED_NATIVE_LIBS'] ??
        path.normalize(path.absolute('../.eg.cache/dev.native.libs/example.so'));

    final ext = path.extension(pattern);
    if (ext.isEmpty) {
      throw Exception('NIX_RETROVIBED_SHARED_NATIVE_LIBS must be of the form \'{dir}/.{ext}\'. is \'${pattern}\'.');
    }

    final dirPath = path.dirname(pattern);
    final dir = Directory(dirPath);

    // No directory means nothing to bundle, not an error: a platform whose
    // native libs are delivered some other way (e.g. android's gradle
    // jniLibs packaging, driven by the committed
    // console/android/app/src/main/jniLibs/{abi} symlinks — see
    // android.JNILibDir in .eg/android/android.go) simply never gets
    // pointed at a populated directory here. This hook only handles
    // platforms that actually need Flutter's own native-assets bundling;
    // it has no platform knowledge of its own.
    if (!dir.existsSync()) {
      return;
    }

    // Recurse and prefer whichever candidates' path names this
    // invocation's targetArchitecture (the build hook is invoked once per
    // target architecture), so a multi-arch directory only contributes its
    // own arch's file. Single-arch layouts (libs sitting flat, no
    // arch-named path segment) never match and the plain fallback (every
    // matching file) applies unchanged.
    final candidates = dir
        .listSync(recursive: true)
        .whereType<File>()
        .where((file) => file.path.endsWith(ext))
        .toList();

    // Match architecture as a whole path segment, not a substring: "arm" is
    // a substring of "arm64", so a plain .contains() would wrongly match
    // an arm64 build for the (32-bit) Architecture.arm invocation too.
    final architecture = input.config.code.targetArchitecture.name;
    final archMatches = candidates.where((file) => path.split(file.path).contains(architecture)).toList();
    final files = archMatches.isNotEmpty ? archMatches : candidates;

    for (final libFile in files) {
      final filename = path.basename(libFile.path);

      output.assets.code.add(
        CodeAsset(
          package: input.packageName,
          name: "${architecture}.${filename}",
          linkMode: DynamicLoadingBundled(),
          file: libFile.uri,
        ),
      );

      output.dependencies.add(libFile.uri);
    }

    output.dependencies.add(dir.uri);
  });
}
