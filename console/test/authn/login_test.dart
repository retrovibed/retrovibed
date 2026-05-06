import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/authn/login.dart' as authn;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();
const _tall = Size(800, 900);

void main() {
  group('Login', () {
    group('layout', () {
      testWidgets('renders child when public key exists', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          authn.Login(
            const Text('authenticated content'),
            publicKey: () => 'ssh-ed25519 AAAA...',
            seed: (_) => '',
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('authenticated content'), findsOneWidget);
        expect(find.byType(TextFormField), findsNothing);
        expect(tester.takeException(), isNull);
      });

      testWidgets('shows error and form when authenticated fails', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: _tall,
          authn.Login(
            const Text('authenticated content'),
            publicKey: () => 'ssh-ed25519 AAAA...',
            seed: (_) => '',
            authenticated: () => Future.error(Exception('daemon failed')),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('authenticated content'), findsNothing);
        expect(find.text('an unexpected problem has occurred'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders form when public key is empty', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: _tall,
          authn.Login(
            const Text('authenticated content'),
            publicKey: () => '',
            seed: (_) => '',
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('authenticated content'), findsNothing);
        expect(find.text('email'), findsOneWidget);
        expect(find.text('password'), findsOneWidget);
        expect(find.text('Login'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in unconstrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: _tall,
          Center(
            child: authn.Login(
              const Text('child'),
              publicKey: () => '',
              seed: (_) => '',
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in constrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: _tall,
          ConstrainedBox(
            constraints: const BoxConstraints(
              maxWidth: 400,
              maxHeight: 900,
            ),
            child: authn.Login(
              const Text('child'),
              publicKey: () => '',
              seed: (_) => '',
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in Column', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: _tall,
          Column(
            children: [
              authn.Login(
                const Text('child'),
                publicKey: () => '',
                seed: (_) => '',
              ),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });
    });

    group('resolutions', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          SingleChildScrollView(
            child: authn.Login(
              const Text('child'),
              publicKey: () => '',
              seed: (_) => '',
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('password'), findsOneWidget);
        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('registration', () {
      testWidgets('auto-checks register when no account exists', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: _tall,
          authn.Login(
            const Text('child'),
            publicKey: () => '',
            seed: (_) => '',
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('register a new account'), findsOneWidget);
        expect(find.widgetWithText(TextFormField, 'confirm password'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('does not auto-check register when account exists', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: _tall,
          authn.Login(
            const Text('child'),
            publicKey: () => 'ssh-ed25519 AAAA...',
            seed: (_) => '',
            authenticated: () => Future.error(Exception('fail')),
          ),
        );
        await tester.pumpAndSettle();

        final checkbox = tester.widget<Checkbox>(
          find.byType(Checkbox).first,
        );
        expect(checkbox.value, isFalse);
        expect(tester.takeException(), isNull);
      });

      testWidgets('shows error when passwords do not match', (
        WidgetTester tester,
      ) async {
        var seedCalled = false;

        await tester.pumpApp(
          physicalSize: _tall,
          authn.Login(
            const Text('child'),
            publicKey: () => '',
            seed: (_) {
              seedCalled = true;
              return '';
            },
          ),
        );
        await tester.pumpAndSettle();

        await tester.enterText(find.widgetWithText(TextFormField, 'email'), 'user');
        await tester.enterText(find.widgetWithText(TextFormField, 'password'), 'pass1');
        await tester.enterText(find.widgetWithText(TextFormField, 'confirm password'), 'pass2');
        await tester.pumpAndSettle();
        await tester.tap(find.byType(Checkbox).at(1));
        await tester.pumpAndSettle();
        await tester.tap(find.text('Login'));
        await tester.pumpAndSettle();

        expect(find.text('passwords do not match'), findsOneWidget);
        expect(seedCalled, isFalse);
        expect(tester.takeException(), isNull);
      });

      testWidgets('calls seed when passwords match', (
        WidgetTester tester,
      ) async {
        String? captured;

        await tester.pumpApp(
          physicalSize: _tall,
          authn.Login(
            const Text('child'),
            publicKey: () => '',
            seed: (p) {
              captured = p;
              return 'error';
            },
          ),
        );
        await tester.pumpAndSettle();

        await tester.enterText(find.widgetWithText(TextFormField, 'email'), 'user');
        await tester.enterText(find.widgetWithText(TextFormField, 'password'), 'pass');
        await tester.enterText(find.widgetWithText(TextFormField, 'confirm password'), 'pass');
        await tester.pumpAndSettle();
        await tester.tap(find.byType(Checkbox).at(1));
        await tester.pumpAndSettle();
        await tester.tap(find.text('Login'));
        await tester.pumpAndSettle();

        expect(captured, 'user:pass');
        expect(tester.takeException(), isNull);
      });

      testWidgets('hides confirm field when register unchecked', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: _tall,
          authn.Login(
            const Text('child'),
            publicKey: () => '',
            seed: (_) => '',
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(TextFormField), findsNWidgets(3));

        await tester.tap(find.byType(Checkbox).first);
        await tester.pumpAndSettle();

        expect(find.byType(TextFormField), findsNWidgets(2));
        expect(tester.takeException(), isNull);
      });

      testWidgets('skips confirm validation when register unchecked', (
        WidgetTester tester,
      ) async {
        var seedCalled = false;

        await tester.pumpApp(
          physicalSize: _tall,
          authn.Login(
            const Text('child'),
            publicKey: () => '',
            seed: (_) {
              seedCalled = true;
              return 'error';
            },
          ),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.byType(Checkbox).first);
        await tester.pumpAndSettle();

        await tester.enterText(find.widgetWithText(TextFormField, 'email'), 'user');
        await tester.enterText(find.widgetWithText(TextFormField, 'password'), 'pass');
        await tester.pumpAndSettle();
        await tester.tap(find.byType(Checkbox).at(1));
        await tester.pumpAndSettle();
        await tester.tap(find.text('Login'));
        await tester.pumpAndSettle();

        expect(seedCalled, isTrue);
        expect(tester.takeException(), isNull);
      });
    });

    group('seed interaction', () {
      testWidgets('shows child after successful seed', (
        WidgetTester tester,
      ) async {
        var seeded = false;

        await tester.pumpApp(
          physicalSize: _tall,
          authn.Login(
            const Text('authenticated content'),
            publicKey: () => seeded ? 'ssh-ed25519 AAAA...' : '',
            seed: (p) {
              seeded = true;
              return '';
            },
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('authenticated content'), findsNothing);

        await tester.enterText(find.widgetWithText(TextFormField, 'email'), 'email');
        await tester.enterText(find.widgetWithText(TextFormField, 'password'), 'password');
        await tester.enterText(find.widgetWithText(TextFormField, 'confirm password'), 'password');
        await tester.pumpAndSettle();
        await tester.tap(find.byType(Checkbox).at(1));
        await tester.pumpAndSettle();
        await tester.tap(find.text('Login'));
        await tester.pumpAndSettle();

        expect(find.text('authenticated content'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('shows error on seed failure', (WidgetTester tester) async {
        await tester.pumpApp(
          physicalSize: _tall,
          authn.Login(
            const Text('authenticated content'),
            publicKey: () => '',
            seed: (_) => 'failed to generate key',
          ),
        );
        await tester.pumpAndSettle();

        await tester.enterText(find.widgetWithText(TextFormField, 'email'), 'email');
        await tester.enterText(find.widgetWithText(TextFormField, 'password'), 'password');
        await tester.enterText(find.widgetWithText(TextFormField, 'confirm password'), 'password');
        await tester.pumpAndSettle();
        await tester.tap(find.byType(Checkbox).at(1));
        await tester.pumpAndSettle();
        await tester.tap(find.text('Login'));
        await tester.pumpAndSettle();

        expect(find.text('login failed'), findsOneWidget);
        expect(find.text('authenticated content'), findsNothing);
        expect(tester.takeException(), isNull);
      });

      testWidgets('login button is disabled when password is empty', (
        WidgetTester tester,
      ) async {
        var seedCalled = false;

        await tester.pumpApp(
          physicalSize: _tall,
          authn.Login(
            const Text('child'),
            publicKey: () => '',
            seed: (_) {
              seedCalled = true;
              return '';
            },
          ),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.text('Login'));
        await tester.pumpAndSettle();

        expect(seedCalled, isFalse);
        expect(tester.takeException(), isNull);
      });

      testWidgets('login button enables after entering password', (
        WidgetTester tester,
      ) async {
        var seedCalled = false;

        await tester.pumpApp(
          physicalSize: _tall,
          authn.Login(
            const Text('child'),
            publicKey: () => '',
            seed: (_) {
              seedCalled = true;
              return 'error';
            },
          ),
        );
        await tester.pumpAndSettle();

        await tester.enterText(find.widgetWithText(TextFormField, 'email'), 'email');
        await tester.enterText(find.widgetWithText(TextFormField, 'password'), 'password');
        await tester.enterText(find.widgetWithText(TextFormField, 'confirm password'), 'password');
        await tester.pumpAndSettle();
        await tester.tap(find.byType(Checkbox).at(1));
        await tester.pumpAndSettle();
        await tester.tap(find.text('Login'));
        await tester.pumpAndSettle();

        expect(seedCalled, isTrue);
        expect(tester.takeException(), isNull);
      });

      testWidgets('passes password to seed function', (
        WidgetTester tester,
      ) async {
        String? capturedPassword;

        await tester.pumpApp(
          physicalSize: _tall,
          authn.Login(
            const Text('child'),
            publicKey: () => '',
            seed: (p) {
              capturedPassword = p;
              return 'error';
            },
          ),
        );
        await tester.pumpAndSettle();

        await tester.enterText(find.widgetWithText(TextFormField, 'email'), 'email');
        await tester.enterText(find.widgetWithText(TextFormField, 'password'), 'password');
        await tester.enterText(find.widgetWithText(TextFormField, 'confirm password'), 'password');
        await tester.pumpAndSettle();
        await tester.tap(find.byType(Checkbox).at(1));
        await tester.pumpAndSettle();
        await tester.tap(find.text('Login'));
        await tester.pumpAndSettle();

        expect(capturedPassword, 'email:password');
      });

      testWidgets('tapping seed error dismisses it', (
        WidgetTester tester,
      ) async {
        var clicked = false;
        await tester.pumpApp(
          physicalSize: _tall,
          authn.Login(
            const Text('authenticated content'),
            publicKey: () => '',
            seed: (_) {
              if (clicked) return '';
              clicked = true;
              return 'failed to generate key';
            },
          ),
        );
        await tester.pumpAndSettle();

        await tester.enterText(find.widgetWithText(TextFormField, 'email'), 'email');
        await tester.enterText(find.widgetWithText(TextFormField, 'password'), 'password');
        await tester.enterText(find.widgetWithText(TextFormField, 'confirm password'), 'password');
        await tester.pumpAndSettle();
        await tester.tap(find.byType(Checkbox).at(1));
        await tester.pumpAndSettle();
        await tester.tap(find.text('Login'));
        await tester.pumpAndSettle();

        expect(find.text('login failed'), findsOneWidget);
        await tester.tap(find.text('login failed'));
        await tester.pumpAndSettle();
        expect(find.text('login failed'), findsNothing);
        expect(tester.takeException(), isNull);
      });

      testWidgets('shows error when seed succeeds but authenticated fails', (
        WidgetTester tester,
      ) async {
        var seeded = false;

        await tester.pumpApp(
          physicalSize: _tall,
          authn.Login(
            const Text('authenticated content'),
            publicKey: () => seeded ? 'ssh-ed25519 AAAA...' : '',
            seed: (p) {
              seeded = true;
              return '';
            },
            authenticated: () => Future.error(Exception('daemon failed')),
          ),
        );
        await tester.pumpAndSettle();

        await tester.enterText(find.widgetWithText(TextFormField, 'email'), 'email');
        await tester.enterText(find.widgetWithText(TextFormField, 'password'), 'password');
        await tester.enterText(find.widgetWithText(TextFormField, 'confirm password'), 'password');
        await tester.pumpAndSettle();
        await tester.tap(find.byType(Checkbox).at(1));
        await tester.pumpAndSettle();
        await tester.tap(find.text('Login'));
        await tester.pumpAndSettle();

        expect(find.text('authenticated content'), findsNothing);
        expect(find.text('an unexpected problem has occurred'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('logout', () {
      testWidgets('returns to login screen after logout', (
        WidgetTester tester,
      ) async {
        var loggedOut = false;

        await tester.pumpApp(
          authn.Login(
            Builder(
              builder: (context) => TextButton(
                onPressed: () {
                  loggedOut = true;
                  authn.Login.logout(context);
                },
                child: const Text('do logout'),
              ),
            ),
            publicKey: () => loggedOut ? '' : 'ssh-ed25519 AAAA...',
            seed: (_) => '',
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('do logout'), findsOneWidget);
        expect(find.text('email'), findsNothing);

        await tester.tap(find.text('do logout'));
        await tester.pumpAndSettle();

        expect(find.text('email'), findsOneWidget);
        expect(find.text('password'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });
  });
}
