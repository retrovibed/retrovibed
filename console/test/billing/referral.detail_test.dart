import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/billing/referral.detail.dart';
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
                attributionCount: Int64(5),
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
  group('billing.ReferralDetail', () {
    group('layout', () {
      testWidgets('renders referred users count', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(_wrap(const ReferralDetail()));
        await tester.pumpAndSettle();

        expect(find.text('5'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders monthly revenue', (WidgetTester tester) async {
        await tester.pumpApp(_wrap(const ReferralDetail()));
        await tester.pumpAndSettle();

        expect(find.text('\$0.50'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

    });

    group('resolutions', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          _wrap(const ReferralDetail()),
        );
        await tester.pumpAndSettle();

        expect(find.text('referred users'), findsOneWidget);
        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });
  });
}
