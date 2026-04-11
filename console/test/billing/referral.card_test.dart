import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/billing/referral.card.dart';
import 'package:retrovibed/billing/registered.dart';
import 'package:retrovibed/billing/api.dart' as api;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

Widget _wrap(Widget child, {api.BillingLookupResponse? response}) {
  return Registered(
    child,
    lookup:
        ({options = const []}) => Future.value(
          response ??
              api.BillingLookupResponse(
                billing: api.Billing(customerId: 'cus_123'),
                attributionCount: Int64(7),
                attributionRate: 10,
              ),
        ),
    create:
        ({options = const []}) => Future.value(
          api.BillingCreateResponse(
            billing: api.Billing(customerId: 'cus_123'),
          ),
        ),
  );
}

void main() {
  group('billing.ReferralCard', () {
    group('layout', () {
      testWidgets('renders referrals title', (WidgetTester tester) async {
        await tester.pumpApp(
          _wrap(ReferralCard(onPressed: (_) {})),
        );
        await tester.pumpAndSettle();

        expect(find.text('Referrals'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders attribution count', (WidgetTester tester) async {
        await tester.pumpApp(
          _wrap(ReferralCard(onPressed: (_) {})),
        );
        await tester.pumpAndSettle();

        expect(find.text('7'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders monthly revenue', (WidgetTester tester) async {
        await tester.pumpApp(
          _wrap(ReferralCard(onPressed: (_) {})),
        );
        await tester.pumpAndSettle();

        expect(find.text('\$0.70/mo'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders zero when no attributions', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          _wrap(
            ReferralCard(onPressed: (_) {}),
            response: api.BillingLookupResponse(
              billing: api.Billing(customerId: 'cus_123'),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('0'), findsOneWidget);
        expect(find.text('\$0.00/mo'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('invokes onPressed with ReferralDetail', (
        WidgetTester tester,
      ) async {
        Widget? received;
        await tester.pumpApp(
          _wrap(ReferralCard(onPressed: (w) => received = w)),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.text('Referrals').first);
        await tester.pumpAndSettle();

        expect(received, isNotNull);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in constrained environment', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          _wrap(
            ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 400, maxHeight: 600),
              child: ReferralCard(onPressed: (_) {}),
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('renders without overflow in SizedBox', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          _wrap(
            SizedBox(
              width: 400,
              height: 400,
              child: ReferralCard(onPressed: (_) {}),
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
          _wrap(ReferralCard(onPressed: (_) {})),
        );
        await tester.pumpAndSettle();

        expect(find.text('Referrals'), findsOneWidget);
        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });
  });
}
