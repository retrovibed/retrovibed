import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/community/content.detail.dart';
import 'package:retrovibed/community/content.display.read.only.dart';
import 'package:retrovibed/community.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

Community _testCommunity({String? lastSyncAt}) => Community(
  id: 'community-1',
  url: 'https://testdomain.community.retrovibe.space',
  description: 'A test community',
  lastSyncAt: lastSyncAt ?? DateTime.now().toUtc().toIso8601String(),
);

Future<PublishedContentSearchResponse> _empty(
  String id, {
  List<httpx.Option> options = const [],
  PublishedContentSearchRequest? req,
}) async {
  final r = req ?? PublishedContentSearchRequest(offset: ds.Int64(0), limit: ds.Int64(100));
  return PublishedContentSearchResponse(items: [], next: r);
}

Future<PublishedContentSearchResponse> _noop(
  String id, {
  List<httpx.Option> options = const [],
  PublishedContentSearchRequest? req,
}) => Completer<PublishedContentSearchResponse>().future;

Future<PublishedContentSearchResponse> _withContent(
  String id, {
  List<httpx.Option> options = const [],
  PublishedContentSearchRequest? req,
}) async {
  final r = req ?? PublishedContentSearchRequest(offset: ds.Int64(0), limit: ds.Int64(100));
  return PublishedContentSearchResponse(
    next: r,
    items: [
      PublishedContent(
        id: 'pc-1',
        communityId: 'community-1',
        title: 'Movie One',
        mimetype: 'video/mp4',
        bytes: ds.Int64(1500000000),
        publishedAt: '2026-01-15T10:00:00Z',
      ),
      PublishedContent(
        id: 'pc-2',
        communityId: 'community-1',
        title: 'Movie Two',
        mimetype: 'video/mp4',
        bytes: ds.Int64(800000000),
        publishedAt: '2026-02-10T14:30:00Z',
      ),
    ],
  );
}

Future<PublishedContentSearchResponse> _withLongContent(
  String id, {
  List<httpx.Option> options = const [],
  PublishedContentSearchRequest? req,
}) async {
  final r = req ?? PublishedContentSearchRequest(offset: ds.Int64(0), limit: ds.Int64(100));
  return PublishedContentSearchResponse(
    next: r,
    items: [
      PublishedContent(
        id: 'pc-long',
        communityId: 'community-1',
        title: 'A' * 200,
        mimetype: 'video/mp4',
        bytes: ds.Int64(1500000000),
        publishedAt: '2026-01-15T10:00:00Z',
      ),
    ],
  );
}

Future<PublishedContentSearchResponse> _404(
  String id, {
  List<httpx.Option> options = const [],
  PublishedContentSearchRequest? req,
}) => Future.error(http.Response('not found', 404));

