import 'package:flutter/gestures.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/library/known.media.display.dart';
import 'package:retrovibed/library/known.media.source.dart';
import 'package:retrovibed/library/api.dart' as api;
import 'package:retrovibed/media/media.pb.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = ValueVariant(Resolutions.all.entries.toSet());

Media _media({String mimetype = 'video/mp4'}) => Media(
  id: uuidx.withSuffix(1),
  description: 'Test Media',
  mimetype: mimetype,
  createdAt: DateTime.now().toIso8601String(),
  archiveId: uuidx.min(),
  torrentId: uuidx.min(),
  knownMediaId: uuidx.min(),
);

Future<void> _hover(WidgetTester tester) async {
  final gesture = await tester.createGesture(kind: PointerDeviceKind.mouse);
  await gesture.addPointer(location: Offset.zero);
  addTearDown(gesture.removePointer);
  await gesture.moveTo(tester.getCenter(find.byType(KnownMediaDisplay)));
  await tester.pumpAndSettle();
}

void _resolutionTests(String description, Media media, [api.Known? known]) {
  Future<api.Known> pending(Media m) =>
      known != null ? Future.value(known) : api.known.autodetect(m);

  group(description, () {
    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;

      await tester.pumpApp(physicalSize: entry.value, KnownMediaDisplay(pending(media), media: media));
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);

      await _hover(tester);
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders without overflow when archive pending', (
      WidgetTester tester,
    ) async {
      final entry = _resolutions.currentValue!;

      final m = media.deepCopy()..archiveId = uuidx.max();
      await tester.pumpApp(physicalSize: entry.value, KnownMediaDisplay(pending(m), media: m));
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);

      await _hover(tester);
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders without overflow when archive active', (
      WidgetTester tester,
    ) async {
      final entry = _resolutions.currentValue!;

      final m = media.deepCopy()..archiveId = uuidx.withSuffix(99);
      await tester.pumpApp(physicalSize: entry.value, KnownMediaDisplay(pending(m), media: m));
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);

      await _hover(tester);
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('renders without overflow when highlighted', (
      WidgetTester tester,
    ) async {
      final entry = _resolutions.currentValue!;

      await tester.pumpApp(
        physicalSize: entry.value,
        KnownMediaDisplay(pending(media), media: media, highlighted: true),
      );
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);

      await _hover(tester);
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });
}

api.Known _known({String id = ''}) => api.Known(
  id: id,
  description: 'Test Known',
  summary: 'Test summary',
  rating: 0.0,
  image: '',
);

void main() {
  _resolutionTests('no source icon', _media());
  _resolutionTests('tmdb source icon', _media(), _known(id: '${KnownMediaSource.tmdb}-0000-0000-0000-000000000000'));

  group('KnownMediaDisplay onTap', () {
    testWidgets('onTap is invoked when tapped', (tester) async {
      final media = _media();
      var tapped = false;
      await tester.pumpApp(
        KnownMediaDisplay(
          api.known.autodetect(media),
          media: media,
          onTap: () { tapped = true; },
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(KnownMediaDisplay));
      await tester.pump();
      expect(tapped, isTrue);
    });

    testWidgets('missing factory forwards onTap', (tester) async {
      final media = _media();
      var tapped = false;
      await tester.pumpApp(
        KnownMediaDisplay.missing(
          media,
          onTap: () { tapped = true; },
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(KnownMediaDisplay));
      await tester.pump();
      expect(tapped, isTrue);
    });
  });

  group('KnownMediaDisplay onDoubleTap', () {
    testWidgets('onDoubleTap is invoked when double tapped', (tester) async {
      final media = _media();
      var doubleTapped = false;
      await tester.pumpApp(
        KnownMediaDisplay(
          api.known.autodetect(media),
          media: media,
          onDoubleTap: () { doubleTapped = true; },
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(KnownMediaDisplay));
      await tester.pump(kDoubleTapMinTime);
      await tester.tap(find.byType(KnownMediaDisplay));
      await tester.pump(kDoubleTapTimeout);
      expect(doubleTapped, isTrue);
    });

    testWidgets('missing factory forwards onDoubleTap', (tester) async {
      final media = _media();
      var doubleTapped = false;
      await tester.pumpApp(
        KnownMediaDisplay.missing(
          media,
          onDoubleTap: () { doubleTapped = true; },
        ),
      );
      await tester.pumpAndSettle();
      await tester.tap(find.byType(KnownMediaDisplay));
      await tester.pump(kDoubleTapMinTime);
      await tester.tap(find.byType(KnownMediaDisplay));
      await tester.pump(kDoubleTapTimeout);
      expect(doubleTapped, isTrue);
    });
  });
}
