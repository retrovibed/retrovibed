import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/usermanagement/display.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('usermanagement.Display', () {
    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        const Display(),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });
}
