import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/meta/device.manual.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('ManualConfiguration loose constraints resolutions', () {
    for (final entry in Resolutions.all.entries) {
      testWidgets('renders without overflow at ${entry.key}', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: entry.value,
          fit: FlexFit.loose,
          ManualConfiguration(connect: (_) {}),
        );

        await tester.pumpAndSettle();
        expect(tester.takeException(), isNull);
        final size = tester.getSize(find.byType(ManualConfiguration));
        expect(size.width, lessThanOrEqualTo(entry.value.width));
        expect(size.height, lessThanOrEqualTo(entry.value.height));
      });
    }
  });

  group('ManualConfiguration tight constraints resolutions', () {
    for (final entry in Resolutions.all.entries) {
      testWidgets('renders without overflow at ${entry.key}', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: entry.value,
          fit: FlexFit.tight,
          ManualConfiguration(connect: (_) {}),
        );

        await tester.pumpAndSettle();
        expect(tester.takeException(), isNull);
        final size = tester.getSize(find.byType(ManualConfiguration));
        expect(size, equals(entry.value));
      });
    }
  });
}
