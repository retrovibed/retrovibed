import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/library/known.media.download.list.dart';
import 'package:retrovibed/library/api.dart' as api;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('KnownMediaDownload', () {
    testWidgets('renders without overflow', (tester) async {
      await tester.pumpApp(
        KnownMediaDownloadList(
          children: [api.Known(description: 'Test', summary: 'summary')],
        ),
        isolatecache: true,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders empty children without overflow', (tester) async {
      await tester.pumpApp(KnownMediaDownloadList());
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders with leading widget', (tester) async {
      await tester.pumpApp(
        KnownMediaDownloadList(
          leading: const Text('leading'),
          children: [api.Known(description: 'Test', summary: 'summary')],
        ),
        isolatecache: true,
      );
      await tester.pumpAndSettle();
      expect(find.text('leading'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
