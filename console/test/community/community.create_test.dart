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
    testWidgets('submits domain and description to create', (tester) async {
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

      await tester.enterText(find.byType(TextFormField).first, 'mydomain');
      await tester.enterText(find.byType(TextFormField).at(1), 'my description');
      await tester.tap(find.text('Create'));
      await tester.pump();

      expect(received, isNotNull);
      expect(received!.domain, equals('mydomain'));
      expect(received!.description, equals('my description'));
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
