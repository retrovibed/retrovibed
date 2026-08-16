import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta/api.dart' as api;
import 'package:retrovibed/meta/daemon.dropdown.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  final local = api.Daemon(id: 'local', hostname: 'localhost:9998', description: 'Local Library');
  final remote = api.Daemon(id: 'remote', hostname: 'example.com:9998', description: 'Remote Library');

  Future<api.DaemonSearchResponse> fakesearch(api.DaemonSearchRequest req) {
    return Future.value(api.DaemonSearchResponse(items: [local, remote]));
  }

  Future<Stream<api.Daemon>> noopDiscover({List<httpx.Option> options = const []}) async {
    return const Stream<api.Daemon>.empty();
  }

  group('DaemonDropdown remoteonly', () {
    testWidgets('includes the local device by default', (WidgetTester tester) async {
      await tester.pumpApp(
        Scaffold(
          body: DaemonDropdown(
            library: ValueNotifier(api.Daemon(id: 'current')),
            search: fakesearch,
            discover: noopDiscover,
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), 'lib');
      await tester.pumpAndSettle();

      expect(find.text('local library'), findsOneWidget);
      expect(find.text('Remote Library'), findsOneWidget);
    });

    testWidgets('excludes the local device when remoteonly is true', (WidgetTester tester) async {
      await tester.pumpApp(
        Scaffold(
          body: DaemonDropdown(
            library: ValueNotifier(api.Daemon(id: 'current')),
            search: fakesearch,
            discover: noopDiscover,
            remoteonly: true,
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextField), 'lib');
      await tester.pumpAndSettle();

      expect(find.text('local library'), findsNothing);
      expect(find.text('Remote Library'), findsOneWidget);
    });
  });

  group('DaemonDropdown discover', () {
    testWidgets('triggers a scan on mount', (WidgetTester tester) async {
      bool discoverCalled = false;
      Future<Stream<api.Daemon>> discover({List<httpx.Option> options = const []}) async {
        discoverCalled = true;
        return const Stream<api.Daemon>.empty();
      }

      await tester.pumpApp(
        Scaffold(
          body: DaemonDropdown(
            library: ValueNotifier(api.Daemon(id: 'current')),
            search: fakesearch,
            discover: discover,
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(discoverCalled, true);
    });

    testWidgets('refreshes an open dropdown as new peers stream in', (WidgetTester tester) async {
      final controller = StreamController<api.Daemon>();
      final discovered = api.Daemon(id: 'discovered', hostname: 'discovered.local:9998', description: 'Discovered Library');
      bool includeDiscovered = false;

      Future<api.DaemonSearchResponse> search(api.DaemonSearchRequest req) {
        final items = includeDiscovered ? [local, remote, discovered] : [local, remote];
        return Future.value(api.DaemonSearchResponse(items: items));
      }

      Future<Stream<api.Daemon>> discover({List<httpx.Option> options = const []}) async {
        return controller.stream;
      }

      // the discover stream is left open for most of this test (to simulate
      // an in-progress scan), which keeps the dropdown's indeterminate
      // "scanning" spinner animating — pumpAndSettle would hang waiting for
      // it, so a bounded number of pumps is used instead until the stream
      // is closed below.
      await tester.pumpApp(
        Scaffold(
          body: DaemonDropdown(
            library: ValueNotifier(api.Daemon(id: 'current')),
            search: search,
            discover: discover,
          ),
        ),
      );
      await tester.pump();
      await tester.pump();

      await tester.enterText(find.byType(TextField), 'lib');
      await tester.pump();
      await tester.pump();

      expect(find.text('Discovered Library'), findsNothing);

      includeDiscovered = true;
      controller.add(discovered);
      // closing ends the scan (stops the loading spinner's indeterminate
      // animation) so pumpAndSettle can actually settle.
      await controller.close();
      await tester.pumpAndSettle();

      expect(find.text('Discovered Library'), findsOneWidget);
    });
  });
}
