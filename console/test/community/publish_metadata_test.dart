import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/community/publish.metadata.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media/media.pb.dart';
import 'package:retrovibed/media/media.known.pb.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

const _nonNilUuid = '12345678-1234-1234-1234-123456789012';

Download _download({
  String description = 'Test Content',
  String mimetype = 'video/mp4',
  String knownMediaId = '',
  String torrentId = '',
}) => Download(
  media: Media(
    description: description,
    mimetype: mimetype,
    knownMediaId: knownMediaId,
    torrentId: torrentId,
  ),
);

Future<KnownLookupResponse> _noopKnownGet(String id, {List<httpx.Option> options = const []}) =>
    Future.value(KnownLookupResponse(known: Known()));

Future<KnownCreateResponse> _defaultKnownCreate(KnownCreateRequest req,
        {List<httpx.Option> options = const []}) =>
    Future.value(KnownCreateResponse(known: Known(id: _nonNilUuid, description: req.known.description)));

Future<MetadataSyncResponse> _noopMetadataSync(String id, Media media,
        {List<httpx.Option> options = const []}) =>
    Future.value(MetadataSyncResponse());

PublishMetadata _build({
  Download? download,
  void Function(Known)? onConfirm,
  Future<KnownLookupResponse> Function(String, {List<httpx.Option> options})? knownGet,
  Future<KnownCreateResponse> Function(KnownCreateRequest, {List<httpx.Option> options})? knownCreate,
  Future<MetadataSyncResponse> Function(String, Media, {List<httpx.Option> options})? metadataSync,
}) => PublishMetadata(
  download: download ?? _download(),
  onConfirm: onConfirm ?? (_) {},
  knownGet: knownGet ?? _noopKnownGet,
  knownCreate: knownCreate ?? _defaultKnownCreate,
  metadataSync: metadataSync ?? _noopMetadataSync,
);

void main() {
  group('PublishMetadata', () {
    testWidgets('renders form fields', (tester) async {
      await tester.pumpApp(_build());
      await tester.pumpAndSettle();

      expect(find.text('Title'), findsOneWidget);
      expect(find.text('Summary'), findsOneWidget);
      expect(find.text('Image URL (optional)'), findsOneWidget);
      expect(find.text('Release Date'), findsOneWidget);
      expect(find.text('Adult content'), findsOneWidget);
      expect(find.text('Continue'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('populates title from download description', (tester) async {
      await tester.pumpApp(_build(download: _download(description: 'My Movie')));
      await tester.pumpAndSettle();

      expect(find.widgetWithText(TextFormField, 'My Movie'), findsOneWidget);
    });

    testWidgets('loads existing metadata when knownMediaId is set', (tester) async {
      bool knownGetCalled = false;

      await tester.pumpApp(
        _build(
          download: _download(knownMediaId: _nonNilUuid),
          knownGet: (id, {options = const []}) {
            knownGetCalled = true;
            return Future.value(KnownLookupResponse(known: Known(id: id, description: 'Loaded Title')));
          },
        ),
      );
      await tester.pumpAndSettle();

      expect(knownGetCalled, isTrue);
      expect(find.widgetWithText(TextFormField, 'Loaded Title'), findsOneWidget);
    });

    testWidgets('submit with empty title shows error and does not call onConfirm',
        (tester) async {
      Known? confirmed;

      await tester.pumpApp(_build(onConfirm: (k) => confirmed = k));
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextFormField).first, '');
      await tester.pump();
      await tester.tap(find.text('Continue'));
      await tester.pumpAndSettle();

      expect(confirmed, isNull);
      expect(find.text('an unexpected problem has occurred'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('submit without existing metadata: calls knownCreate then metadataSync',
        (tester) async {
      bool knownCreateCalled = false;
      bool metadataSyncCalled = false;
      Known? confirmed;

      await tester.pumpApp(
        _build(
          download: _download(torrentId: _nonNilUuid),
          onConfirm: (k) => confirmed = k,
          knownCreate: (req, {options = const []}) {
            knownCreateCalled = true;
            return Future.value(KnownCreateResponse(known: Known(id: _nonNilUuid, description: req.known.description)));
          },
          metadataSync: (id, media, {options = const []}) {
            metadataSyncCalled = true;
            return Future.value(MetadataSyncResponse());
          },
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('Continue'));
      await tester.pumpAndSettle();

      expect(knownCreateCalled, isTrue);
      expect(metadataSyncCalled, isTrue);
      expect(confirmed, isNotNull);
      expect(tester.takeException(), isNull);
    });

    testWidgets('submit with existing metadata: calls metadataSync only',
        (tester) async {
      bool knownCreateCalled = false;
      bool metadataSyncCalled = false;
      Known? confirmed;

      await tester.pumpApp(
        _build(
          download: _download(knownMediaId: _nonNilUuid, torrentId: _nonNilUuid),
          onConfirm: (k) => confirmed = k,
          knownGet: (id, {options = const []}) =>
              Future.value(KnownLookupResponse(known: Known(id: id, description: 'Existing'))),
          knownCreate: (req, {options = const []}) {
            knownCreateCalled = true;
            return Future.value(KnownCreateResponse(known: Known()));
          },
          metadataSync: (id, media, {options = const []}) {
            metadataSyncCalled = true;
            return Future.value(MetadataSyncResponse());
          },
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('Continue'));
      await tester.pumpAndSettle();

      expect(knownCreateCalled, isFalse);
      expect(metadataSyncCalled, isTrue);
      expect(confirmed, isNotNull);
      expect(tester.takeException(), isNull);
    });

    testWidgets('submit without torrent: calls onConfirm directly, skips metadataSync',
        (tester) async {
      bool metadataSyncCalled = false;
      Known? confirmed;

      await tester.pumpApp(
        _build(
          download: _download(torrentId: ''),
          onConfirm: (k) => confirmed = k,
          metadataSync: (id, media, {options = const []}) {
            metadataSyncCalled = true;
            return Future.value(MetadataSyncResponse());
          },
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('Continue'));
      await tester.pumpAndSettle();

      expect(metadataSyncCalled, isFalse);
      expect(confirmed, isNotNull);
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
