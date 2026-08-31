import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/community/list.display.dart';
import 'package:retrovibed/community/community.pb.dart';
import 'package:retrovibed/community/api.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

Future<CommunitySearchResponse> _empty(
  CommunitySearchRequest req, {
  List<httpx.Option> options = const [],
}) async {
  return CommunitySearchResponse(items: []);
}

Future<CommunitySearchResponse> _noop(
  CommunitySearchRequest req, {
  List<httpx.Option> options = const [],
}) => Completer<CommunitySearchResponse>().future;

Future<CommunitySearchResponse> _withCommunities(
  CommunitySearchRequest req, {
  List<httpx.Option> options = const [],
}) async {
  return CommunitySearchResponse(
    items: [
      Community(id: '1', description: 'Alpha community', url: 'https://alpha.community.retrovibe.space'),
      Community(id: '2', description: 'Beta community', url: 'https://beta.community.retrovibe.space'),
    ],
  );
}

Future<CommunitySearchResponse> _withSubscription(
  CommunitySearchRequest req, {
  List<httpx.Option> options = const [],
}) async {
  return CommunitySearchResponse(
    items: [
      Community(id: '1', url: 'https://alpha.community.retrovibe.space', accountId: 'other', subscribedAt: '2026-03-20T00:00:00Z'),
      Community(id: '2', url: 'https://beta.community.retrovibe.space', accountId: 'other'),
    ],
  );
}

void main() {
  group('ListDisplay', () {
    group('resolutions - loading state', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          ListDisplay(search: _noop),
        );
        await tester.pump();

        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('resolutions - empty state', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          ListDisplay(search: _empty),
        );
        await tester.pumpAndSettle();

        expect(find.text('No communities found'), findsOneWidget);
        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('resolutions - with communities', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          ListDisplay(search: _withCommunities),
        );
        await tester.pumpAndSettle();

        expect(
          find.text('https://alpha.community.retrovibe.space'),
          findsOneWidget,
        );
        expect(find.text('https://beta.community.retrovibe.space'), findsOneWidget);
        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('subscription status', () {
      testWidgets('subscribed community shows check_circle', (tester) async {
        await tester.pumpApp(
          physicalSize: Size(1280, 720),
          ListDisplay(search: _withSubscription),
        );
        await tester.pumpAndSettle();

        expect(find.byIcon(Icons.check_circle), findsOneWidget);
        expect(find.byIcon(Icons.add_circle_outline), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('tapping subscribe icon updates to subscribed', (tester) async {
        int calls = 0;
        Future<CommunitySearchResponse> search(
          CommunitySearchRequest req, {
          List<httpx.Option> options = const [],
        }) async {
          calls++;
          return CommunitySearchResponse(
            items: [
              Community(
                id: '1',
                url: 'https://alpha.community.retrovibe.space',
                accountId: 'other',
                subscribedAt: calls > 1 ? '2026-03-20T00:00:00Z' : '',
              ),
            ],
          );
        }

        Future<CommunitySubscribeResponse> subscribe(String id, {List<httpx.Option> options = const []}) async => CommunitySubscribeResponse();

        await tester.pumpApp(
          physicalSize: Size(1280, 720),
          ListDisplay(search: search, subscribe: subscribe),
        );
        await tester.pumpAndSettle();

        expect(find.byIcon(Icons.add_circle_outline), findsOneWidget);
        expect(find.byIcon(Icons.check_circle), findsNothing);

        await tester.tap(find.byIcon(Icons.add_circle_outline));
        await tester.pumpAndSettle();

        expect(find.byIcon(Icons.check_circle), findsOneWidget);
        expect(find.byIcon(Icons.add_circle_outline), findsNothing);
        expect(tester.takeException(), isNull);
      });
    });
  });
}
