import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/authn/login.dart' as authn;
import 'package:retrovibed/debug/developer.settings.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

Widget _authenticated(Widget child) {
  return authn.Login(
    child,
    publicKey: () => 'ssh-ed25519 AAAA...',
    seed: (_) => '',
  );
}

void main() {
  group('DeveloperSettings', () {
    group('layout', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          _authenticated(const DeveloperSettings()),
          physicalSize: entry.value,
        );
        await tester.pumpAndSettle();
        expect(tester.takeException(), isNull);
      }, variant: _resolutions);

      testWidgets('displays all developer flag checkboxes', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(_authenticated(const DeveloperSettings()));
        await tester.pumpAndSettle();

        expect(find.text('Networking'), findsOneWidget);
        expect(find.text('Subscription'), findsOneWidget);
        expect(find.text('Recommendations'), findsOneWidget);
        expect(find.text('Releases'), findsOneWidget);
        expect(find.text('Debug'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('flags', () {
      testWidgets('toggling Networking updates the cached flag', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          _authenticated(
            Builder(
              builder: (context) {
                final flags = authn.Login.cached(context).flags;
                return Column(
                  children: [
                    const DeveloperSettings(),
                    Text('networking:${flags.networking}'),
                  ],
                );
              },
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('networking:false'), findsOneWidget);

        await tester.tap(find.byType(Checkbox).first);
        await tester.pumpAndSettle();

        expect(find.text('networking:true'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('toggling Subscription updates the cached flag', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          _authenticated(
            Builder(
              builder: (context) {
                final flags = authn.Login.cached(context).flags;
                return Column(
                  children: [
                    const DeveloperSettings(),
                    Text('subscription:${flags.subscription}'),
                  ],
                );
              },
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('subscription:false'), findsOneWidget);

        await tester.tap(find.byType(Checkbox).at(1));
        await tester.pumpAndSettle();

        expect(find.text('subscription:true'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('toggling one flag does not affect the other', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          _authenticated(
            Builder(
              builder: (context) {
                final flags = authn.Login.cached(context).flags;
                return Column(
                  children: [
                    const DeveloperSettings(),
                    Text('networking:${flags.networking}'),
                    Text('subscription:${flags.subscription}'),
                  ],
                );
              },
            ),
          ),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.byType(Checkbox).first);
        await tester.pumpAndSettle();

        expect(find.text('networking:true'), findsOneWidget);
        expect(find.text('subscription:false'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('toggling Recommendations updates the cached flag', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          _authenticated(
            Builder(
              builder: (context) {
                final flags = authn.Login.cached(context).flags;
                return Column(
                  children: [
                    const DeveloperSettings(),
                    Text('recommendations:${flags.recommendations}'),
                  ],
                );
              },
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('recommendations:true'), findsOneWidget);

        await tester.tap(find.byType(Checkbox).at(2));
        await tester.pumpAndSettle();

        expect(find.text('recommendations:false'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('toggling Releases updates the cached flag', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          _authenticated(
            Builder(
              builder: (context) {
                final flags = authn.Login.cached(context).flags;
                return Column(
                  children: [
                    const DeveloperSettings(),
                    Text('releases:${flags.releases}'),
                  ],
                );
              },
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('releases:true'), findsOneWidget);

        await tester.tap(find.byType(Checkbox).at(3));
        await tester.pumpAndSettle();

        expect(find.text('releases:false'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('toggling Debug updates the cached flag', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          _authenticated(
            Builder(
              builder: (context) {
                final flags = authn.Login.cached(context).flags;
                return Column(
                  children: [
                    const DeveloperSettings(),
                    Text('debug:${flags.debug}'),
                  ],
                );
              },
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('debug:true'), findsOneWidget);

        await tester.tap(find.byType(Checkbox).at(4));
        await tester.pumpAndSettle();

        expect(find.text('debug:false'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });
  });
}
