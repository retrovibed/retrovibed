import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/community/community.create.dart';
import 'package:retrovibed/community/community.pb.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

Future<CommunityCreateResponse> _noop(Community _) => Completer<CommunityCreateResponse>().future;

void main() {
  group('CommunityCreate', () {
    testWidgets('submits url and description to create', (tester) async {
      Community? received;

      await tester.pumpApp(
        physicalSize: const Size(800, 620),
        CommunityCreate(
          create: (c) {
            received = c;
            return Completer<CommunityCreateResponse>().future;
          },
          onCreate: (_) {},
          onCancel: () {},
        ),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextFormField).first, 'https://mydomain.community.retrovibe.space');
      await tester.enterText(find.byType(TextFormField).at(1), 'my description');
      await tester.tap(find.text('Create'));
      await tester.pump();

      expect(received, isNotNull);
      expect(received!.url, equals('https://mydomain.community.retrovibe.space'));
      expect(received!.description, equals('my description'));
      expect(tester.takeException(), isNull);
    });

    testWidgets('allows create with a blank url (server generates a default)', (tester) async {
      Community? received;

      await tester.pumpApp(
        physicalSize: const Size(800, 620),
        CommunityCreate(
          create: (c) {
            received = c;
            return Completer<CommunityCreateResponse>().future;
          },
          onCreate: (_) {},
          onCancel: () {},
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('Create'));
      await tester.pump();

      expect(received, isNotNull);
      expect(tester.takeException(), isNull);
    });

    group('renders without overflow', () {
      testWidgets('default', (tester) async {
        final entry = _resolutions.currentValue!;

        await tester.pumpApp(
          physicalSize: entry.value,
          SingleChildScrollView(
            child: CommunityCreate(
              create: _noop,
              onCreate: (_) {},
              onCancel: () {},
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });
  });
}
