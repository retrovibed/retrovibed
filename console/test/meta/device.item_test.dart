import 'dart:io';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/meta/device.item.dart';
import 'package:retrovibed/meta/api.dart' as api;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  final daemon = api.Daemon(description: 'Test Library');
  Future<api.Daemon> onSelect(BuildContext context, api.Daemon daemon) async => daemon;

  Future<void> pumpFailing(WidgetTester tester, DaemonOnSelect failingSelect) {
    return tester.pumpApp(
      Scaffold(
        body: SizedBox(
          width: 200,
          height: 100,
          child: DaemonDropdownItem(library: daemon, onTap: (_) {}, onSelect: failingSelect),
        ),
      ),
    );
  }

  group('DaemonDropdownItem error handling', () {
    testWidgets('offline: connection refused shows the daemon-unreachable message', (
      WidgetTester tester,
    ) async {
      Future<api.Daemon> failingSelect(BuildContext context, api.Daemon daemon) {
        return Future.error(
          SocketException(
            'Connection refused',
            osError: OSError('Connection refused', 111),
            address: InternetAddress.loopbackIPv4,
            port: 9998,
          ),
        );
      }

      await pumpFailing(tester, failingSelect);
      await tester.pumpAndSettle();

      await tester.tap(find.text('Test Library'));
      await tester.pumpAndSettle();

      expect(
        find.text('unable to connect to daemon, is it running? check 127.0.0.1:9998.'),
        findsOneWidget,
      );
      expect(tester.takeException(), isNull);
    });

    testWidgets('offline: invalid argument (unreachable/malformed endpoint) shows the daemon-unreachable message', (
      WidgetTester tester,
    ) async {
      // reproduces: SocketException: Connection failed (OS Error: Invalid
      // argument, errno = 22), address = eg, port = 9998
      Future<api.Daemon> failingSelect(BuildContext context, api.Daemon daemon) {
        return Future.error(
          SocketException(
            'Connection failed',
            osError: OSError('Invalid argument', 22),
            address: InternetAddress.loopbackIPv4,
            port: 9998,
          ),
        );
      }

      await pumpFailing(tester, failingSelect);
      await tester.pumpAndSettle();

      await tester.tap(find.text('Test Library'));
      await tester.pumpAndSettle();

      expect(
        find.text('unable to connect to daemon, is it running? check 127.0.0.1:9998.'),
        findsOneWidget,
      );
      expect(tester.takeException(), isNull);
    });

    testWidgets('offline error can be dismissed by tapping it, restoring the row', (
      WidgetTester tester,
    ) async {
      Future<api.Daemon> failingSelect(BuildContext context, api.Daemon daemon) {
        return Future.error(
          SocketException(
            'Connection refused',
            osError: OSError('Connection refused', 111),
            address: InternetAddress.loopbackIPv4,
            port: 9998,
          ),
        );
      }

      await pumpFailing(tester, failingSelect);
      await tester.pumpAndSettle();

      await tester.tap(find.text('Test Library'));
      await tester.pumpAndSettle();

      final message = find.text('unable to connect to daemon, is it running? check 127.0.0.1:9998.');
      expect(message, findsOneWidget);

      await tester.tap(message);
      await tester.pumpAndSettle();

      expect(message, findsNothing);
      expect(find.text('Test Library'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('httpauto: 404 shows "not found"', (WidgetTester tester) async {
      Future<api.Daemon> failingSelect(BuildContext context, api.Daemon daemon) {
        return Future.error(http.Response('not found', 404));
      }

      await pumpFailing(tester, failingSelect);
      await tester.pumpAndSettle();

      await tester.tap(find.text('Test Library'));
      await tester.pumpAndSettle();

      expect(find.text('not found'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('httpauto: 401 shows the permissions message', (
      WidgetTester tester,
    ) async {
      Future<api.Daemon> failingSelect(BuildContext context, api.Daemon daemon) {
        return Future.error(http.Response('unauthorized', 401));
      }

      await pumpFailing(tester, failingSelect);
      await tester.pumpAndSettle();

      await tester.tap(find.text('Test Library'));
      await tester.pumpAndSettle();

      expect(find.text('you lack sufficient permissions'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('unknown: unrelated errors fall back to the generic message', (
      WidgetTester tester,
    ) async {
      Future<api.Daemon> failingSelect(BuildContext context, api.Daemon daemon) {
        return Future.error(Exception('something unrelated broke'));
      }

      await pumpFailing(tester, failingSelect);
      await tester.pumpAndSettle();

      await tester.tap(find.text('Test Library'));
      await tester.pumpAndSettle();

      expect(find.text('an unexpected problem has occurred'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('DaemonDropdownItem constrained parent', () {
    testWidgets('renders within fixed SizedBox constraints', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            width: 200,
            height: 150,
            child: DaemonDropdownItem(library: daemon, onTap: (_) {}, onSelect: onSelect),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Test Library'), findsOneWidget);
      expect(tester.takeException(), isNull);

      final size = tester.getSize(find.byType(DaemonDropdownItem));
      expect(size.width, equals(200));
      expect(size.height, equals(150));
    });

    testWidgets('renders in Column with fixed height', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: Column(
            children: [
              SizedBox(
                height: 100,
                child: DaemonDropdownItem(library: daemon, onTap: (_) {}, onSelect: onSelect),
              ),
              Expanded(child: Container()),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Test Library'), findsOneWidget);
      expect(tester.takeException(), isNull);

      final size = tester.getSize(find.byType(DaemonDropdownItem));
      expect(size.height, equals(100));
    });

    testWidgets('renders in Row with fixed width', (WidgetTester tester) async {
      await tester.pumpApp(
        Scaffold(
          body: Row(
            children: [
              SizedBox(
                width: 150,
                child: DaemonDropdownItem(library: daemon, onTap: (_) {}, onSelect: onSelect),
              ),
              Expanded(child: Container()),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Test Library'), findsOneWidget);
      expect(tester.takeException(), isNull);

      final size = tester.getSize(find.byType(DaemonDropdownItem));
      expect(size.width, equals(150));
    });

    testWidgets('renders with small dimensions', (WidgetTester tester) async {
      await tester.pumpApp(
        Scaffold(
          body: Center(
            child: SizedBox(
              width: 80,
              height: 64,
              child: DaemonDropdownItem(library: daemon, onTap: (_) {}, onSelect: onSelect),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);

      final size = tester.getSize(find.byType(DaemonDropdownItem));
      expect(size.width, equals(80));
      expect(size.height, equals(64));
    });

    testWidgets('renders with zero width constraint', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: Center(
            child: SizedBox(
              width: 0,
              height: 100,
              child: DaemonDropdownItem(library: daemon, onTap: (_) {}, onSelect: onSelect),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });

    testWidgets('renders with zero height constraint', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: Center(
            child: SizedBox(
              width: 100,
              height: 0,
              child: DaemonDropdownItem(library: daemon, onTap: (_) {}, onSelect: onSelect),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });
  });

  group('DaemonDropdownItem unconstrained parent', () {
    testWidgets('renders in ListView with fixed height child', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: ListView(
            children: [
              SizedBox(
                height: 200,
                child: DaemonDropdownItem(library: daemon, onTap: (_) {}, onSelect: onSelect),
              ),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Test Library'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in SingleChildScrollView with fixed height', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SingleChildScrollView(
            child: SizedBox(
              height: 300,
              child: DaemonDropdownItem(library: daemon, onTap: (_) {}, onSelect: onSelect),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Test Library'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in horizontal ListView with fixed width child', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            height: 100,
            child: ListView(
              scrollDirection: Axis.horizontal,
              children: [
                SizedBox(
                  width: 200,
                  child: DaemonDropdownItem(library: daemon, onTap: (_) {}, onSelect: onSelect),
                ),
              ],
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Test Library'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
