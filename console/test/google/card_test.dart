import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/google/card.dart' as google;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('google.Card', () {
    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        google.Card(onPressed: (_) {}),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });

  group('Integrations Card', () {
    testWidgets('renders title and description', (WidgetTester tester) async {
      await tester.pumpApp(google.Card(onPressed: (_) {}));
      await tester.pumpAndSettle();

      expect(find.text('Google'), findsOneWidget);
      expect(find.text('Connect Google services'), findsOneWidget);
    });

    testWidgets('tapping card calls onPressed with widget', (
      WidgetTester tester,
    ) async {
      Widget? received;
      await tester.pumpApp(google.Card(onPressed: (w) => received = w));
      await tester.pumpAndSettle();
      await tester.tap(find.text('Google'));
      await tester.pumpAndSettle();

      expect(received, isNotNull);
    });
  });
}
