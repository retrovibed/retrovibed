import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/profiles.dart' as profiles;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('profiles.Current', () {
    testWidgets('renders id and name fields once the session loads', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(profiles.Current());
      await tester.pumpAndSettle();

      expect(find.text('id'), findsOneWidget);
      expect(find.text('name'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        physicalSize: entry.value,
        profiles.Current(),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('surfaces an error when session fetch fails', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        authn.Authenticated(
          profiles.Current(),
          apissh: () => Future.error(Exception('ssh failed')),
          apisignup: () async => authn.Session(),
          apicurrent: (token) async => authn.Session(),
        ),
      );
      await tester.pumpAndSettle();

      // Both the Authenticated ancestor and Current's own error state
      // surface the same generic message.
      expect(find.text('an unexpected problem has occurred'), findsWidgets);
      expect(tester.takeException(), isNull);
    });
  });
}
