import 'dart:async';
import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/community/content.detail.dart';
import 'package:retrovibed/community/content.display.dart';
import 'package:retrovibed/community/community.pb.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

Community _testCommunity() => Community(
  id: 'community-1',
  domain: 'testdomain',
  description: 'A test community',
);

Future<PublishedContentListResponse> _empty(
  String id, {
  List<httpx.Option> options = const [],
  int offset = 0,
  int limit = 100,
}) async {
  return PublishedContentListResponse(
    items: [],
    next: PublishedContentListRequest(offset: Int64(offset), limit: Int64(limit)),
  );
}

Future<PublishedContentListResponse> _noop(
  String id, {
  List<httpx.Option> options = const [],
  int offset = 0,
  int limit = 100,
}) => Completer<PublishedContentListResponse>().future;

Future<PublishedContentListResponse> _withContent(
  String id, {
  List<httpx.Option> options = const [],
  int offset = 0,
  int limit = 100,
}) async {
  return PublishedContentListResponse(
    next: PublishedContentListRequest(offset: Int64(offset), limit: Int64(limit)),
    items: [
      PublishedContent(
        id: 'pc-1',
        communityId: 'community-1',
        knownMediaId: 'magnet:?xt=urn:btih:abc123',
        magnetUri: 'magnet:?xt=urn:btih:abc123',
        title: 'Movie One',
        mimetype: 'video/mp4',
        bytes: Int64(1500000000),
        publishedAt: '2026-01-15T10:00:00Z',
      ),
      PublishedContent(
        id: 'pc-2',
        communityId: 'community-1',
        knownMediaId: 'magnet:?xt=urn:btih:def456',
        magnetUri: 'magnet:?xt=urn:btih:def456',
        title: 'Movie Two',
        mimetype: 'video/mp4',
        bytes: Int64(800000000),
        publishedAt: '2026-02-10T14:30:00Z',
      ),
    ],
  );
}

Future<PublishedContentListResponse> _withArchivedContent(
  String id, {
  List<httpx.Option> options = const [],
  int offset = 0,
  int limit = 100,
}) async {
  return PublishedContentListResponse(
    next: PublishedContentListRequest(offset: Int64(offset), limit: Int64(limit)),
    items: [
      PublishedContent(
        id: 'pc-1',
        communityId: 'community-1',
        knownMediaId: 'magnet:?xt=urn:btih:abc123',
        magnetUri: 'magnet:?xt=urn:btih:abc123',
        title: 'Archived Movie',
        mimetype: 'video/mp4',
        bytes: Int64(1500000000),
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
        bytes: Int64(800000000),
        publishedAt: '2026-02-10T14:30:00Z',
        archivedId: uuidx.min(),
      ),
    ],
  );
}

Future<PublishedContentListResponse> _withLongContent(
  String id, {
  List<httpx.Option> options = const [],
  int offset = 0,
  int limit = 100,
}) async {
  return PublishedContentListResponse(
    next: PublishedContentListRequest(offset: Int64(0), limit: Int64(limit)),
    items: [
      PublishedContent(
        id: 'pc-long',
        communityId: 'community-1',
        knownMediaId: 'magnet:?xt=urn:btih:${'x' * 40}',
        magnetUri: 'magnet:?xt=urn:btih:${'x' * 40}',
        title: 'A' * 200,
        mimetype: 'video/mp4',
        bytes: Int64(1500000000),
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

      testWidgets('detail is offstage before tap', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          CommunityContentDisplay(
            community: _testCommunity(),
            apipublished: _withContent,
          ),
        );
        await tester.pumpAndSettle();

        // with maintainState the widget is in tree but offstage - skipOffstage: false finds it
        expect(find.byType(PublishedContentDetail, skipOffstage: false), findsWidgets);
        // default skipOffstage: true should find nothing visible
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
  });
}