void main() {
  group('ContentDisplayReadOnly', () {
    group('resolutions - loading state', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          ContentDisplayReadOnly(
            community: _testCommunity(),
            apipublished: _noop,
          ),
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
          ContentDisplayReadOnly(
            community: _testCommunity(),
            apipublished: _empty,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('No content published yet'), findsOneWidget);
        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('resolutions - with content', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          ContentDisplayReadOnly(
            community: _testCommunity(),
            apipublished: _withContent,
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('resolutions - long content text overflow', () {
      testWidgets('ellipsizes without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          ContentDisplayReadOnly(
            community: _testCommunity(),
            apipublished: _withLongContent,
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      }, variant: _resolutions);
    });

    group('constrained parent', () {
      testWidgets('renders inside tight SizedBox without overflow', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          SizedBox(
            width: 320,
            height: 480,
            child: ContentDisplayReadOnly(
              community: _testCommunity(),
              apipublished: _withContent,
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('empty state inside tight SizedBox', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          SizedBox(
            width: 320,
            height: 480,
            child: ContentDisplayReadOnly(
              community: _testCommunity(),
              apipublished: _empty,
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('No content published yet'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('unconstrained parent', () {
      testWidgets('renders inside SingleChildScrollView without overflow', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          SingleChildScrollView(
            child: ContentDisplayReadOnly(
              community: _testCommunity(),
              apipublished: _withContent,
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(tester.takeException(), isNull);
      });

      testWidgets('empty state inside SingleChildScrollView', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          SingleChildScrollView(
            child: ContentDisplayReadOnly(
              community: _testCommunity(),
              apipublished: _empty,
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('No content published yet'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('content display', () {
      testWidgets('shows content rows after loading', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          ContentDisplayReadOnly(
            community: _testCommunity(),
            apipublished: _withContent,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('Movie One'), findsOneWidget);
        expect(find.text('Movie Two'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('no delete buttons shown', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          ContentDisplayReadOnly(
            community: _testCommunity(),
            apipublished: _withContent,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byIcon(Icons.delete), findsNothing);
        expect(tester.takeException(), isNull);
      });

      testWidgets('detail is not built before tap', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          ContentDisplayReadOnly(
            community: _testCommunity(),
            apipublished: _withContent,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(PublishedContentDetail, skipOffstage: false), findsNothing);
        expect(tester.takeException(), isNull);
      });

      testWidgets('tapping row makes detail visible', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          ContentDisplayReadOnly(
            community: _testCommunity(),
            apipublished: _withContent,
          ),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.text('Movie One'));
        await tester.pump();

        expect(find.byType(PublishedContentDetail), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('tapping row twice collapses detail', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          ContentDisplayReadOnly(
            community: _testCommunity(),
            apipublished: _withContent,
          ),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.text('Movie One'));
        await tester.pump();
        expect(find.byType(PublishedContentDetail), findsOneWidget);

        await tester.tap(find.text('Movie One'));
        await tester.pump();
        expect(find.byType(PublishedContentDetail), findsNothing);
        expect(tester.takeException(), isNull);
      });
    });

    group('404 response', () {
      testWidgets('shows empty state without error', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          ContentDisplayReadOnly(
            community: _testCommunity(),
            apipublished: _404,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('No content published yet'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('staleness / auto-resync', () {
      testWidgets('stale community triggers apiresync on mount, not apipublished', (tester) async {
        bool apipublishedCalled = false;
        Future<PublishedContentSearchResponse> failPublished(
          String id, {
          List<httpx.Option> options = const [],
          PublishedContentSearchRequest? req,
        }) async {
          apipublishedCalled = true;
          throw StateError('apipublished should not be called for a stale community');
        }

        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          ContentDisplayReadOnly(
            community: _testCommunity(
              lastSyncAt: DateTime.now().toUtc().subtract(const Duration(hours: 2)).toIso8601String(),
            ),
            apipublished: failPublished,
            apiresync: _withContent,
          ),
        );
        await tester.pumpAndSettle();

        expect(apipublishedCalled, isFalse);
        expect(find.text('Movie One'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('fresh community calls apipublished, not apiresync', (tester) async {
        bool apiresyncCalled = false;
        Future<PublishedContentSearchResponse> failResync(
          String id, {
          List<httpx.Option> options = const [],
          PublishedContentSearchRequest? req,
        }) async {
          apiresyncCalled = true;
          throw StateError('apiresync should not be called for a fresh community');
        }

        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          ContentDisplayReadOnly(
            community: _testCommunity(
              lastSyncAt: DateTime.now().toUtc().subtract(const Duration(minutes: 10)).toIso8601String(),
            ),
            apipublished: _withContent,
            apiresync: failResync,
          ),
        );
        await tester.pumpAndSettle();

        expect(apiresyncCalled, isFalse);
        expect(find.text('Movie One'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('never-synced community (empty last_sync_at) triggers apiresync', (tester) async {
        bool apiresyncCalled = false;
        Future<PublishedContentSearchResponse> captureResync(
          String id, {
          List<httpx.Option> options = const [],
          PublishedContentSearchRequest? req,
        }) {
          apiresyncCalled = true;
          return _withContent(id, options: options, req: req);
        }

        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          ContentDisplayReadOnly(
            community: _testCommunity(lastSyncAt: ''),
            apipublished: _noop,
            apiresync: captureResync,
          ),
        );
        await tester.pumpAndSettle();

        expect(apiresyncCalled, isTrue);
        expect(find.text('Movie One'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });
  });
}
