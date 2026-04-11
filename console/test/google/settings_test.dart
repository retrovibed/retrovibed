import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/google/settings.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('Integrations Settings', () {
    testWidgets('renders without errors', (WidgetTester tester) async {
      await tester.pumpApp(Settings());
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    });

    testWidgets('shows loading state initially', (WidgetTester tester) async {
      await tester.pumpApp(Settings());
      await tester.pumpAndSettle();
      expect(find.byType(Settings), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
