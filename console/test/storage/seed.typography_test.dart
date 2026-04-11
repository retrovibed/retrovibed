import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/storage/seed.typography.dart' as seed_typography;
import 'package:retrovibed/storage/seed.dart' as seed;
import 'package:retrovibed/uuidx.dart' as uuidx;

void main() {
  group('SeedTypography', () {
    // Create a test classifier using proper UUID patterns
    final testClassifier = seed.Classifier(
      community: uuidx.withSuffix(1),
      personal: uuidx.withSuffix(2),
    );

    testWidgets('renders correctly without onChange callback', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        seed_typography.SeedTypography(
          uuidx.withSuffix(1), // community seed
          classifier: testClassifier,
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('community'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders correctly with onChange callback', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        seed_typography.SeedTypography(
          uuidx.withSuffix(1), // community seed
          classifier: testClassifier,
          onChange: (seed) {
            // onChange callback test - just ensure it doesn't crash
          },
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(DropdownButton<String>), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    group('layout - finite constraints', () {
      testWidgets('renders in Container with fixed dimensions', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Container(
            width: 200,
            height: 100,
            child: seed_typography.SeedTypography(
              uuidx.withSuffix(1), // community seed
              classifier: testClassifier,
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('community'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders in SizedBox with tight constraints', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 100,
            height: 50,
            child: seed_typography.SeedTypography(
              uuidx.withSuffix(1), // community seed
              classifier: testClassifier,
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('community'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders with Expanded in Row', (WidgetTester tester) async {
        await tester.pumpApp(
          Row(
            children: [
              Expanded(
                child: seed_typography.SeedTypography(
                  uuidx.withSuffix(1), // community seed
                  classifier: testClassifier,
                ),
              ),
              const SizedBox(width: 50),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('community'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders with Expanded in Column', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Column(
            children: [
              Expanded(
                child: seed_typography.SeedTypography(
                  uuidx.withSuffix(1), // community seed
                  classifier: testClassifier,
                ),
              ),
              const SizedBox(height: 50),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('community'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('layout - flex containers', () {
      testWidgets('renders in Row without overflow', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Row(
            children: [
              seed_typography.SeedTypography(
                uuidx.withSuffix(1), // community seed
                classifier: testClassifier,
              ),
              const Text('other content'),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('community'), findsOneWidget);
        expect(find.text('other content'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders in Column without overflow', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Column(
            children: [
              seed_typography.SeedTypography(
                uuidx.withSuffix(1), // community seed
                classifier: testClassifier,
              ),
              const Text('other content'),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('community'), findsOneWidget);
        expect(find.text('other content'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders in ListView', (WidgetTester tester) async {
        await tester.pumpApp(
          ListView(
            children: [
              seed_typography.SeedTypography(
                uuidx.withSuffix(1), // community seed
                classifier: testClassifier,
              ),
              seed_typography.SeedTypography(
                uuidx.withSuffix(2), // personal seed
                classifier: testClassifier,
              ),
              seed_typography.SeedTypography(
                uuidx.min(), // global seed
                classifier: testClassifier,
              ),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('community'), findsOneWidget);
        expect(find.text('personal'), findsOneWidget);
        expect(find.text('global'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders in horizontal ScrollView', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            child: Row(
              children: [
                seed_typography.SeedTypography(
                  uuidx.withSuffix(1), // community seed
                  classifier: testClassifier,
                ),
                const Text('scrollable content'),
              ],
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('community'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('layout - intrinsic dimensions', () {
      testWidgets('provides valid intrinsic width', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          IntrinsicWidth(
            child: seed_typography.SeedTypography(
              uuidx.withSuffix(1), // community seed
              classifier: testClassifier,
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('community'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('provides valid intrinsic height', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          IntrinsicHeight(
            child: seed_typography.SeedTypography(
              uuidx.withSuffix(1), // community seed
              classifier: testClassifier,
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('community'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('dropdown interaction', () {
      testWidgets('selecting global seed calls onChange with global seed', (
        WidgetTester tester,
      ) async {
        seed.Seed? selected;

        await tester.pumpApp(
          seed_typography.SeedTypography(
            uuidx.min(), // start with global seed
            classifier: testClassifier,
            onChange: (s) => selected = s,
          ),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.byType(DropdownButton<String>));
        await tester.pumpAndSettle();
        await tester.tap(find.text('community').last);
        await tester.pumpAndSettle();

        expect(selected!.id, equals(uuidx.withSuffix(1)));

        await tester.tap(find.byType(DropdownButton<String>));
        await tester.pumpAndSettle();
        await tester.tap(find.text('global').last);
        await tester.pumpAndSettle();

        expect(selected!.id, equals(uuidx.min()));
      });
    });

    group('seed classification', () {
      testWidgets('renders global seed correctly', (WidgetTester tester) async {
        await tester.pumpApp(
          seed_typography.SeedTypography(
            uuidx.min(), // Empty string should be treated as global
            classifier: testClassifier,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('global'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders community seed correctly', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          seed_typography.SeedTypography(
            uuidx.withSuffix(1), // community seed
            classifier: testClassifier,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('community'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders personal seed correctly', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          seed_typography.SeedTypography(
            uuidx.withSuffix(2), // personal seed
            classifier: testClassifier,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('personal'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders unique seed correctly', (WidgetTester tester) async {
        await tester.pumpApp(
          seed_typography.SeedTypography(
            uuidx.random(), // random unique seed
            classifier: testClassifier,
          ),
        );
        await tester.pumpAndSettle();

        expect(
          find.text('private'),
          findsOneWidget,
        ); // Note: unique seed has label "private" in seed.dart
        expect(tester.takeException(), isNull);
      });
    });
  });
}
