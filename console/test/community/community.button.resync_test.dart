import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/community/community.button.resync.dart';
import 'package:retrovibed/community/community.pb.dart';
import 'package:retrovibed/community/community.publish.pb.dart';

final _community = Community(id: 'c1', url: 'https://example-community.community.retrovibe.space');

void main() {
  group('ResyncButton', () {
    testWidgets('renders refresh icon', (tester) async {
      await tester.pumpApp(
        ResyncButton(community: _community, apiresync: (id, {options = const [], req}) async => throw UnimplementedError()),
      );
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.refresh), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('tapping calls apiresync with the community id', (tester) async {
      String? calledWith;

      await tester.pumpApp(
        ResyncButton(
          community: _community,
          apiresync: (id, {options = const [], req}) async {
            calledWith = id;
            return PublishedContentSearchResponse(community: _community);
          },
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.refresh));
      await tester.pumpAndSettle();

      expect(calledWith, equals('c1'));
      expect(tester.takeException(), isNull);
    });

    testWidgets('shows loading spinner while the request is pending', (tester) async {
      final completer = Completer<PublishedContentSearchResponse>();

      await tester.pumpApp(
        ResyncButton(community: _community, apiresync: (id, {options = const [], req}) => completer.future),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.refresh));
      await tester.pump();

      expect(find.byType(CircularProgressIndicator), findsOneWidget);
      expect(tester.takeException(), isNull);

      completer.complete(PublishedContentSearchResponse(community: _community));
      await tester.pumpAndSettle();

      expect(find.byType(CircularProgressIndicator), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('invokes onResynced with the refreshed community', (tester) async {
      Community? resynced;
      final refreshed = Community(id: 'c1', url: 'https://refreshed-domain.community.retrovibe.space');

      await tester.pumpApp(
        ResyncButton(
          community: _community,
          onResynced: (c) => resynced = c,
          apiresync: (id, {options = const [], req}) async => PublishedContentSearchResponse(community: refreshed),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.refresh));
      await tester.pumpAndSettle();

      expect(resynced?.url, equals('https://refreshed-domain.community.retrovibe.space'));
      expect(tester.takeException(), isNull);
    });

    testWidgets('failed request clears the loading spinner without calling onResynced', (tester) async {
      bool onResyncedCalled = false;

      await tester.pumpApp(
        Scaffold(
          body: ResyncButton(
            community: _community,
            onResynced: (_) => onResyncedCalled = true,
            apiresync: (id, {options = const [], req}) => Future<PublishedContentSearchResponse>.error(Exception('boom')),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.refresh));
      await tester.pumpAndSettle();

      expect(onResyncedCalled, isFalse);
      expect(find.byType(CircularProgressIndicator), findsNothing);
    });
  });
}
