import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/deeplink.dart';
import 'package:retrovibed/community/api.dart' as community;
import 'package:retrovibed/billing/api.dart' as billing;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _community = community.Community(
  url: 'https://testdomain.community.retrovibe.space',
  description: 'A test community',
  createdAt: '2024-01-15T14:30:00Z',
);

final _searchResponse = community.CommunitySearchResponse(items: [_community]);
final _emptySearch = community.CommunitySearchResponse();
Future<Uri?> _noInitial() => Future.value(null);

void main() {
  group('DeepLink', () {
    group('community URL', () {
      testWidgets('shows subscribe confirmation', (tester) async {
        final controller = StreamController<Uri>();

        await tester.pumpApp(
          DeepLink(
            const Text('child'),
            uriStream: () => controller.stream,
            initialUri: _noInitial,
            search: (req, {options = const []}) => Future.value(_searchResponse),
            consumeAttribution: (token, {options = const []}) => Future.value(billing.AttributionConsumeResponse()),
            subscribe: (ctx, c, a) => Future.value(),
          ),
        );
        await tester.pumpAndSettle();

        controller.add(Uri.parse('https://testdomain.community.retrovibe.space'));
        await tester.pumpAndSettle();

        expect(find.text('Yes'), findsOneWidget);
        expect(find.text('No'), findsOneWidget);
        expect(find.text('https://testdomain.community.retrovibe.space'), findsOneWidget);
        expect(tester.takeException(), isNull);

        controller.close();
      });

      testWidgets('does not consume attribution', (tester) async {
        final controller = StreamController<Uri>();
        bool consumed = false;

        await tester.pumpApp(
          DeepLink(
            const Text('child'),
            uriStream: () => controller.stream,
            initialUri: _noInitial,
            search: (req, {options = const []}) => Future.value(_searchResponse),
            consumeAttribution: (token, {options = const []}) {
              consumed = true;
              return Future.value(billing.AttributionConsumeResponse());
            },
            subscribe: (ctx, c, a) => Future.value(),
          ),
        );
        await tester.pumpAndSettle();

        controller.add(Uri.parse('https://testdomain.community.retrovibe.space'));
        await tester.pumpAndSettle();

        expect(consumed, isFalse);
        expect(tester.takeException(), isNull);

        controller.close();
      });

      testWidgets('searches by domain from URI host', (tester) async {
        final controller = StreamController<Uri>();
        String? capturedQuery;

        await tester.pumpApp(
          DeepLink(
            const Text('child'),
            uriStream: () => controller.stream,
            initialUri: _noInitial,
            search: (req, {options = const []}) {
              capturedQuery = req.query;
              return Future.value(_searchResponse);
            },
            consumeAttribution: (token, {options = const []}) => Future.value(billing.AttributionConsumeResponse()),
            subscribe: (ctx, c, a) => Future.value(),
          ),
        );
        await tester.pumpAndSettle();

        controller.add(Uri.parse('https://mycommunity.community.retrovibe.space'));
        await tester.pumpAndSettle();

        expect(capturedQuery, equals('mycommunity'));
        expect(tester.takeException(), isNull);

        controller.close();
      });

      testWidgets('subscribes on confirm', (tester) async {
        final controller = StreamController<Uri>();
        community.Community? subscribed;

        await tester.pumpApp(
          DeepLink(
            const Text('child'),
            uriStream: () => controller.stream,
            initialUri: _noInitial,
            search: (req, {options = const []}) => Future.value(_searchResponse),
            consumeAttribution: (token, {options = const []}) => Future.value(billing.AttributionConsumeResponse()),
            subscribe: (ctx, c, a) {
              subscribed = c;
              return Future.value();
            },
          ),
        );
        await tester.pumpAndSettle();

        controller.add(Uri.parse('https://testdomain.community.retrovibe.space'));
        await tester.pumpAndSettle();

        await tester.tap(find.text('Yes'));
        await tester.pumpAndSettle();

        expect(subscribed, isNotNull);
        expect(subscribed!.url, equals('https://testdomain.community.retrovibe.space'));
        expect(tester.takeException(), isNull);

        controller.close();
      });

      testWidgets('dismisses on cancel', (tester) async {
        final controller = StreamController<Uri>();

        await tester.pumpApp(
          DeepLink(
            const Text('child'),
            uriStream: () => controller.stream,
            initialUri: _noInitial,
            search: (req, {options = const []}) => Future.value(_searchResponse),
            consumeAttribution: (token, {options = const []}) => Future.value(billing.AttributionConsumeResponse()),
            subscribe: (ctx, c, a) => Future.value(),
          ),
        );
        await tester.pumpAndSettle();

        controller.add(Uri.parse('https://testdomain.community.retrovibe.space'));
        await tester.pumpAndSettle();

        expect(find.text('Yes'), findsOneWidget);

        await tester.tap(find.text('No'));
        await tester.pumpAndSettle();

        expect(find.text('Yes'), findsNothing);
        expect(tester.takeException(), isNull);

        controller.close();
      });

      testWidgets('shows error when community not found', (tester) async {
        final controller = StreamController<Uri>();

        await tester.pumpApp(
          DeepLink(
            const Text('child'),
            uriStream: () => controller.stream,
            initialUri: _noInitial,
            search: (req, {options = const []}) => Future.value(_emptySearch),
            consumeAttribution: (token, {options = const []}) => Future.value(billing.AttributionConsumeResponse()),
            subscribe: (ctx, c, a) => Future.value(),
          ),
        );
        await tester.pumpAndSettle();

        controller.add(Uri.parse('https://unknown.community.retrovibe.space'));
        await tester.pumpAndSettle();

        expect(find.text('Yes'), findsNothing);
        expect(tester.takeException(), isNull);

        controller.close();
      });
    });

    group('invite URL', () {
      testWidgets('consumes attribution from invite URL', (tester) async {
        final controller = StreamController<Uri>();
        String? capturedToken;

        await tester.pumpApp(
          DeepLink(
            const Text('child'),
            uriStream: () => controller.stream,
            initialUri: _noInitial,
            search: (req, {options = const []}) => Future.value(_searchResponse),
            consumeAttribution: (token, {options = const []}) {
              capturedToken = token;
              return Future.value(billing.AttributionConsumeResponse());
            },
            subscribe: (ctx, c, a) => Future.value(),
          ),
        );
        await tester.pumpAndSettle();

        controller.add(Uri.parse('https://invite.retrovibe.space/?a=eyJinvite'));
        await tester.pumpAndSettle();

        expect(capturedToken, equals('eyJinvite'));
        expect(tester.takeException(), isNull);

        controller.close();
      });

      testWidgets('does not search community on invite URL', (tester) async {
        final controller = StreamController<Uri>();
        bool searched = false;

        await tester.pumpApp(
          DeepLink(
            const Text('child'),
            uriStream: () => controller.stream,
            initialUri: _noInitial,
            search: (req, {options = const []}) {
              searched = true;
              return Future.value(_searchResponse);
            },
            consumeAttribution: (token, {options = const []}) => Future.value(billing.AttributionConsumeResponse()),
            subscribe: (ctx, c, a) => Future.value(),
          ),
        );
        await tester.pumpAndSettle();

        controller.add(Uri.parse('https://invite.retrovibe.space/?a=tok'));
        await tester.pumpAndSettle();

        expect(searched, isFalse);
        expect(tester.takeException(), isNull);

        controller.close();
      });

      testWidgets('ignores invite URL with empty attribution', (tester) async {
        final controller = StreamController<Uri>();
        bool consumed = false;

        await tester.pumpApp(
          DeepLink(
            const Text('child'),
            uriStream: () => controller.stream,
            initialUri: _noInitial,
            search: (req, {options = const []}) => Future.value(_searchResponse),
            consumeAttribution: (token, {options = const []}) {
              consumed = true;
              return Future.value(billing.AttributionConsumeResponse());
            },
            subscribe: (ctx, c, a) => Future.value(),
          ),
        );
        await tester.pumpAndSettle();

        controller.add(Uri.parse('https://invite.retrovibe.space/'));
        await tester.pumpAndSettle();

        expect(consumed, isFalse);
        expect(tester.takeException(), isNull);

        controller.close();
      });
    });

    testWidgets('ignores non-retrovibe URLs', (tester) async {
      final controller = StreamController<Uri>();
      bool searchCalled = false;
      bool consumed = false;

      await tester.pumpApp(
        DeepLink(
          const Text('child'),
          uriStream: () => controller.stream,
          initialUri: _noInitial,
          search: (req, {options = const []}) {
            searchCalled = true;
            return Future.value(_searchResponse);
          },
          consumeAttribution: (token, {options = const []}) {
            consumed = true;
            return Future.value(billing.AttributionConsumeResponse());
          },
          subscribe: (ctx, c, a) => Future.value(),
        ),
      );
      await tester.pumpAndSettle();

      controller.add(Uri.parse('https://example.com/something'));
      await tester.pumpAndSettle();

      expect(searchCalled, isFalse);
      expect(consumed, isFalse);
      expect(tester.takeException(), isNull);

      controller.close();
    });

    testWidgets('processes initial URI', (tester) async {
      await tester.pumpApp(
        DeepLink(
          const Text('child'),
          uriStream: () => const Stream.empty(),
          initialUri: () => Future.value(Uri.parse('https://testdomain.community.retrovibe.space')),
          search: (req, {options = const []}) => Future.value(_searchResponse),
          consumeAttribution: (token, {options = const []}) => Future.value(billing.AttributionConsumeResponse()),
          subscribe: (ctx, c, a) => Future.value(),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('https://testdomain.community.retrovibe.space'), findsOneWidget);
      expect(find.text('Yes'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
