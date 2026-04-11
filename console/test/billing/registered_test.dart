import 'dart:async';
import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/billing/registered.dart';
import 'package:retrovibed/billing/api.dart' as api;
import 'package:retrovibed/billing/plan.summary.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('billing.Registered', () {
    group('initialization', () {
      testWidgets('shows loading spinner before lookup completes', (
        WidgetTester tester,
      ) async {
        final completer = Completer<api.BillingLookupResponse>();

        await tester.pumpApp(
          Registered(
            const Text('child content'),
            lookup: ({options = const []}) => completer.future,
            create: ({options = const []}) => Future.error(Exception('not called')),
          ),
        );
        await tester.pump();

        expect(find.byType(CircularProgressIndicator), findsOneWidget);
        expect(find.text('child content'), findsNothing);

        completer.complete(
          api.BillingLookupResponse(
            billing: api.Billing(customerId: 'cus_123'),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(CircularProgressIndicator), findsNothing);
        expect(find.text('child content'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders child after successful lookup', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Registered(
            const Text('child content'),
            lookup:
                ({options = const []}) => Future.value(
                  api.BillingLookupResponse(
                    billing: api.Billing(
                      customerId: 'cus_123',
                      planId: 'price_abc',
                    ),
                  ),
                ),
            create: ({options = const []}) => Future.error(Exception('not called')),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('child content'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('sets current billing from lookup response', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Registered(
            const Text('child'),
            lookup:
                ({options = const []}) => Future.value(
                  api.BillingLookupResponse(
                    billing: api.Billing(
                      customerId: 'cus_123',
                      planId: 'price_abc',
                    ),
                  ),
                ),
            create: ({options = const []}) => Future.error(Exception('not called')),
          ),
        );
        await tester.pumpAndSettle();

        final state = tester.state<RegisteredState>(find.byType(Registered));
        expect(state.current.customerId, equals('cus_123'));
        expect(state.current.planId, equals('price_abc'));
        expect(tester.takeException(), isNull);
      });

      testWidgets('populates attribution fields from lookup', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Registered(
            const Text('child'),
            lookup:
                ({options = const []}) => Future.value(
                  api.BillingLookupResponse(
                    billing: api.Billing(customerId: 'cus_123'),
                    attributionCount: Int64(5),
                    attributionRate: 10,
                  ),
                ),
            create: ({options = const []}) => Future.error(Exception('not called')),
          ),
        );
        await tester.pumpAndSettle();

        final state = tester.state<RegisteredState>(find.byType(Registered));
        expect(state.attributionCount, equals(5));
        expect(state.attributionRate, equals(10));
        expect(tester.takeException(), isNull);
      });

      testWidgets('calls create when lookup returns empty customerId', (
        WidgetTester tester,
      ) async {
        var createCalled = false;

        await tester.pumpApp(
          Registered(
            const Text('child content'),
            lookup:
                ({options = const []}) => Future.value(
                  api.BillingLookupResponse(
                    billing: api.Billing(planId: 'price_abc'),
                  ),
                ),
            create: ({options = const []}) {
              createCalled = true;
              return Future.value(
                api.BillingCreateResponse(
                  billing: api.Billing(
                    customerId: 'cus_new',
                    planId: 'price_abc',
                  ),
                ),
              );
            },
          ),
        );
        await tester.pumpAndSettle();

        expect(createCalled, isTrue);
        expect(find.text('child content'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('sets billing from create response', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Registered(
            const Text('child'),
            lookup:
                ({options = const []}) => Future.value(
                  api.BillingLookupResponse(billing: api.Billing()),
                ),
            create:
                ({options = const []}) => Future.value(
                  api.BillingCreateResponse(
                    billing: api.Billing(
                      customerId: 'cus_new',
                      planId: 'price_new',
                    ),
                  ),
                ),
          ),
        );
        await tester.pumpAndSettle();

        final state = tester.state<RegisteredState>(find.byType(Registered));
        expect(state.current.customerId, equals('cus_new'));
        expect(tester.takeException(), isNull);
      });

      testWidgets('defaults to free plan on 404', (WidgetTester tester) async {
        await tester.pumpApp(
          Registered(
            const Text('child content'),
            lookup: ({options = const []}) => Future.error(http.Response('', 404)),
            create: ({options = const []}) => Future.error(Exception('not called')),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('child content'), findsOneWidget);
        final state = tester.state<RegisteredState>(find.byType(Registered));
        expect(state.current.planId, equals(free().id));
        expect(tester.takeException(), isNull);
      });

      testWidgets('shows error on unexpected lookup failure', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Registered(
            const Text('child content'),
            lookup: ({options = const []}) => Future.error(Exception('network error')),
            create: ({options = const []}) => Future.error(Exception('not called')),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('an unexpected problem has occurred'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('replace', () {
      testWidgets('updates current', (WidgetTester tester) async {
        await tester.pumpApp(
          Registered(
            const Text('child'),
            lookup:
                ({options = const []}) => Future.value(
                  api.BillingLookupResponse(
                    billing: api.Billing(
                      customerId: 'cus_123',
                      planId: 'price_old',
                    ),
                  ),
                ),
            create: ({options = const []}) => Future.error(Exception('not called')),
          ),
        );
        await tester.pumpAndSettle();

        final state = tester.state<RegisteredState>(find.byType(Registered));
        state.replace(api.BillingLookupResponse(billing: api.Billing(customerId: 'cus_123', planId: 'price_new')));
        await tester.pumpAndSettle();

        expect(state.current.planId, equals('price_new'));
        expect(tester.takeException(), isNull);
      });

      testWidgets('updates refresh notifier', (WidgetTester tester) async {
        await tester.pumpApp(
          Registered(
            const Text('child'),
            lookup:
                ({options = const []}) => Future.value(
                  api.BillingLookupResponse(
                    billing: api.Billing(customerId: 'cus_123'),
                  ),
                ),
            create: ({options = const []}) => Future.error(Exception('not called')),
          ),
        );
        await tester.pumpAndSettle();

        final state = tester.state<RegisteredState>(find.byType(Registered));
        api.Billing? notified;
        state.refresh.addListener(() => notified = state.refresh.value);

        state.replace(api.BillingLookupResponse(billing: api.Billing(customerId: 'cus_123', planId: 'price_new')));
        await tester.pumpAndSettle();

        expect(notified?.planId, equals('price_new'));
        expect(tester.takeException(), isNull);
      });
    });

    group('of', () {
      testWidgets('returns RegisteredState from descendant context', (
        WidgetTester tester,
      ) async {
        RegisteredState? found;

        await tester.pumpApp(
          Registered(
            Builder(
              builder: (context) {
                found = Registered.of(context);
                return const Text('child');
              },
            ),
            lookup:
                ({options = const []}) => Future.value(
                  api.BillingLookupResponse(
                    billing: api.Billing(customerId: 'cus_123'),
                  ),
                ),
            create: ({options = const []}) => Future.error(Exception('not called')),
          ),
        );
        await tester.pumpAndSettle();

        expect(found, isNotNull);
        expect(tester.takeException(), isNull);
      });
    });
  });
}
