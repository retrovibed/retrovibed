import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/authn/login.dart' as authn;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

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
            seed: (_, __) => Future.value(),
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
          authn.Login(
            const Text('authenticated content'),
            publicKey: () => 'ssh-ed25519 AAAA...',
            seed: (_, __) => Future.value(),
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
          authn.Login(
            const Text('authenticated content'),
            publicKey: () => '',
            seed: (_, __) => Future.value(),
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
          Center(
            child: authn.Login(
              const Text('child'),
              publicKey: () => '',
              seed: (_, __) => Future.value(),
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
          ConstrainedBox(
            constraints: const BoxConstraints(
              maxWidth: 400,
              maxHeight: 600,
            ),
            child: authn.Login(
              const Text('child'),
              publicKey: () => '',
              seed: (_, __) => Future.value(),
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
          Column(
            children: [
              authn.Login(
                const Text('child'),
                publicKey: () => '',
                seed: (_, __) => Future.value(),
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
              seed: (_, __) => Future.value(),
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
          authn.Login(
            const Text('child'),
            publicKey: () => '',
            seed: (_, __) => Future.value(),
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
          authn.Login(
            const Text('child'),
            publicKey: () => 'ssh-ed25519 AAAA...',
            seed: (_, __) => Future.value(),
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

      testWidgets('login disabled when passwords do not match', (
        WidgetTester tester,
      ) async {
        var seedCalled = false;

        await tester.pumpApp(
          authn.Login(
            const Text('child'),
            publicKey: () => '',
            seed: (_, __) {
              seedCalled = true;
              return Future.value();
            },
          ),
        );
        await tester.pumpAndSettle();

        // register defaults to unchecked; turn it on to require a matching confirm field
        await tester.tap(find.byType(Checkbox).first);
        await tester.pumpAndSettle();

        await tester.enterText(find.widgetWithText(TextFormField, 'email'), 'user');
        await tester.enterText(find.widgetWithText(TextFormField, 'password'), 'pass1');
        await tester.enterText(find.widgetWithText(TextFormField, 'confirm password'), 'pass2');
        await tester.pumpAndSettle();
        await tester.tap(find.byType(Checkbox).at(1));
        await tester.pumpAndSettle();
        await tester.tap(find.text('Login'));
        await tester.pumpAndSettle();

        expect(seedCalled, isFalse);
        expect(tester.takeException(), isNull);
      });

      testWidgets('calls seed when passwords match', (
        WidgetTester tester,
      ) async {
        String? captured;

        await tester.pumpApp(
          authn.Login(
            const Text('child'),
            publicKey: () => '',
            seed: (u, p) {
              captured = '$u:$p';
              return Future.error('error');
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
          authn.Login(
            const Text('child'),
            publicKey: () => '',
            seed: (_, __) => Future.value(),
          ),
        );
        await tester.pumpAndSettle();

        // register defaults to unchecked: the confirm field stays mounted
        // (Visibility.maintainState is true) but isn't visible.
        expect(find.byType(TextFormField), findsNWidgets(3));
        expect(
          tester
              .widget<Visibility>(
                find.ancestor(
                  of: find.widgetWithText(TextFormField, 'confirm password'),
                  matching: find.byType(Visibility),
                ),
              )
              .visible,
          isFalse,
        );

        // checking register reveals the confirm field
        await tester.tap(find.byType(Checkbox).first);
        await tester.pumpAndSettle();
        expect(find.byType(TextFormField), findsNWidgets(3));
        expect(
          tester
              .widget<Visibility>(
                find.ancestor(
                  of: find.widgetWithText(TextFormField, 'confirm password'),
                  matching: find.byType(Visibility),
                ),
              )
              .visible,
          isTrue,
        );

        // unchecking it again hides the confirm field once more
        await tester.tap(find.byType(Checkbox).first);
        await tester.pumpAndSettle();
        expect(find.byType(TextFormField), findsNWidgets(3));
        expect(
          tester
              .widget<Visibility>(
                find.ancestor(
                  of: find.widgetWithText(TextFormField, 'confirm password'),
                  matching: find.byType(Visibility),
                ),
              )
              .visible,
          isFalse,
        );
        expect(tester.takeException(), isNull);
      });

      testWidgets('skips confirm validation when register unchecked', (
        WidgetTester tester,
      ) async {
        var seedCalled = false;

        await tester.pumpApp(
          authn.Login(
            const Text('child'),
            publicKey: () => '',
            seed: (_, __) {
              seedCalled = true;
              return Future.error('error');
            },
          ),
        );
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

      testWidgets(
        'confirm field keeps its value (and stays validated against it) across register toggling',
        (
          WidgetTester tester,
        ) async {
          String? captured;

          await tester.pumpApp(
            authn.Login(
              const Text('child'),
              publicKey: () => '',
              seed: (u, p) {
                captured = '$u:$p';
                return Future.value();
              },
            ),
          );
          await tester.pumpAndSettle();

          // register defaults to unchecked; turn it on to reveal the confirm field
          await tester.tap(find.byType(Checkbox).first);
          await tester.pumpAndSettle();

          await tester.enterText(find.widgetWithText(TextFormField, 'email'), 'user');
          await tester.enterText(find.widgetWithText(TextFormField, 'password'), 'pass1');
          await tester.enterText(find.widgetWithText(TextFormField, 'confirm password'), 'pass1');
          await tester.pumpAndSettle();

          // uncheck register: the confirm field is hidden but stays mounted
          // (Visibility.maintainState is true), so its value isn't lost.
          await tester.tap(find.byType(Checkbox).first);
          await tester.pumpAndSettle();
          expect(find.byType(TextFormField), findsNWidgets(3));

          // recheck register: the confirm field reappears still showing what
          // was typed, not blank, so what's validated matches what's displayed.
          await tester.tap(find.byType(Checkbox).first);
          await tester.pumpAndSettle();
          expect(find.text('pass1'), findsNWidgets(2));

          await tester.tap(find.byType(Checkbox).at(1));
          await tester.pumpAndSettle();

          await tester.tap(find.text('Login'));
          await tester.pumpAndSettle();

          expect(find.text('passwords do not match'), findsNothing);
          expect(captured, 'user:pass1');
          expect(tester.takeException(), isNull);
        },
      );
    });

    group('seed interaction', () {
      testWidgets('shows child after successful seed', (
        WidgetTester tester,
      ) async {
        var seeded = false;

        await tester.pumpApp(
          authn.Login(
            const Text('authenticated content'),
            publicKey: () => seeded ? 'ssh-ed25519 AAAA...' : '',
            seed: (_, __) {
              seeded = true;
              return Future.value();
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
          authn.Login(
            const Text('authenticated content'),
            publicKey: () => '',
            seed: (_, __) => Future.error('failed to generate key'),
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
          authn.Login(
            const Text('child'),
            publicKey: () => '',
            seed: (_, __) {
              seedCalled = true;
              return Future.value();
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
          authn.Login(
            const Text('child'),
            publicKey: () => '',
            seed: (_, __) {
              seedCalled = true;
              return Future.error('error');
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
          authn.Login(
            const Text('child'),
            publicKey: () => '',
            seed: (u, p) {
              capturedPassword = '$u:$p';
              return Future.error('error');
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

      testWidgets('displays entered username after login failure is dismissed', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          authn.Login(
            const Text('authenticated content'),
            publicKey: () => '',
            seed: (_, __) => Future.error('authentication failed'),
          ),
        );
        await tester.pumpAndSettle();

        await tester.enterText(find.widgetWithText(TextFormField, 'email'), 'user@example.com');
        await tester.enterText(find.widgetWithText(TextFormField, 'password'), 'secret');
        await tester.pumpAndSettle();
        await tester.tap(find.byType(Checkbox).at(1));
        await tester.pumpAndSettle();
        await tester.tap(find.text('Login'));
        await tester.pumpAndSettle();

        expect(find.text('login failed'), findsOneWidget);

        await tester.tap(find.text('login failed'));
        await tester.pumpAndSettle();

        expect(find.text('user@example.com'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('tapping seed error dismisses it', (
        WidgetTester tester,
      ) async {
        var clicked = false;
        await tester.pumpApp(
          authn.Login(
            const Text('authenticated content'),
            publicKey: () => '',
            seed: (_, __) {
              if (clicked) return Future.value();
              clicked = true;
              return Future.error('failed to generate key');
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
          authn.Login(
            const Text('authenticated content'),
            publicKey: () => seeded ? 'ssh-ed25519 AAAA...' : '',
            seed: (_, __) {
              seeded = true;
              return Future.value();
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

      testWidgets('ignores resubmission while the previous attempt is still awaiting authentication', (
        WidgetTester tester,
      ) async {
        final captured = <String>[];
        var seeded = false;

        await tester.pumpApp(
          authn.Login(
            const Text('authenticated content'),
            // publicKey only becomes non-empty after seed succeeds, so
            // _checkKey() actually reaches the pending authenticated()
            // call below instead of short-circuiting on an empty key.
            publicKey: () => seeded ? 'ssh-ed25519 AAAA...' : '',
            seed: (u, p) {
              captured.add('$u:$p');
              seeded = true;
              return Future.value();
            },
            authenticated: () => Completer<void>().future,
          ),
        );
        await tester.pumpAndSettle();

        // register defaults to unchecked, so the confirm-password mismatch
        // check can't interfere with isolating the in-flight-submission guard.
        await tester.enterText(find.widgetWithText(TextFormField, 'email'), 'user1');
        await tester.enterText(find.widgetWithText(TextFormField, 'password'), 'pass1');
        await tester.pumpAndSettle();
        await tester.tap(find.byType(Checkbox).at(1));
        await tester.pumpAndSettle();

        await tester.tap(find.text('Login'));
        // authenticated() never resolves, so the LoadingButton's spinner
        // animates forever; pumpAndSettle would time out here.
        await tester.pump();

        expect(find.text('authenticated content'), findsNothing);

        await tester.enterText(find.widgetWithText(TextFormField, 'email'), 'user2');
        await tester.enterText(find.widgetWithText(TextFormField, 'password'), 'pass2');
        await tester.pump();

        // the button is disabled/obscured while loading, so this tap is
        // expected not to land on it.
        await tester.tap(find.text('Login'), warnIfMissed: false);
        await tester.pump();

        expect(captured, equals(['user1:pass1']));
        expect(tester.takeException(), isNull);
      });
    });

    group('guest interaction', () {
      testWidgets('shows child after successful guest login', (WidgetTester tester) async {
        var guestCalled = false;

        await tester.pumpApp(
          authn.Login(
            const Text('authenticated content'),
            publicKey: () => guestCalled ? 'ssh-ed25519 AAAA...' : '',
            guest: () {
              guestCalled = true;
              return Future.value();
            },
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('authenticated content'), findsNothing);

        await tester.tap(find.byTooltip('continue as guest'));
        await tester.pumpAndSettle();

        expect(guestCalled, isTrue);
        expect(find.text('authenticated content'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('shows error on guest login failure', (WidgetTester tester) async {
        await tester.pumpApp(
          authn.Login(
            const Text('authenticated content'),
            publicKey: () => '',
            guest: () => Future.error('guest login failed'),
          ),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.byTooltip('continue as guest'));
        await tester.pumpAndSettle();

        expect(find.text('guest login failed'), findsOneWidget);
        expect(find.text('authenticated content'), findsNothing);
        expect(tester.takeException(), isNull);
      });

      testWidgets('does not require username, password, or terms acceptance', (WidgetTester tester) async {
        var guestCalled = false;

        await tester.pumpApp(
          authn.Login(
            const Text('authenticated content'),
            publicKey: () => guestCalled ? 'ssh-ed25519 AAAA...' : '',
            guest: () {
              guestCalled = true;
              return Future.value();
            },
          ),
        );
        await tester.pumpAndSettle();

        // no text entered, no checkboxes ticked
        await tester.tap(find.byTooltip('continue as guest'));
        await tester.pumpAndSettle();

        expect(guestCalled, isTrue);
        expect(find.text('authenticated content'), findsOneWidget);
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
            seed: (_, __) => Future.value(),
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

      testWidgets('does not display previous session credentials after logout', (
        WidgetTester tester,
      ) async {
        var loggedOut = false;
        var seeded = false;

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
            publicKey: () => (seeded && !loggedOut) ? 'ssh-ed25519 AAAA...' : '',
            seed: (_, __) {
              seeded = true;
              return Future.value();
            },
          ),
        );
        await tester.pumpAndSettle();

        // register defaults to unchecked, so email/password is all that's needed
        await tester.enterText(find.widgetWithText(TextFormField, 'email'), 'user@example.com');
        await tester.enterText(find.widgetWithText(TextFormField, 'password'), 'secret');
        await tester.pumpAndSettle();
        await tester.tap(find.byType(Checkbox).at(1));
        await tester.pumpAndSettle();
        await tester.tap(find.text('Login'));
        await tester.pumpAndSettle();

        expect(find.text('do logout'), findsOneWidget);

        await tester.tap(find.text('do logout'));
        await tester.pumpAndSettle();

        expect(find.text('email'), findsOneWidget);
        expect(find.text('user@example.com'), findsNothing);
        expect(find.text('secret'), findsNothing);
        expect(tester.takeException(), isNull);
      });
    });
  });
}
