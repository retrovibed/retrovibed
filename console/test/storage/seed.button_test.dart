import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/storage/seed.button.dart' as seed_button;
import 'package:retrovibed/uuidx.dart' as uuidx;

void main() {
  group('SeedButton', () {
    testWidgets('global current seed shows the globe', (WidgetTester tester) async {
      await tester.pumpApp(seed_button.SeedButton(uuidx.min(), onChange: (s) {}));
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.public), findsOneWidget);
      expect(find.byIcon(Icons.shield), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('global current seed has no current option, only the globe', (WidgetTester tester) async {
      await tester.pumpApp(seed_button.SeedButton(uuidx.min(), onChange: (s) {}));
      await tester.pumpAndSettle();

      await tester.tap(find.byType(PopupMenuButton<String>));
      await tester.pumpAndSettle();

      // the button and the global option.
      expect(find.byIcon(Icons.public), findsNWidgets(2));
      expect(find.byIcon(Icons.shield), findsOneWidget);
      expect(find.text('current'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('unique current seed has no current option', (WidgetTester tester) async {
      final current = uuidx.random();
      await tester.pumpApp(seed_button.SeedButton(current, onChange: (s) {}));
      await tester.pumpAndSettle();

      await tester.tap(find.byType(PopupMenuButton<String>));
      await tester.pumpAndSettle();

      // the button and the new private option.
      expect(find.byIcon(Icons.shield), findsNWidgets(2));
      expect(find.text('current'), findsNothing);
      expect(find.text(current), findsNothing);
      expect(tester.takeException(), isNull);
    });
  });
}
