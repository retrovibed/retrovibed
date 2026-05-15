import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/billing/cancellation.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/design.kit/modals.dart' as modals;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

Future<void> _noop({List<httpx.Option> options = const []}) => Future.value();
Future<void> _pending({List<httpx.Option> options = const []}) => Completer<void>().future;
Future<void> _failing({List<httpx.Option> options = const []}) => Future.error(Exception('api error'));

Widget _wrap(Widget child) => modals.Node(child);

// asyncfn keeps the LoadingIconButton in loading state (indefinite CircularProgressIndicator)
// until the modal is dismissed. Use pump() after tapping the delete icon so the modal
// renders, then pumpAndSettle() after Yes/No which resolves the future and stops the spinner.

void main() {
  group('billing.CancellationButton', () {
    group('layout', () {
      testWidgets('renders delete icon', (tester) async {
        await tester.pumpApp(_wrap(CancellationButton()));
        await tester.pumpAndSettle();

        expect(find.byIcon(Icons.delete), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in unconstrained environment', (tester) async {
        await tester.pumpApp(_wrap(CancellationButton()));
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in constrained environment', (tester) async {
        await tester.pumpApp(
          ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 400, maxHeight: 600),
            child: _wrap(CancellationButton()),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });
    });

    group('resolutions', () {
      testWidgets('renders without overflow', (tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          _wrap(CancellationButton()),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('without billing permissions', () {
      group('confirmation', () {
        testWidgets('tapping delete icon shows yes/no confirmation', (tester) async {
          await tester.pumpApp(_wrap(CancellationButton(apiidentitydelete: _pending)));
          await tester.pumpAndSettle();

          await tester.tap(find.byIcon(Icons.delete));
          await tester.pump();

          expect(find.text('Yes'), findsOneWidget);
          expect(find.text('No'), findsOneWidget);
          expect(tester.takeException(), isNull);
        });

        testWidgets('tapping No dismisses confirmation without calling API', (tester) async {
          bool billingCalled = false;
          bool identityCalled = false;
          await tester.pumpApp(
            _wrap(
              CancellationButton(
                apibillingdelete: ({options = const []}) {
                  billingCalled = true;
                  return Future.value();
                },
                apiidentitydelete: ({options = const []}) {
                  identityCalled = true;
                  return Future.value();
                },
              ),
            ),
          );
          await tester.pumpAndSettle();

          await tester.tap(find.byIcon(Icons.delete));
          await tester.pump();

          await tester.tap(find.text('No'));
          await tester.pumpAndSettle();

          expect(billingCalled, isFalse);
          expect(identityCalled, isFalse);
          expect(find.text('Yes'), findsNothing);
          expect(find.text('No'), findsNothing);
          expect(tester.takeException(), isNull);
        });

        testWidgets('tapping Yes calls apiidentitydelete', (tester) async {
          bool called = false;
          await tester.pumpApp(
            _wrap(
              CancellationButton(
                apiidentitydelete: ({options = const []}) {
                  called = true;
                  return Future.value();
                },
              ),
            ),
          );
          await tester.pumpAndSettle();

          await tester.tap(find.byIcon(Icons.delete));
          await tester.pump();

          await tester.tap(find.text('Yes'));
          await tester.pumpAndSettle();

          expect(called, isTrue);
          expect(tester.takeException(), isNull);
        });

        testWidgets('tapping Yes does not call billing API', (tester) async {
          bool called = false;
          await tester.pumpApp(
            _wrap(
              CancellationButton(
                apibillingdelete: ({options = const []}) {
                  called = true;
                  return Future.value();
                },
              ),
            ),
          );
          await tester.pumpAndSettle();

          await tester.tap(find.byIcon(Icons.delete));
          await tester.pump();

          await tester.tap(find.text('Yes'));
          await tester.pumpAndSettle();

          expect(called, isFalse);
          expect(tester.takeException(), isNull);
        });

        testWidgets('tapping Yes dismisses confirmation', (tester) async {
          await tester.pumpApp(_wrap(CancellationButton()));
          await tester.pumpAndSettle();

          await tester.tap(find.byIcon(Icons.delete));
          await tester.pump();

          await tester.tap(find.text('Yes'));
          await tester.pumpAndSettle();

          expect(find.text('Yes'), findsNothing);
          expect(find.text('No'), findsNothing);
          expect(tester.takeException(), isNull);
        });
      });
    });

    group('with billing permissions', () {
      Future<void> pumpWithBilling(WidgetTester tester, Widget child) {
        return tester.pumpApp(
          child,
          authzCurrent: authn.AuthzCache.fakeWith(meta.Token()..billingModify = true),
        );
      }

      group('confirmation', () {
        testWidgets('tapping delete icon shows yes/no confirmation', (tester) async {
          await pumpWithBilling(tester, _wrap(CancellationButton(apibillingdelete: _pending)));
          await tester.pumpAndSettle();

          await tester.tap(find.byIcon(Icons.delete));
          await tester.pump();

          expect(find.text('Yes'), findsOneWidget);
          expect(find.text('No'), findsOneWidget);
          expect(tester.takeException(), isNull);
        });

        testWidgets('tapping No dismisses confirmation without calling API', (tester) async {
          bool called = false;
          await pumpWithBilling(
            tester,
            _wrap(
              CancellationButton(
                apibillingdelete: ({options = const []}) {
                  called = true;
                  return Future.value();
                },
              ),
            ),
          );
          await tester.pumpAndSettle();

          await tester.tap(find.byIcon(Icons.delete));
          await tester.pump();

          await tester.tap(find.text('No'));
          await tester.pumpAndSettle();

          expect(called, isFalse);
          expect(find.text('Yes'), findsNothing);
          expect(find.text('No'), findsNothing);
          expect(tester.takeException(), isNull);
        });

        testWidgets('tapping Yes calls apibillingdelete', (tester) async {
          bool called = false;
          await pumpWithBilling(
            tester,
            _wrap(
              CancellationButton(
                apibillingdelete: ({options = const []}) {
                  called = true;
                  return Future.value();
                },
              ),
            ),
          );
          await tester.pumpAndSettle();

          await tester.tap(find.byIcon(Icons.delete));
          await tester.pump();

          await tester.tap(find.text('Yes'));
          await tester.pumpAndSettle();

          expect(called, isTrue);
          expect(tester.takeException(), isNull);
        });

        testWidgets('tapping Yes dismisses confirmation', (tester) async {
          await pumpWithBilling(tester, _wrap(CancellationButton(apibillingdelete: _noop)));
          await tester.pumpAndSettle();

          await tester.tap(find.byIcon(Icons.delete));
          await tester.pump();

          await tester.tap(find.text('Yes'));
          await tester.pumpAndSettle();

          expect(find.text('Yes'), findsNothing);
          expect(find.text('No'), findsNothing);
          expect(tester.takeException(), isNull);
        });
      });

      group('errors', () {
        testWidgets('API failure dismisses confirmation', (tester) async {
          await pumpWithBilling(tester, _wrap(CancellationButton(apibillingdelete: _failing)));
          await tester.pumpAndSettle();

          await tester.tap(find.byIcon(Icons.delete));
          await tester.pump();

          await tester.tap(find.text('Yes'));
          await tester.pumpAndSettle();

          expect(find.text('Yes'), findsNothing);
          expect(find.text('No'), findsNothing);
          expect(tester.takeException(), isNull);
        });

        testWidgets('API failure calls API exactly once', (tester) async {
          int callCount = 0;
          await pumpWithBilling(
            tester,
            _wrap(
              CancellationButton(
                apibillingdelete: ({options = const []}) {
                  callCount++;
                  return Future.error(Exception('fail'));
                },
              ),
            ),
          );
          await tester.pumpAndSettle();

          await tester.tap(find.byIcon(Icons.delete));
          await tester.pump();
          await tester.tap(find.text('Yes'));
          await tester.pumpAndSettle();

          expect(callCount, equals(1));
          expect(tester.takeException(), isNull);
        });
      });
    });
  });
}
