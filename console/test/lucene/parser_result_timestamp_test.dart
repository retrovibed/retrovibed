import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/lucene.dart' as lucene;
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/timex.dart' as timex;

void main() {
  group('ParserResultTimestamp inside chip wrapper', () {
    testWidgets('does not throw during layout', (tester) async {
      final field = lucene.Timestamp.auto('published', timex.epoch, (_) {});
      final result = field.of(DateTime.utc(2025, 6, 1));

      await tester.pumpApp(
        Wrap(children: [result]),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('renders field name and date text', (tester) async {
      final field = lucene.Timestamp.auto('published', timex.epoch, (_) {});
      final result = field.of(DateTime.utc(2025, 6, 1));

      await tester.pumpApp(
        Wrap(children: [result]),
      );
      await tester.pumpAndSettle();

      expect(find.textContaining('published'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
