import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/main.dart';

void main() {
  testWidgets('smoke test', (WidgetTester tester) async {
    // Build our app and trigger a frame.
    await tester.pumpWidget(Retrovibed());
    expect(tester.takeException(), isNull);
  });
}
