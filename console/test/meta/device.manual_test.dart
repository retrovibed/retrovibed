import 'dart:io';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/meta/api.dart' as api;
import 'package:retrovibed/meta/device.manual.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('ManualConfiguration loose constraints resolutions', () {
    for (final entry in Resolutions.all.entries) {
      testWidgets('renders without overflow at ${entry.key}', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: entry.value,
          fit: FlexFit.loose,
          ManualConfiguration((_) {}, apiconnect: (_) => Future.value(api.Daemon())),
        );

        await tester.pumpAndSettle();
        expect(tester.takeException(), isNull);
        final size = tester.getSize(find.byType(ManualConfiguration));
        expect(size.width, lessThanOrEqualTo(entry.value.width));
        expect(size.height, lessThanOrEqualTo(entry.value.height));
      });
    }
  });

  group('ManualConfiguration tight constraints resolutions', () {
    for (final entry in Resolutions.all.entries) {
      testWidgets('renders without overflow at ${entry.key}', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          physicalSize: entry.value,
          fit: FlexFit.tight,
          ManualConfiguration((_) {}, apiconnect: (_) => Future.value(api.Daemon())),
        );

        await tester.pumpAndSettle();
        expect(tester.takeException(), isNull);
        final size = tester.getSize(find.byType(ManualConfiguration));
        expect(size, equals(entry.value));
      });
    }
  });

  group('ManualConfiguration _connect error handling', () {
    testWidgets('offline error shows a create-anyway confirmation and creates on confirm', (
      WidgetTester tester,
    ) async {
      bool created = false;
      final daemon = api.Daemon(hostname: 'localhost:9998');

      await tester.pumpApp(
        ds.Node(
          ManualConfiguration(
            (_) {},
            apiconnect: (_) async {
              throw SocketException('', osError: OSError('', 111));
            },
            apicreate: (req) async {
              created = true;
              return api.DaemonCreateResponse(daemon: daemon);
            },
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('connect'));
      await tester.pumpN(3);

      expect(
        find.textContaining('unable to connect to daemon'),
        findsOneWidget,
      );

      await tester.tap(find.text('Create Anyway'));
      await tester.pumpAndSettle();

      expect(created, true);
    });

    testWidgets('offline error does not create when cancelled', (
      WidgetTester tester,
    ) async {
      bool created = false;

      await tester.pumpApp(
        ds.Node(
          ManualConfiguration(
            (_) {},
            apiconnect: (_) async {
              throw SocketException('', osError: OSError('', 111));
            },
            apicreate: (req) async {
              created = true;
              return api.DaemonCreateResponse(daemon: req.daemon);
            },
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('connect'));
      await tester.pumpN(3);

      await tester.tap(find.text('Cancel'));
      await tester.pumpAndSettle();

      expect(created, false);
    });

    testWidgets('dns resolution error shows a distinct confirmation and creates on confirm', (
      WidgetTester tester,
    ) async {
      bool created = false;
      final daemon = api.Daemon(hostname: 'localhost:9998');

      await tester.pumpApp(
        ds.Node(
          ManualConfiguration(
            (_) {},
            apiconnect: (_) async {
              throw SocketException('lookup failed', osError: OSError('', -2));
            },
            apicreate: (req) async {
              created = true;
              return api.DaemonCreateResponse(daemon: daemon);
            },
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('connect'));
      await tester.pumpN(3);

      expect(
        find.textContaining('unable to resolve hostname'),
        findsOneWidget,
      );
      expect(find.textContaining('unable to connect to daemon'), findsNothing);

      await tester.tap(find.text('Create Anyway'));
      await tester.pumpAndSettle();

      expect(created, true);
    });

    testWidgets('http 404 shows the httpauto not-found error', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ManualConfiguration(
          (_) {},
          apiconnect: (_) async {
            throw http.Response('', 404);
          },
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('connect'));
      await tester.pumpAndSettle();

      expect(find.text('not found'), findsOneWidget);
    });

    testWidgets('http 409 shows the httpauto conflict error', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ManualConfiguration(
          (_) {},
          apiconnect: (_) async {
            throw http.Response('', 409);
          },
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('connect'));
      await tester.pumpAndSettle();

      expect(
        find.text('a conflict occurred, the resource may already exist'),
        findsOneWidget,
      );
    });

    testWidgets('unexpected error shows the generic unknown error', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ManualConfiguration(
          (_) {},
          apiconnect: (_) async {
            throw Exception('boom');
          },
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('connect'));
      await tester.pumpAndSettle();

      expect(
        find.text('an unexpected problem has occurred'),
        findsOneWidget,
      );
    });

    testWidgets('successful connect creates the daemon and calls onConnected', (
      WidgetTester tester,
    ) async {
      bool createCalled = false;
      api.Daemon? connected;

      await tester.pumpApp(
        ManualConfiguration(
          (d) => connected = d,
          apiconnect: (d) async => d,
          apicreate: (req) async {
            createCalled = true;
            return api.DaemonCreateResponse(daemon: req.daemon);
          },
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('connect'));
      await tester.pumpAndSettle();

      expect(createCalled, true);
      expect(connected, equals(api.Daemon(hostname: httpx.localhost())));
    });
  });
}
