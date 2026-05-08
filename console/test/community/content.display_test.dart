import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/community/content.display.dart';
import 'package:retrovibed/community/community.pb.dart';
import 'package:retrovibed/media/media.pb.dart' as media_pb;
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
  return PublishedContentListResponse(items: []);
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
    items: [
      PublishedContent(
        id: 'pc-1',
        communityId: 'community-1',
        knownMediaId: 'magnet:?xt=urn:btih:abc123&dn=Movie+One',
        magnetUri: 'magnet:?xt=urn:btih:abc123&dn=Movie+One',
      ),
      PublishedContent(
        id: 'pc-2',
        communityId: 'community-1',
        knownMediaId: 'magnet:?xt=urn:btih:def456&dn=Movie+Two',
        magnetUri: 'magnet:?xt=urn:btih:def456&dn=Movie+Two',
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
  final longId = 'A' * 200;
  final longUri = 'magnet:?xt=urn:btih:${'x' * 160}&dn=${'Very+Long+Title+' * 10}';
  return PublishedContentListResponse(
    items: [
      PublishedContent(
        id: 'pc-long',
        communityId: 'community-1',
        knownMediaId: longId,
        magnetUri: longUri,
      ),
    ],
  );
}

Future<media_pb.MagnetCreateResponse> _magnetNoop(
  media_pb.MagnetCreateRequest req, {
  List<httpx.Option> options = const [],
}) async {
  return media_pb.MagnetCreateResponse();
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
            apimagnet: _magnetNoop,
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
            apimagnet: _magnetNoop,
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
            apimagnet: _magnetNoop,
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
            apimagnet: _magnetNoop,
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
              apimagnet: _magnetNoop,
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
              apimagnet: _magnetNoop,
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
              apimagnet: _magnetNoop,
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
              apimagnet: _magnetNoop,
            ),
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('No content published yet'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    group('content display', () {
      testWidgets('shows Published Content header', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          CommunityContentDisplay(
            community: _testCommunity(),
            apipublished: _empty,
            apimagnet: _magnetNoop,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('Published Content'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('shows Publish button', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          CommunityContentDisplay(
            community: _testCommunity(),
            apipublished: _empty,
            apimagnet: _magnetNoop,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.text('Publish'), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('shows content rows after loading', (tester) async {
        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          CommunityContentDisplay(
            community: _testCommunity(),
            apipublished: _withContent,
            apimagnet: _magnetNoop,
          ),
        );
        await tester.pumpAndSettle();

        expect(find.byIcon(Icons.movie), findsNWidgets(2));
        expect(find.byIcon(Icons.download), findsNWidgets(2));
        expect(tester.takeException(), isNull);
      });

      testWidgets('archive button calls apimagnet', (tester) async {
        var magnetCalled = false;
        Future<media_pb.MagnetCreateResponse> trackMagnet(
          media_pb.MagnetCreateRequest req, {
          List<httpx.Option> options = const [],
        }) async {
          magnetCalled = true;
          return media_pb.MagnetCreateResponse();
        }

        await tester.pumpApp(
          physicalSize: const Size(1280, 720),
          CommunityContentDisplay(
            community: _testCommunity(),
            apipublished: _withContent,
            apimagnet: trackMagnet,
          ),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.byIcon(Icons.download).first);
        await tester.pumpAndSettle();

        expect(magnetCalled, isTrue);
        expect(tester.takeException(), isNull);
      });
    });
  });
}
