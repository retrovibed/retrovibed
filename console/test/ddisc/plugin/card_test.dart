import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/ddisc/plugin/card.dart' as ddiscplugin;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('ddiscplugin.Card', () {
    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        ddiscplugin.Card(onPressed: (_) {}),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });
}
