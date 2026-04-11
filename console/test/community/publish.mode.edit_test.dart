import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/billing/meta.billing.pb.dart';
import 'package:retrovibed/billing/plan.summary.dart';
import 'package:retrovibed/billing/registered.dart';
import 'package:retrovibed/community/community.pb.dart';
import 'package:retrovibed/community/publish.mode.edit.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

Widget _withPlan(String planId, Widget child) => Registered(
  child,
  lookup:
      ({options = const []}) => Future.value(
        BillingLookupResponse(
          billing: Billing(customerId: 'cust', planId: planId),
          plan: PlanSummary.plan(family()),
        ),
      ),
  create: ({options = const []}) => Future.value(BillingCreateResponse()),
);

void main() {
  group('PublishModeSelector', () {
    group('no plan (UNLISTED max)', () {
      for (final mode in PublishMode.values) {
        testWidgets('renders without overflow — $mode', (tester) async {
          final entry = _resolutions.currentValue!;

          await tester.pumpApp(
            physicalSize: entry.value,
            PublishModeEdit(publishMode: mode, onChanged: (_) {}),
          );
          await tester.pumpAndSettle();

          expect(tester.takeException(), isNull);
        }, variant: _resolutions);
      }
    });

    group('personal plan (LISTED max)', () {
      for (final mode in PublishMode.values) {
        testWidgets('renders without overflow — $mode', (tester) async {
          final entry = _resolutions.currentValue!;

          await tester.pumpApp(
            physicalSize: entry.value,
            _withPlan(
              personal().id,
              PublishModeEdit(publishMode: mode, onChanged: (_) {}),
            ),
          );
          await tester.pumpAndSettle();

          expect(tester.takeException(), isNull);
        }, variant: _resolutions);
      }
    });

    group('family plan (SYNDICATED max)', () {
      for (final mode in PublishMode.values) {
        testWidgets('renders without overflow — $mode', (tester) async {
          final entry = _resolutions.currentValue!;

          await tester.pumpApp(
            physicalSize: entry.value,
            _withPlan(
              family().id,
              PublishModeEdit(publishMode: mode, onChanged: (_) {}),
            ),
          );
          await tester.pumpAndSettle();

          expect(tester.takeException(), isNull);
        }, variant: _resolutions);
      }
    });
  });
}
