import 'dart:io';

import 'package:desktop_drop/desktop_drop.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/mimex.dart' as mimex;

void main() {
  group('FileDropWell.resolve', () {
    test('directory resolves as a directory when accepted', () async {
      final dir = Directory.systemTemp.createTempSync('resolve');
      addTearDown(() => dir.deleteSync(recursive: true));

      // desktop_drop reports directories as files on linux and windows.
      final evt = await ds.FileDropWell.resolve([DropItemFile(dir.path)], mimetypes: mimex.folders);

      expect(evt.files, hasLength(1));
      expect(evt.files.first, isA<DropItemDirectory>());
      expect(evt.files.first.path, dir.path);
      expect(evt.files.first.mimeType, mimex.directory);
    });

    test('folders only well discards files', () async {
      final dir = Directory.systemTemp.createTempSync('resolve');
      addTearDown(() => dir.deleteSync(recursive: true));
      final file = File('${dir.path}/example.torrent')..writeAsStringSync('example');

      final evt = await ds.FileDropWell.resolve(
        [DropItemFile(dir.path), DropItemFile(file.path, name: 'example.torrent')],
        mimetypes: mimex.folders,
      );

      expect(evt.files.map((f) => f.path), [dir.path]);
    });

    test('other wells discard directories and pass files through unfiltered', () async {
      final dir = Directory.systemTemp.createTempSync('resolve');
      addTearDown(() => dir.deleteSync(recursive: true));
      final file = File('${dir.path}/example.txt')..writeAsStringSync('example');

      for (final mimetypes in [
        <String>[],
        [mimex.bittorrent],
      ]) {
        final evt = await ds.FileDropWell.resolve(
          [DropItemFile(dir.path), DropItemFile(file.path, name: 'example.txt')],
          mimetypes: mimetypes,
        );

        expect(evt.files, hasLength(1), reason: '$mimetypes');
        expect(evt.files.first, isA<DropItemFile>(), reason: '$mimetypes');
        expect(evt.files.first.path, file.path, reason: '$mimetypes');
      }
    });
  });
}
