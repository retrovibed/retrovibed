import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/meta/api.dart' as api;
import 'package:retrovibed/meta/daemon.mdns.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

Widget _subject({
  required Future<String> Function() discover,
  Widget Function(void Function(api.Daemon) connect, {void Function()? retry})?
  preamble,
}) {
  return MDNSDiscovery(
    daemon: (_) {},
    discover: discover,
    preamble:
        preamble ??
        (connect, {retry}) => InitialSetup(connect: connect, retry: retry),
  );
}

void main() {
  group('MDNSDiscovery resolutions', () {
    for (final entry in Resolutions.all.entries) {
      testWidgets('NoLocalService renders without overflow at ${entry.key}', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: entry.value,
          _subject(
            discover: () async => throw Exception('no service'),
            preamble:
                (connect, {retry}) =>
                    NoLocalService(connect: connect, retry: retry),
          ),
        );

        await tester.pumpAndSettle();
        expect(tester.takeException(), isNull);
        final w = tester.getSize(find.byType(MDNSDiscovery));
        expect(w.width, lessThanOrEqualTo(entry.value.width));
        expect(w.height, lessThanOrEqualTo(entry.value.height));
      });

      testWidgets('InitialSetup renders without overflow at ${entry.key}', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: entry.value,
          _subject(discover: () async => throw Exception('no service')),
        );

        await tester.pumpAndSettle();
        expect(tester.takeException(), isNull);
        final initialSetupSize = tester.getSize(find.byType(MDNSDiscovery));
        expect(initialSetupSize.width, lessThanOrEqualTo(entry.value.width));
        expect(initialSetupSize.height, lessThanOrEqualTo(entry.value.height));
      });
    }
  });

  group('MDNSDiscovery discover results', () {
    testWidgets('shows preamble after discover fails', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        _subject(discover: () async => throw Exception('no service found')),
      );

      await tester.pumpAndSettle();

      expect(find.byType(InitialSetup), findsOneWidget);
    });

    testWidgets('shows preamble after discover succeeds', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(_subject(discover: () async => 'localhost:9998'));

      await tester.pumpAndSettle();

      expect(find.byType(InitialSetup), findsOneWidget);
    });

    testWidgets('shows preamble with NoLocalService after discover fails', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        _subject(
          discover: () async => throw Exception('no service found'),
          preamble:
              (connect, {retry}) =>
                  NoLocalService(connect: connect, retry: retry),
        ),
      );

      await tester.pumpAndSettle();

      expect(find.byType(NoLocalService), findsOneWidget);
    });
  });
}
