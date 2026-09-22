import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

meta.DaemonSearchResponse _response({List<meta.Daemon> items = const []}) {
  return meta.DaemonSearchResponse(items: items, next: meta.DaemonSearchRequest());
}

meta.Daemon _daemon({String id = 'id-1', String description = 'Test Daemon'}) {
  return meta.Daemon(id: id, description: description, hostname: 'host-$id');
}

Future<meta.DaemonSearchResponse> _mockSearchWithItems(meta.DaemonSearchRequest req) {
  return Future.value(
    _response(items: [_daemon(id: 'id-1', description: 'Alice Server'), _daemon(id: 'id-2', description: 'Bob Server')]),
  );
}

Future<meta.DaemonSearchResponse> _mockSearchEmpty(meta.DaemonSearchRequest req) {
  return Future.value(_response());
}

Future<meta.DaemonSearchResponse> _mockSearchFailure(meta.DaemonSearchRequest req) {
  return Future.error(Exception('search failed'));
}

final _resolutions = Resolutions.variant();

void main() {
  group('meta.DaemonList', () {
    testWidgets('renders items once loaded', (WidgetTester tester) async {
      await tester.pumpApp(meta.DaemonList(search: _mockSearchWithItems));
      await tester.pumpAndSettle();

      expect(find.text('Alice Server'), findsOneWidget);
      expect(find.text('Bob Server'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders empty state when there are no devices', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(meta.DaemonList(search: _mockSearchEmpty));
      await tester.pumpAndSettle();

      expect(find.text('no other devices found'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('surfaces an error when the search fails', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(meta.DaemonList(search: _mockSearchFailure));
      await tester.pumpAndSettle();

      expect(find.text('an unexpected problem has occurred'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('hides add control in readonly mode', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        meta.DaemonList(search: _mockSearchWithItems, readonly: true),
      );
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.add), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        physicalSize: entry.value,
        meta.DaemonList(search: _mockSearchWithItems),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });
}
