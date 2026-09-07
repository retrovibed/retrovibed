import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/billing/api.dart' as api;
import 'package:retrovibed/billing/registered.dart';
import 'package:retrovibed/billing/settings.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

api.BillingPlansResponse _plansResponse(List<api.Plan> plans) {
  return api.BillingPlansResponse(plans: plans);
}

api.Plan _plan(String id, {bool hidden = false}) {
  return api.Plan(id: id, token: 'tok.$id', stripeId: id, hidden: hidden);
}

// Settings gates the plan dropdown on authn.developer(context).subscription,
// which is read from the Login ancestor's cached flags. Without a real Login
// ancestor (as in plain pumpApp), those flags default to all-false.
Widget _withLogin(Widget child) {
  return authn.Login(
    child,
    publicKey: () => 'ssh-ed25519 AAAA...',
    seed: (_, __) => Future.value(),
  );
}

Widget _withRegistered(Widget child, {String planId = '', String customerId = ''}) {
  return Registered(
    child,
    lookup:
        ({options = const []}) => Future.value(
          api.BillingLookupResponse(billing: api.Billing(planId: planId, customerId: customerId)),
        ),
    create: ({options = const []}) => Future.error(Exception('not called')),
  );
}

void main() {
  group('billing.Settings', () {
    group('layout', () {
      testWidgets('renders plan label', (WidgetTester tester) async {
        await tester.pumpApp(
          _withLogin(
            Settings(
              apibillingplans: ({options = const []}) => Future.value(_plansResponse([])),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('plan'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in constrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 400, maxHeight: 600),
            child: _withLogin(
              Settings(
                apibillingplans: ({options = const []}) => Future.value(_plansResponse([])),
              ),
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
              Expanded(
                child: SingleChildScrollView(
                  child: _withLogin(
                    Settings(
                      apibillingplans: ({options = const []}) => Future.value(_plansResponse([])),
                    ),
                  ),
                ),
              ),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in Row', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Row(
            children: [
              Expanded(
                child: _withLogin(
                  Settings(
                    apibillingplans: ({options = const []}) => Future.value(_plansResponse([])),
                  ),
                ),
              ),
            ],
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in SizedBox', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 400,
            height: 600,
            child: SingleChildScrollView(
              child: _withLogin(
                Settings(
                  apibillingplans: ({options = const []}) => Future.value(_plansResponse([])),
                ),
              ),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders with custom margin and padding', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          _withLogin(
            Settings(
              margin: const EdgeInsets.all(8),
              padding: const EdgeInsets.all(8),
              apibillingplans: ({options = const []}) => Future.value(_plansResponse([])),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders with topLeft alignment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          _withLogin(
            Settings(
              alignment: Alignment.topLeft,
              apibillingplans: ({options = const []}) => Future.value(_plansResponse([])),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders with center alignment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          _withLogin(
            Settings(
              alignment: Alignment.center,
              apibillingplans: ({options = const []}) => Future.value(_plansResponse([])),
            ),
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
          _withLogin(
            Settings(
              apibillingplans: ({options = const []}) => Future.value(_plansResponse([])),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('plan'), findsOneWidget);
        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('plans', () {
      testWidgets('populates dropdown from API response', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          _withLogin(
            Settings(
              apibillingplans:
                  ({options = const []}) => Future.value(
                    _plansResponse([
                      _plan('plans.free'),
                      _plan('plans.founder'),
                    ]),
                  ),
            ),
          ),
        );
        await tester.pumpAndSettle();
        expect(find.text('free'), findsOneWidget);
        expect(find.text('founder'), findsNothing);

        await tester.tap(find.text('free'));
        await tester.pumpAndSettle();
        expect(find.text('free'), findsExactly(2));
        expect(find.text('founder'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('hides hidden plans from dropdown', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          _withLogin(
            Settings(
              apibillingplans:
                  ({options = const []}) => Future.value(
                    _plansResponse([
                      _plan('plans.free'),
                      _plan('plans.personal.2025', hidden: true),
                      _plan('plans.founder'),
                    ]),
                  ),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('free'), findsOneWidget);
        expect(find.text('founder'), findsNothing);

        await tester.tap(find.text('free'));
        await tester.pumpAndSettle();
        expect(find.text('free'), findsExactly(2));
        expect(find.text('founder'), findsOneWidget);
        expect(find.text('personal'), findsNothing);
        expect(tester.takeException(), isNull);
      });

      testWidgets('shows current hidden plan in dropdown', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          _withLogin(
            _withRegistered(
              Settings(
                apibillingplans:
                    ({options = const []}) => Future.value(
                      _plansResponse([
                        _plan('plans.free'),
                        _plan('plans.personal.2025', hidden: true),
                        _plan('plans.founder'),
                      ]),
                    ),
              ),
              planId: 'plans.personal.2025',
              customerId: 'derpy',
            ),
          ),
        );
        await tester.pumpAndSettle();
        await tester.tap(find.text('personal'));
        await tester.pumpAndSettle();
        expect(find.text('personal'), findsExactly(2));
        expect(tester.takeException(), isNull);
      });

      testWidgets('shows error widget on API failure', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          _withLogin(
            Settings(
              apibillingplans: ({options = const []}) => Future.error(Exception('plans unavailable')),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('an unexpected problem has occurred'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('clears error on retry tap', (WidgetTester tester) async {
        var callCount = 0;
        await tester.pumpApp(
          _withLogin(
            Settings(
              apibillingplans: ({options = const []}) {
                callCount++;
                if (callCount == 1) return Future.error(Exception('fail'));
                return Future.value(_plansResponse([_plan('plans.free')]));
              },
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('an unexpected problem has occurred'), findsOneWidget);

        await tester.tap(find.text('an unexpected problem has occurred'));
        await tester.pumpAndSettle();

        expect(find.text('an unexpected problem has occurred'), findsNothing);
        expect(tester.takeException(), isNull);
      });
    });
  });
}
