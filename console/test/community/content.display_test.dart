import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/community/content.detail.dart';
import 'package:retrovibed/community/content.display.dart';
import 'package:retrovibed/community.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/modals.dart' as modals;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

Community _testCommunity({String? lastSyncAt}) => Community(
  id: 'community-1',
  domain: 'testdomain',
  description: 'A test community',
  lastSyncAt: lastSyncAt ?? DateTime.now().toUtc().toIso8601String(),
);

Future<PublishedContentSearchResponse> _empty(
  String id, {
  List<httpx.Option> options = const [],
  PublishedContentSearchRequest? req,
}) async {
  final r = req ?? PublishedContentSearchRequest(offset: ds.Int64(0), limit: ds.Int64(100));
  return PublishedContentSearchResponse(
    items: [],
    next: r,
  );
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
        knownMediaId: 'magnet:?xt=urn:btih:abc123',
        magnetUri: 'magnet:?xt=urn:btih:abc123',
        title: 'Movie One',
        mimetype: 'video/mp4',
        bytes: ds.Int64(1500000000),
        publishedAt: '2026-01-15T10:00:00Z',
      ),
      PublishedContent(
        id: 'pc-2',
        communityId: 'community-1',
        knownMediaId: 'magnet:?xt=urn:btih:def456',
        magnetUri: 'magnet:?xt=urn:btih:def456',
        title: 'Movie Two',
        mimetype: 'video/mp4',
        bytes: ds.Int64(800000000),
        publishedAt: '2026-02-10T14:30:00Z',
      ),
    ],
  );
}

Future<PublishedContentSearchResponse> _withArchivedContent(
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
        knownMediaId: 'magnet:?xt=urn:btih:abc123',
        magnetUri: 'magnet:?xt=urn:btih:abc123',
        title: 'Archived Movie',
        mimetype: 'video/mp4',
        bytes: ds.Int64(1500000000),
        publishedAt: '2026-01-15T10:00:00Z',
        archivedId: uuidx.withSuffix(1),
      ),
      PublishedContent(
        id: 'pc-2',
        communityId: 'community-1',
        knownMediaId: 'magnet:?xt=urn:btih:def456',
        magnetUri: 'magnet:?xt=urn:btih:def456',
        title: 'Unarchived Movie',
        mimetype: 'video/mp4',
        bytes: ds.Int64(800000000),
        publishedAt: '2026-02-10T14:30:00Z',
        archivedId: uuidx.min(),
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
        knownMediaId: 'magnet:?xt=urn:btih:${'x' * 40}',
        magnetUri: 'magnet:?xt=urn:btih:${'x' * 40}',
        title: 'A' * 200,
        mimetype: 'video/mp4',
        bytes: ds.Int64(1500000000),
        publishedAt: '2026-01-15T10:00:00Z',
      ),
    ],
  );
}

