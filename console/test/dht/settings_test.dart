import 'package:material_ui/material_ui.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/dht/settings.dart';
import 'package:retrovibed/dht/api.dart' as api;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('Settings', () {
    testWidgets('renders without overflow', (WidgetTester tester) async {
      await tester.pumpApp(
        physicalSize: const Size(800, 800),
        Settings(api.DHTSettings(), onChange: (settings) async => settings),
      );
      await tester.pumpAndSettle();

      expect(find.text('passive'), findsOneWidget);
      expect(find.text('bootstrap'), findsOneWidget);
      expect(find.text('info'), findsOneWidget);
      expect(find.text('debug'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('toggling a checkbox emits the updated settings', (WidgetTester tester) async {
      api.DHTSettings? changed;
      await tester.pumpApp(
        physicalSize: const Size(800, 800),
        Settings(
          api.DHTSettings(),
          onChange: (settings) async => changed = settings,
        ),
      );
      await tester.pumpAndSettle();
      expect(changed, isNull);

      await tester.tap(find.text('passive'));
      await tester.pumpAndSettle();
      expect(changed?.passive, isTrue);
      expect(changed?.autoBootstrap, isFalse);

      await tester.tap(find.text('debug'));
      await tester.pumpAndSettle();
      expect(changed?.passive, isTrue);
      expect(changed?.debug, isTrue);
    });
  });
}
