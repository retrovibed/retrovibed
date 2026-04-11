import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/torrents/settings.dart';
import 'package:retrovibed/torrents/api.dart' as api;
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:nock/nock.dart';

void main() {
  group('Settings', () {
    nock('https://api.retrovibe.space')
        .get(RegExp(r'/s/torrents/'))
        .reply(
          200,
          api.TorrentSettings.create().toProto3Json(),
          headers: {'Content-Type': 'application/json'},
        );

    testWidgets('renders without overflow', (WidgetTester tester) async {
      final initial = api.TorrentSettings();

      await tester.pumpApp(
        physicalSize: const Size(800, 800),
        Settings(
          initial,
          onChange: (settings) async {
            return settings;
          },
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('download rate'), findsOneWidget);
      expect(find.text('upload rate'), findsOneWidget);
      expect(find.text('peers'), findsOneWidget);
      expect(find.text('log'), findsOneWidget);
      expect(find.text('debug'), findsOneWidget);
      expect(find.text('firewall'), findsOneWidget);

      expect(tester.takeException(), isNull);
    });

    testWidgets('renders correctly with initial data', (
      WidgetTester tester,
    ) async {
      final initialSettings = api.TorrentSettings();

      await tester.pumpApp(
        physicalSize: const Size(800, 800),
        Settings(
          initialSettings,
          onChange: (settings) async {
            return settings;
          },
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('download rate'), findsOneWidget);
      expect(find.text('upload rate'), findsOneWidget);
      expect(find.text('peers'), findsOneWidget);
      expect(find.text('log'), findsOneWidget);
      expect(find.text('debug'), findsOneWidget);
      expect(find.text('firewall'), findsOneWidget);
    });

    testWidgets('handles settings changes correctly', (
      WidgetTester tester,
    ) async {
      final initialSettings = api.TorrentSettings();

      bool onChangeCalled = false;

      await tester.pumpApp(
        physicalSize: const Size(800, 800),
        Settings(
          initialSettings,
          onChange: (settings) async {
            onChangeCalled = true;
            return settings;
          },
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('download rate'), findsOneWidget);
      expect(find.text('upload rate'), findsOneWidget);
      expect(find.text('peers'), findsOneWidget);
      expect(find.text('log'), findsOneWidget);
      expect(find.text('debug'), findsOneWidget);
      expect(find.text('firewall'), findsOneWidget);

      expect(onChangeCalled, isFalse);
    });
  });
}
