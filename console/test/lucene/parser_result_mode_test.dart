import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/lucene.dart' as lucene;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('ParserResultMode inside chip wrapper', () {
    testWidgets('does not throw during layout', (tester) async {
      final field = lucene.Mode.auto('hd', false, (_) {});
      final result = field.of(true);

      await tester.pumpApp(
        Wrap(children: [result]),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('renders field name text', (tester) async {
      final field = lucene.Mode.auto('hd', false, (_) {});
      final result = field.of(true);

      await tester.pumpApp(
        Wrap(children: [result]),
      );
      await tester.pumpAndSettle();

      expect(find.text('hd'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('ParserResultMode edit', () {
    test('returns null because it cannot be edited, only removed', () {
      final field = lucene.Mode.auto('hd', false, (_) {});
      final result = field.of(true);

      expect(result.edit((_) {}), isNull);
    });
  });
}
