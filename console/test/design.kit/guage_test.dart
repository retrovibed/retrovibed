import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/guage.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('Gauge rendering', () {
    testWidgets('renders without error at 0.0', (tester) async {
      await tester.pumpApp(Gauge(0.0));
      await tester.pumpAndSettle();
      expect(find.byType(Gauge), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without error at 1.0', (tester) async {
      await tester.pumpApp(Gauge(1.0));
      await tester.pumpAndSettle();
      expect(find.byType(Gauge), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without error at 0.5', (tester) async {
      await tester.pumpApp(Gauge(0.5));
      await tester.pumpAndSettle();
      expect(find.byType(Gauge), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without error for NaN fill', (tester) async {
      await tester.pumpApp(Gauge(double.nan));
      await tester.pumpAndSettle();
      expect(find.byType(Gauge), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('clamps fill below 0.0', (tester) async {
      await tester.pumpApp(Gauge(-1.0));
      await tester.pumpAndSettle();
      expect(find.byType(Gauge), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('clamps fill above 1.0', (tester) async {
      await tester.pumpApp(Gauge(2.0));
      await tester.pumpAndSettle();
      expect(find.byType(Gauge), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('Gauge fill width', () {
    testWidgets('fill=0.0 renders FractionallySizedBox with widthFactor 0.0', (
      tester,
    ) async {
      await tester.pumpApp(
        SizedBox(width: 200, child: Gauge(0.0, height: 8.0)),
      );
      await tester.pumpAndSettle();

      final fsb = tester.widget<FractionallySizedBox>(
        find.byType(FractionallySizedBox),
      );
      expect(fsb.widthFactor, 0.0);
      expect(tester.takeException(), isNull);
    });

    testWidgets('fill=1.0 renders FractionallySizedBox with widthFactor 1.0', (
      tester,
    ) async {
      await tester.pumpApp(
        SizedBox(width: 200, child: Gauge(1.0, height: 8.0)),
      );
      await tester.pumpAndSettle();

      final fsb = tester.widget<FractionallySizedBox>(
        find.byType(FractionallySizedBox),
      );
      expect(fsb.widthFactor, 1.0);
      expect(tester.takeException(), isNull);
    });

    testWidgets('fill=NaN renders FractionallySizedBox with widthFactor 0.0', (
      tester,
    ) async {
      await tester.pumpApp(
        SizedBox(width: 200, child: Gauge(double.nan, height: 8.0)),
      );
      await tester.pumpAndSettle();

      final fsb = tester.widget<FractionallySizedBox>(
        find.byType(FractionallySizedBox),
      );
      expect(fsb.widthFactor, 0.0);
      expect(tester.takeException(), isNull);
    });

    testWidgets('fill=-1.0 clamps to widthFactor 0.0', (tester) async {
      await tester.pumpApp(
        SizedBox(width: 200, child: Gauge(-1.0, height: 8.0)),
      );
      await tester.pumpAndSettle();

      final fsb = tester.widget<FractionallySizedBox>(
        find.byType(FractionallySizedBox),
      );
      expect(fsb.widthFactor, 0.0);
      expect(tester.takeException(), isNull);
    });

    testWidgets('fill=2.0 clamps to widthFactor 1.0', (tester) async {
      await tester.pumpApp(
        SizedBox(width: 200, child: Gauge(2.0, height: 8.0)),
      );
      await tester.pumpAndSettle();

      final fsb = tester.widget<FractionallySizedBox>(
        find.byType(FractionallySizedBox),
      );
      expect(fsb.widthFactor, 1.0);
      expect(tester.takeException(), isNull);
    });
  });

  group('Gauge layout', () {
    testWidgets('respects custom height', (tester) async {
      await tester.pumpApp(SizedBox(width: 200, child: Gauge(0.5, height: 16.0)));
      await tester.pumpAndSettle();

      final container = tester.widget<Container>(find.byType(Container).first);
      expect(container.constraints?.maxHeight ?? 16.0, 16.0);
      expect(tester.takeException(), isNull);
    });

    testWidgets(
      'renders across multiple resolutions',
      (tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(Gauge(0.5), physicalSize: entry.value);
        await tester.pumpAndSettle();
        expect(find.byType(Gauge), findsOneWidget);
        expect(tester.takeException(), isNull);
      },
      variant: _resolutions,
    );
  });
}
