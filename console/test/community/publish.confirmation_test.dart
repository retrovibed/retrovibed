import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/community.dart';
import 'package:retrovibed/community/publish.confirmation.dart';
import 'package:retrovibed/media/media.pb.dart';
import 'package:retrovibed/media/media.known.pb.dart';

final _resolutions = Resolutions.variant();

const _communityId = 'comm-id-123';
const _knownMediaId = 'known-id-456';
const _libraryId = 'lib-id-789';

Download _videoDownload({String knownMediaId = _knownMediaId}) => Download(
  media: Media(
    id: _libraryId,
    description: 'My Video',
    mimetype: 'video/mp4',
    knownMediaId: knownMediaId,
  ),
);

Download _audioDownload() => Download(
  media: Media(
    id: _libraryId,
    description: 'My Audio',
    mimetype: 'audio/mp3',
  ),
);

Community _community() => Community(id: _communityId, url: 'https://test.community');

Future<YouTubeStatus> _noYouTube({List<httpx.Option> options = const []}) => Future.value(YouTubeStatus(id: ''));

Future<PublishContentResponse> _noopPublish(
  String communityId,
  PublishContentRequest req, {
  List<httpx.Option> options = const [],
}) => Future.value(PublishContentResponse());

PublishConfirmation _build({
  Download? download,
  Community? community,
  Known? knownMedia,
  VoidCallback? onPublished,
  Future<YouTubeStatus> Function({List<httpx.Option> options})? youtubeStatus,
  Future<PublishContentResponse> Function(String, PublishContentRequest, {List<httpx.Option> options})? publish,
}) => PublishConfirmation(
  download: download ?? _videoDownload(),
  community: community ?? _community(),
  knownMedia: knownMedia,
  onPublished: onPublished ?? () {},
  youtubeStatus: youtubeStatus ?? _noYouTube,
  apicommunitypublish: publish ?? _noopPublish,
);

void main() {
  group('PublishConfirmation', () {
    testWidgets('renders confirmation rows for content and community', (tester) async {
      await tester.pumpApp(_build());
      await tester.pumpAndSettle();

      expect(find.text('My Video'), findsOneWidget);
      expect(find.text('https://test.community'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders media info row when knownMedia is provided', (tester) async {
      await tester.pumpApp(
        _build(knownMedia: Known(id: 'k', description: 'Known Title')),
      );
      await tester.pumpAndSettle();

      expect(find.text('Known Title'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('cross-post checkbox visible for video content', (tester) async {
      await tester.pumpApp(_build(download: _videoDownload()));
      await tester.pumpAndSettle();

      expect(find.text('Cross-post to YouTube'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('cross-post checkbox hidden for non-video content', (tester) async {
      await tester.pumpApp(_build(download: _audioDownload()));
      await tester.pumpAndSettle();

      expect(find.text('Cross-post to YouTube'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('YouTube status with id enables the cross-post checkbox', (tester) async {
      await tester.pumpApp(
        _build(
          youtubeStatus: ({options = const []}) => Future.value(YouTubeStatus(id: 'google-123')),
        ),
      );
      await tester.pumpAndSettle();

      final checkbox = tester.widget<CheckboxListTile>(find.byType(CheckboxListTile));
      expect(checkbox.onChanged, isNotNull);
    });

    testWidgets('YouTube status with empty id disables the cross-post checkbox', (tester) async {
      await tester.pumpApp(
        _build(
          youtubeStatus: ({options = const []}) => Future.value(YouTubeStatus(id: '')),
        ),
      );
      await tester.pumpAndSettle();

      final checkbox = tester.widget<CheckboxListTile>(find.byType(CheckboxListTile));
      expect(checkbox.onChanged, isNull);
    });

    testWidgets('publish invoked with correct communityId, knownMediaId, libraryId', (tester) async {
      String? capturedCommunityId;
      PublishContentRequest? capturedRequest;

      await tester.pumpApp(
        _build(
          publish: (communityId, req, {options = const []}) {
            capturedCommunityId = communityId;
            capturedRequest = req;
            return Future.value(PublishContentResponse());
          },
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('Publish'));
      await tester.pumpAndSettle();

      expect(capturedCommunityId, equals(_communityId));
      expect(capturedRequest, isNotNull);
      expect(capturedRequest!.publishedContent.communityId, equals(_communityId));
      expect(capturedRequest!.publishedContent.knownMediaId, equals(_knownMediaId));
      expect(capturedRequest!.publishedContent.libraryId, equals(_libraryId));
      expect(tester.takeException(), isNull);
    });

    testWidgets('onPublished called after successful publish', (tester) async {
      bool published = false;

      await tester.pumpApp(
        _build(
          onPublished: () => published = true,
          publish: (id, req, {options = const []}) => Future.value(PublishContentResponse()),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('Publish'));
      await tester.pumpAndSettle();

      expect(published, isTrue);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders without overflow', (tester) async {
      final entry = _resolutions.currentValue!;

      await tester.pumpApp(
        physicalSize: entry.value,
        SingleChildScrollView(child: _build()),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });
}