void main() {
  group('CommunityContentDisplay', () {
    group('resolutions - loading state', () {
      testWidgets('renders without overflow', (WidgetTester tester) async {
        final entry = _resolutions.currentValue!;
        await tester.pumpApp(
          physicalSize: entry.value,
          CommunityContentDisplay(
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
          CommunityContentDisplay(
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
          CommunityContentDisplay(
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
          CommunityContentDisplay(
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
            child: CommunityContentDisplay(
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
            child: CommunityContentDisplay(
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
            child: CommunityContentDisplay(
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
            child: CommunityContentDisplay(
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
          CommunityContentDisplay(
            community: _testCommunity(),
            apipublished: _withContent,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('Movie One'), findsOneWidget);
        expect(find.text('Movie Two'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('detail is not built before tap', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          CommunityContentDisplay(
            community: _testCommunity(),
            apipublished: _withContent,
          ),
        );
        await tester.pumpAndSettle();

        // lazily built: not present in the tree at all until first expanded,
        // so it doesn't fire its API call for every row on initial load.
        expect(find.byType(PublishedContentDetail, skipOffstage: false), findsNothing);
        expect(find.byType(PublishedContentDetail), findsNothing);
        expect(tester.takeException(), isNull);
      });

      testWidgets('tapping row title makes detail visible', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          CommunityContentDisplay(
            community: _testCommunity(),
            apipublished: _withContent,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byType(PublishedContentDetail), findsNothing);

        await tester.tap(find.text('Movie One'));
        await tester.pump();

        expect(find.byType(PublishedContentDetail), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('tapping row twice collapses detail', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          CommunityContentDisplay(
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

      testWidgets('archived content shows archive icon', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          CommunityContentDisplay(
            community: _testCommunity(),
            apipublished: _withArchivedContent,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byIcon(Icons.archive), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('deletion', () {
      Future<PublishContentDeleteResponse> tombstone(
        String id, {
        List<httpx.Option> options = const [],
      }) async => PublishContentDeleteResponse();

      testWidgets('delete button shown for each item', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          modals.Node(
            CommunityContentDisplay(
              community: _testCommunity(),
              apipublished: _withContent,
              apitombstone: tombstone,
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byIcon(Icons.delete), findsNWidgets(2));
        expect(tester.takeException(), isNull);
      });

      testWidgets('tapping delete shows confirmation dialog', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          modals.Node(
            CommunityContentDisplay(
              community: _testCommunity(),
              apipublished: _withContent,
              apitombstone: tombstone,
            ),
          ),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.byIcon(Icons.delete).first);
        await tester.pumpAndSettle();

        expect(find.text('Yes'), findsOneWidget);
        expect(find.text('No'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('confirming deletion removes item from list', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          modals.Node(
            CommunityContentDisplay(
              community: _testCommunity(),
              apipublished: _withContent,
              apitombstone: tombstone,
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('Movie One'), findsOneWidget);

        await tester.tap(find.byIcon(Icons.delete).first);
        await tester.pumpAndSettle();

        await tester.tap(find.text('Yes'));
        await tester.pumpAndSettle();

        expect(find.text('Movie One'), findsNothing);
        expect(find.text('Movie Two'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('cancelling deletion leaves item in list', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          modals.Node(
            CommunityContentDisplay(
              community: _testCommunity(),
              apipublished: _withContent,
              apitombstone: tombstone,
            ),
          ),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.byIcon(Icons.delete).first);
        await tester.pumpAndSettle();

        await tester.tap(find.text('No'));
        await tester.pumpAndSettle();

        expect(find.text('Movie One'), findsOneWidget);
        expect(find.text('Movie Two'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('tombstone is called with item id on confirm', (tester) async {
        String? calledWith;
        Future<PublishContentDeleteResponse> captureTombstone(
          String id, {
          List<httpx.Option> options = const [],
        }) async {
          calledWith = id;
          return PublishContentDeleteResponse();
        }

        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          modals.Node(
            CommunityContentDisplay(
              community: _testCommunity(),
              apipublished: _withContent,
              apitombstone: captureTombstone,
            ),
          ),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.byIcon(Icons.delete).first);
        await tester.pumpAndSettle();

        await tester.tap(find.text('Yes'));
        await tester.pumpAndSettle();

        expect(calledWith, 'pc-1');
        expect(tester.takeException(), isNull);
      });

      testWidgets('tombstone not called on cancel', (tester) async {
        bool called = false;
        Future<PublishContentDeleteResponse> captureTombstone(
          String id, {
          List<httpx.Option> options = const [],
        }) async {
          called = true;
          return PublishContentDeleteResponse();
        }

        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          modals.Node(
            CommunityContentDisplay(
              community: _testCommunity(),
              apipublished: _withContent,
              apitombstone: captureTombstone,
            ),
          ),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.byIcon(Icons.delete).first);
        await tester.pumpAndSettle();

        await tester.tap(find.text('No'));
        await tester.pumpAndSettle();

        expect(called, isFalse);
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
          CommunityContentDisplay(
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
          CommunityContentDisplay(
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
          CommunityContentDisplay(
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
