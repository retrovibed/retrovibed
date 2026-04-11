// import 'dart:async';

// import 'package:flutter/material.dart';
// import 'package:flutter_test/flutter_test.dart';
// import 'package:retrovibed/billing/purchase.dart';
// import 'package:retrovibed/billing/api.dart' as api;
// import 'package:retrovibed/billing/plan.summary.dart';
// import 'package:retrovibed/testing/widget_tester_extensions.dart';

// Future<void> _noop(Future<api.Billing> _) => Future.value();
// Future<void> _noopSession(String plan, {List<dynamic> options = const []}) => Future.value();

void main() {
  // Widget _purchase({
  //   PlanSummary? current,
  //   PlanSummary? desired,
  //   Future<void> Function(Future<api.Billing>) onChange = _noop,
  //   Future<void> Function(String plan, {List<dynamic> options}) session = _noopSession,
  //   Future<api.BillingLookupResponse> Function({List<dynamic> options})? lookup,
  //   Duration interval = Duration.zero,
  // }) {
  //   final c = current ?? free();
  //   final d = desired ?? founder();
  //   return Purchase(
  //     current: c,
  //     desired: PlanSummary.plan(d),
  //     onChange: onChange,
  //     session: session,
  //     lookup:
  //         lookup ??
  //         ({options = const []}) => Future.value(
  //           api.BillingLookupResponse(billing: api.Billing(planId: d.id)),
  //         ),
  //     interval: interval,
  //   );
  // }

  // group('billing.Purchase', () {
  // group('button state', () {
  //   testWidgets('button is disabled when current equals desired', (
  //     WidgetTester tester,
  //   ) async {
  //     var changed = false;
  //     await tester.pumpApp(
  //       _purchase(
  //         current: free(),
  //         desired: free(),
  //         onChange: (_) async => changed = true,
  //       ),
  //     );
  //     await tester.pumpAndSettle();

  //     await tester.tap(find.text('upgrade'));
  //     await tester.pumpAndSettle();

  //     expect(changed, isFalse);
  //     expect(tester.takeException(), isNull);
  //   });

  //   testWidgets('button is enabled when plans differ', (
  //     WidgetTester tester,
  //   ) async {
  //     var changed = false;
  //     await tester.pumpApp(
  //       _purchase(
  //         current: free(),
  //         desired: founder(),
  //         onChange: (_) async => changed = true,
  //       ),
  //     );
  //     await tester.pumpAndSettle();

  //     await tester.tap(find.text('upgrade'));
  //     await tester.pumpAndSettle();

  //     expect(changed, isTrue);
  //     expect(tester.takeException(), isNull);
  //   });
  // });

  //   group('poll', () {
  //     testWidgets(
  //       'resolves immediately when lookup matches desired on first call',
  //       (WidgetTester tester) async {
  //         api.Billing? result;

  //         await tester.pumpApp(
  //           _purchase(
  //             current: free(),
  //             desired: founder(),
  //             onChange: (future) async => result = await future,
  //             lookup:
  //                 ({options = const []}) => Future.value(
  //                   api.BillingLookupResponse(
  //                     billing: api.Billing(planId: founder().id),
  //                   ),
  //                 ),
  //           ),
  //         );
  //         await tester.pumpAndSettle();

  //         await tester.tap(find.text('upgrade'));
  //         await tester.pumpAndSettle();

  //         expect(result?.planId, equals(founder().id));
  //         expect(tester.takeException(), isNull);
  //       },
  //     );

  //     testWidgets('retries until lookup returns desired plan', (
  //       WidgetTester tester,
  //     ) async {
  //       var callCount = 0;
  //       api.Billing? result;

  //       await tester.pumpApp(
  //         _purchase(
  //           current: free(),
  //           desired: founder(),
  //           onChange: (future) async => result = await future,
  //           lookup: ({options = const []}) {
  //             callCount++;
  //             final planId = callCount >= 3 ? founder().id : free().id;
  //             return Future.value(
  //               api.BillingLookupResponse(billing: api.Billing(planId: planId)),
  //             );
  //           },
  //         ),
  //       );
  //       await tester.pumpAndSettle();

  //       await tester.tap(find.text('upgrade'));
  //       await tester.pumpAndSettle();

  //       expect(callCount, equals(3));
  //       expect(result?.planId, equals(founder().id));
  //       expect(tester.takeException(), isNull);
  //     });

  //     testWidgets('calls onChange with the resolved billing', (
  //       WidgetTester tester,
  //     ) async {
  //       final billing = api.Billing(
  //         planId: founder().id,
  //         customerId: 'cus_abc',
  //       );
  //       api.Billing? received;

  //       await tester.pumpApp(
  //         _purchase(
  //           current: free(),
  //           desired: founder(),
  //           onChange: (future) async => received = await future,
  //           lookup: ({options = const []}) => Future.value(api.BillingLookupResponse(billing: billing)),
  //         ),
  //       );
  //       await tester.pumpAndSettle();

  //       await tester.tap(find.text('upgrade'));
  //       await tester.pumpAndSettle();

  //       expect(received?.customerId, equals('cus_abc'));
  //       expect(tester.takeException(), isNull);
  //     });
  //   });

  //   group('dismount', () {
  //     testWidgets('poll stops when widget is removed during delay', (
  //       WidgetTester tester,
  //     ) async {
  //       var lookupCallCount = 0;
  //       final firstLookup = Completer<api.BillingLookupResponse>();

  //       // visible flag controls whether Purchase is in the tree
  //       bool visible = true;
  //       late StateSetter setVisible;

  //       await tester.pumpWidget(
  //         MaterialApp(
  //           home: Material(
  //             child: StatefulBuilder(
  //               builder: (context, setState) {
  //                 setVisible = setState;
  //                 return visible
  //                     ? _purchase(
  //                       current: free(),
  //                       desired: founder(),
  //                       onChange: (_) async {},
  //                       lookup: ({options = const []}) {
  //                         lookupCallCount++;
  //                         return firstLookup.future;
  //                       },
  //                       interval: Duration.zero,
  //                     )
  //                     : const SizedBox();
  //               },
  //             ),
  //           ),
  //         ),
  //       );

  //       await tester.tap(find.text('upgrade'));
  //       await tester.pump();

  //       // dismount Purchase while lookup is in flight
  //       setVisible(() => visible = false);
  //       await tester.pump();

  //       // complete with non-matching plan — would trigger next poll iteration
  //       firstLookup.complete(
  //         api.BillingLookupResponse(billing: api.Billing(planId: free().id)),
  //       );
  //       await tester.pumpAndSettle();

  //       expect(lookupCallCount, equals(1));
  //       expect(tester.takeException(), isNull);
  //     });
  //   });

  //   group('error handling', () {
  //     testWidgets('session error propagates to onChange', (
  //       WidgetTester tester,
  //     ) async {
  //       Object? caught;

  //       await tester.pumpApp(
  //         _purchase(
  //           current: free(),
  //           desired: founder(),
  //           session: (plan, {options = const []}) => Future.error(Exception('session failed')),
  //           onChange: (future) async {
  //             try {
  //               await future;
  //             } catch (e) {
  //               caught = e;
  //             }
  //           },
  //         ),
  //       );
  //       await tester.pumpAndSettle();

  //       await tester.tap(find.text('upgrade'));
  //       await tester.pumpAndSettle();

  //       expect(caught, isA<Exception>());
  //       expect(caught.toString(), contains('session failed'));
  //       expect(tester.takeException(), isNull);
  //     });

  //     testWidgets('lookup error propagates to onChange', (
  //       WidgetTester tester,
  //     ) async {
  //       Object? caught;

  //       await tester.pumpApp(
  //         _purchase(
  //           current: free(),
  //           desired: founder(),
  //           onChange: (future) async {
  //             try {
  //               await future;
  //             } catch (e) {
  //               caught = e;
  //             }
  //           },
  //           lookup: ({options = const []}) => Future.error(Exception('lookup failed')),
  //         ),
  //       );
  //       await tester.pumpAndSettle();

  //       await tester.tap(find.text('upgrade'));
  //       await tester.pumpAndSettle();

  //       expect(caught, isA<Exception>());
  //       expect(caught.toString(), contains('lookup failed'));
  //       expect(tester.takeException(), isNull);
  //     });
  //   });
  // });
}
