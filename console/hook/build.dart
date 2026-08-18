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

    final files = dir.listSync().whereType<File>().where((file) => file.path.endsWith(ext));

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
