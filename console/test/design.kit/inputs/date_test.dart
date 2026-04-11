import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/inputs/date.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('DateInput', () {
    testWidgets('renders with an initial date', (WidgetTester tester) async {
      final initial = DateTime.utc(2024, 6, 15);

      await tester.pumpApp(
        Scaffold(body: DateInput(value: initial, onChanged: (_) {})),
      );
      await tester.pumpAndSettle();

      expect(find.byType(CalendarDatePicker), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without an initial date', (WidgetTester tester) async {
      await tester.pumpApp(Scaffold(body: DateInput(onChanged: (_) {})));
      await tester.pumpAndSettle();

      expect(find.byType(CalendarDatePicker), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
