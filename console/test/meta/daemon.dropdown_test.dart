import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/meta/api.dart' as api;
import 'package:retrovibed/meta/daemon.dropdown.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  final local = api.Daemon(id: 'local', hostname: 'localhost:9998', description: 'Local Library');
  final remote = api.Daemon(id: 'remote', hostname: 'example.com:9998', description: 'Remote Library');

  Future<api.DaemonSearchResponse> fakesearch(api.DaemonSearchRequest req) {
    return Future.value(api.DaemonSearchResponse(items: [local, remote]));
  }

  group('DaemonDropdown remoteonly', () {
    testWidgets('includes the local device by default', (WidgetTester tester) async {
      await tester.pumpApp(
        Scaffold(
          body: DaemonDropdown(
            library: ValueNotifier(api.Daemon(id: 'current')),
            search: fakesearch,
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
}
