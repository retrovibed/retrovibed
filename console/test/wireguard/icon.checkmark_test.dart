import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/wireguard/icon.checkmark.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('IconCheckmark', () {
    testWidgets('renders without overflow', (WidgetTester tester) async {
      await tester.pumpApp(IconCheckmark(true));
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
      expect(find.byType(IconCheckmark), findsOneWidget);
    });

    testWidgets('shows full opacity when checked', (WidgetTester tester) async {
      await tester.pumpApp(IconCheckmark(true));
      await tester.pumpAndSettle();

      final opacity = tester.widget<Opacity>(find.byType(Opacity));
      expect(opacity.opacity, 1.0);
    });

    testWidgets('shows reduced opacity when unchecked', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(IconCheckmark(false));
      await tester.pumpAndSettle();

      final opacity = tester.widget<Opacity>(find.byType(Opacity));
      expect(opacity.opacity, 0.1);
    });

    testWidgets('calls onTap when tapped', (WidgetTester tester) async {
      var tapped = false;
      await tester.pumpApp(
        IconCheckmark(false, onTap: () async => tapped = true),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(IconButton));
      await tester.pump();

      expect(tapped, isTrue);
    });

    testWidgets('does not call through when no onTap provided', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(IconCheckmark(false));
      await tester.pumpAndSettle();

      await tester.tap(find.byType(IconButton));
      await tester.pump();

      expect(tester.takeException(), isNull);
    });
  });
}
