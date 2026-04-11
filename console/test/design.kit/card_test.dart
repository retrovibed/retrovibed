import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('Card cursor', () {
    testWidgets('shows click cursor when onTap is set', (tester) async {
      await tester.pumpApp(
        ds.Card(const Text('content'), onTap: () {}),
      );
      await tester.pumpAndSettle();

      expect(
        tester.resolvedCursorAt(find.byType(ds.Card)),
        SystemMouseCursors.click,
      );
    });

    testWidgets('shows click cursor when onDoubleTap is set', (tester) async {
      await tester.pumpApp(
        ds.Card(const Text('content'), onDoubleTap: () {}),
      );
      await tester.pumpAndSettle();

      expect(
        tester.resolvedCursorAt(find.byType(ds.Card)),
        SystemMouseCursors.click,
      );
    });

    testWidgets('shows click cursor when onLongPress is set', (tester) async {
      await tester.pumpApp(
        ds.Card(const Text('content'), onLongPress: () {}),
      );
      await tester.pumpAndSettle();

      expect(
        tester.resolvedCursorAt(find.byType(ds.Card)),
        SystemMouseCursors.click,
      );
    });

    testWidgets('shows basic cursor when no gesture handler is set', (tester) async {
      await tester.pumpApp(
        ds.Card(const Text('content')),
      );
      await tester.pumpAndSettle();

      expect(
        tester.resolvedCursorAt(find.byType(ds.Card)),
        SystemMouseCursors.basic,
      );
    });
  });
}
