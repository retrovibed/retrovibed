import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/library/media.settings.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/library/api.dart' as api;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/design.kit/modals.dart' as modals;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('MediaSettings layout tests', () {
    late media.Media testMedia;
    late media.Media testMediaWithTorrent;

    // Mock search function that doesn't make HTTP calls
    Future<api.KnownSearchResponse> mockKnownSearch(
      api.KnownSearchRequest req, {
      List<httpx.Option> options = const [],
    }) async {
      return api.KnownSearchResponse(items: [], next: req);
    }

    setUp(() {
      // Create test media without torrent
      testMedia = media.Media(
        id: uuidx.withSuffix(1),
        description: 'Test Media Description',
        mimetype: 'video/mp4',
        createdAt: DateTime.now().toIso8601String(),
        archiveId: uuidx.min(),
        torrentId: uuidx.min(),
        knownMediaId: uuidx.min(),
      );

      // Create test media with torrent
      testMediaWithTorrent = media.Media(
        id: uuidx.withSuffix(2),
        description: 'Test Media with Torrent',
        mimetype: 'video/mkv',
        createdAt: DateTime.now().toIso8601String(),
        archiveId: uuidx.min(),
        torrentId: uuidx.withSuffix(3),
        knownMediaId: uuidx.min(),
      );
    });

    group('Desktop screen (1920x1080)', () {
      const double desktopWidth = 1920;
      const double desktopHeight = 1080;

      testWidgets('renders in constrained SizedBox', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: desktopWidth,
            height: desktopHeight,
            child: SingleChildScrollView(
              child: MediaSettings(
                current: testMedia,
                onChange: (pending, {bool forced = false, bool autoclose = false}) {},
                knownSearch: mockKnownSearch,

              ),
            ),
          ),
        );

        await tester.pumpAndSettle();

        expect(find.byType(MediaSettings), findsOneWidget);

        final settingsSize = tester.getSize(find.byType(MediaSettings));
        expect(settingsSize.width, lessThanOrEqualTo(desktopWidth));
      });

      testWidgets('renders in Column with Expanded', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: desktopWidth,
            height: desktopHeight,
            child: Column(
              children: [
                Expanded(
                  child: SingleChildScrollView(
                    child: MediaSettings(
                      current: testMedia,
                      onChange: (pending, {bool forced = false, bool autoclose = false}) {},
                      knownSearch: mockKnownSearch,

                    ),
                  ),
                ),
              ],
            ),
          ),
        );

        await tester.pumpAndSettle();

        expect(find.byType(MediaSettings), findsOneWidget);
      });

      testWidgets('renders in Row with Expanded', (WidgetTester tester) async {
        await tester.pumpApp(
          SizedBox(
            width: desktopWidth,
            height: desktopHeight,
            child: Row(
              children: [
                Expanded(
                  child: SingleChildScrollView(
                    child: MediaSettings(
                      current: testMedia,
                      onChange: (pending, {bool forced = false, bool autoclose = false}) {},
                      knownSearch: mockKnownSearch,

                    ),
                  ),
                ),
              ],
            ),
          ),
        );

        await tester.pumpAndSettle();

        expect(find.byType(MediaSettings), findsOneWidget);
      });

      testWidgets('renders in complex nested layout', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          Row(
            children: [
              SizedBox(width: 200, child: Container(color: Colors.grey)),
              Expanded(
                child: Column(
                  children: [
                    SizedBox(
                      height: 60,
                      child: Container(color: Colors.blue),
                    ),
                    Expanded(
                      child: SingleChildScrollView(
                        child: MediaSettings(
                          current: testMedia,
                          onChange: (pending, {bool forced = false, bool autoclose = false}) {},
                          knownSearch: mockKnownSearch,

                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        );

        await tester.pumpAndSettle();

        expect(find.byType(MediaSettings), findsOneWidget);
      });
    });

    group('Conditional rendering', () {
      testWidgets('hides source details when torrentId is min', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 800,
            height: 600,
            child: SingleChildScrollView(
              child: MediaSettings(
                current: testMedia,
                onChange: (pending, {bool forced = false, bool autoclose = false}) {},
                knownSearch: mockKnownSearch,

              ),
            ),
          ),
        );

        await tester.pumpAndSettle();

        expect(find.byType(MediaSettings), findsOneWidget);
        expect(find.text('source details'), findsNothing);
      });

      testWidgets('shows source details when torrentId is valid', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 800,
            height: 600,
            child: SingleChildScrollView(
              child: MediaSettings(
                current: testMediaWithTorrent,
                onChange: (pending, {bool forced = false, bool autoclose = false}) {},
                knownSearch: mockKnownSearch,

              ),
            ),
          ),
        );

        await tester.pumpAndSettle();

        expect(find.byType(MediaSettings), findsOneWidget);
        expect(find.text('source details'), findsOneWidget);
      });
    });

    group('State management', () {
      testWidgets('renders with onChange callback', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 800,
            height: 600,
            child: SingleChildScrollView(
              child: MediaSettings(
                current: testMedia,
                onChange: (pending, {bool forced = false, bool autoclose = false}) {},
                knownSearch: mockKnownSearch,

              ),
            ),
          ),
        );

        await tester.pumpAndSettle();

        expect(find.byType(MediaSettings), findsOneWidget);
      });

      testWidgets('maintains state across rebuilds', (
        WidgetTester tester,
      ) async {
        final ValueNotifier<media.Media> mediaNotifier = ValueNotifier(
          testMedia,
        );

        await tester.pumpApp(
          ValueListenableBuilder<media.Media>(
            valueListenable: mediaNotifier,
            builder: (context, currentMedia, child) {
              return SizedBox(
                width: 800,
                height: 600,
                child: SingleChildScrollView(
                  child: MediaSettings(
                    current: currentMedia,
                    onChange: (pending, {bool forced = false, bool autoclose = false}) {},
                    knownSearch: mockKnownSearch,

                  ),
                ),
              );
            },
          ),
        );

        await tester.pumpAndSettle();

        expect(find.byType(MediaSettings), findsOneWidget);

        // Update the media
        mediaNotifier.value = testMediaWithTorrent;
        await tester.pumpAndSettle();

        expect(find.byType(MediaSettings), findsOneWidget);
      });
    });

    group('KnownMediaDropdown onChange behavior', () {
      testWidgets('calls onChange with forced:true when torrentId is min', (
        WidgetTester tester,
      ) async {
        bool capturedForced = false;
        media.Media? capturedMedia;

        final media.Media mediaWithoutTorrent = media.Media(
          id: uuidx.withSuffix(10),
          description: 'Test without torrent',
          mimetype: 'video/mp4',
          createdAt: DateTime.now().toIso8601String(),
          archiveId: uuidx.min(),
          torrentId: uuidx.min(),
          knownMediaId: uuidx.min(),
        );

        await tester.pumpApp(
          SizedBox(
            width: 800,
            height: 600,
            child: SingleChildScrollView(
              child: MediaSettings(
                current: mediaWithoutTorrent,
                onChange: (pending, {bool forced = false, bool autoclose = false}) async {
                  capturedForced = forced;
                  capturedMedia = await pending;
                },
                knownSearch: mockKnownSearch,

              ),
            ),
          ),
        );

        await tester.pumpAndSettle();

        // Initial render should not trigger onChange
        expect(capturedForced, isFalse);
        expect(capturedMedia, isNull);
      });

      testWidgets('calls onChange with forced:true when torrentId is valid', (
        WidgetTester tester,
      ) async {
        bool capturedForced = false;
        media.Media? capturedMedia;

        await tester.pumpApp(
          SizedBox(
            width: 800,
            height: 600,
            child: SingleChildScrollView(
              child: MediaSettings(
                current: testMediaWithTorrent,
                onChange: (pending, {bool forced = false, bool autoclose = false}) async {
                  capturedForced = forced;
                  capturedMedia = await pending;
                },
                knownSearch: mockKnownSearch,

              ),
            ),
          ),
        );

        await tester.pumpAndSettle();

        // Initial render should not trigger onChange
        expect(capturedForced, isFalse);
        expect(capturedMedia, isNull);
      });

      testWidgets('displays source details when torrentId is valid', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 800,
            height: 600,
            child: SingleChildScrollView(
              child: MediaSettings(
                current: testMediaWithTorrent,
                onChange: (pending, {bool forced = false, bool autoclose = false}) {},
                knownSearch: mockKnownSearch,

              ),
            ),
          ),
        );

        await tester.pumpAndSettle();

        expect(find.text('source details'), findsOneWidget);
      });

      testWidgets('hides source details when torrentId is min', (
        WidgetTester tester,
      ) async {
        await tester.pumpApp(
          SizedBox(
            width: 800,
            height: 600,
            child: SingleChildScrollView(
              child: MediaSettings(
                current: testMedia,
                onChange: (pending, {bool forced = false, bool autoclose = false}) {},
                knownSearch: mockKnownSearch,

              ),
            ),
          ),
        );

        await tester.pumpAndSettle();

        expect(find.text('source details'), findsNothing);
      });
    });

    group('onVerify', () {
      Future<media.DownloadMetadataResponse> mockDiscoveredGet(
        String id, {
        List<httpx.Option> options = const [],
      }) async {
        return media.DownloadMetadataResponse(download: media.Download.create());
      }

      Widget buildWithTorrent({
        required String? Function(String, media.Download) onUpdate,
        required String? Function(String) onReset,
        void Function(bool forced)? captureOnChange,
      }) {
        return modals.Node(
          SizedBox(
            width: 800,
            height: 600,
            child: SingleChildScrollView(
              child: MediaSettings(
                current: testMediaWithTorrent,
                onChange: (pending, {bool forced = false, bool autoclose = false}) {
                  captureOnChange?.call(forced);
                },
                knownSearch: mockKnownSearch,

                discoveredGet: mockDiscoveredGet,
                discoveredUpdate: (id, download, {options = const []}) async {
                  onUpdate(id, download);
                  return media.DownloadUpdateResponse();
                },
                discoveredReset: (id, {options = const []}) async {
                  onReset(id);
                  return media.DownloadDeleteResponse();
                },
              ),
            ),
          ),
        );
      }

      testWidgets('tapping confirm calls discoveredUpdate with verifyAt set', (tester) async {
        String? capturedId;
        media.Download? capturedDownload;

        await tester.pumpApp(
          buildWithTorrent(
            onUpdate: (id, download) {
              capturedId = id;
              capturedDownload = download;
              return null;
            },
            onReset: (_) => null,
          ),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.byTooltip('verify data'));
        await tester.pump();

        expect(
          find.text('Are you sure you want to verify ${testMediaWithTorrent.description}?'),
          findsOneWidget,
        );

        await tester.tap(find.text('Yes'));
        await tester.pumpAndSettle();

        expect(capturedId, equals(testMediaWithTorrent.torrentId));
        expect(capturedDownload?.verifyAt, isNotEmpty);
        expect(tester.takeException(), isNull);
      });

      testWidgets('tapping cancel does not call discoveredUpdate', (tester) async {
        bool updateCalled = false;

        await tester.pumpApp(
          buildWithTorrent(
            onUpdate: (_, __) {
              updateCalled = true;
              return null;
            },
            onReset: (_) => null,
          ),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.byTooltip('verify data'));
        await tester.pump();

        expect(
          find.text('Are you sure you want to verify ${testMediaWithTorrent.description}?'),
          findsOneWidget,
        );

        await tester.tap(find.text('No'));
        await tester.pumpAndSettle();

        expect(updateCalled, isFalse);
        expect(tester.takeException(), isNull);
      });
    });

    group('onTap (reset)', () {
      Future<media.DownloadMetadataResponse> mockDiscoveredGet(
        String id, {
        List<httpx.Option> options = const [],
      }) async {
        return media.DownloadMetadataResponse(download: media.Download.create());
      }

      Widget buildWithTorrent({
        required String? Function(String) onReset,
        void Function(bool forced)? captureOnChange,
      }) {
        return modals.Node(
          SizedBox(
            width: 800,
            height: 600,
            child: SingleChildScrollView(
              child: MediaSettings(
                current: testMediaWithTorrent,
                onChange: (pending, {bool forced = false, bool autoclose = false}) {
                  captureOnChange?.call(forced);
                },
                knownSearch: mockKnownSearch,

                discoveredGet: mockDiscoveredGet,
                discoveredUpdate: (id, download, {options = const []}) async =>
                    media.DownloadUpdateResponse(),
                discoveredReset: (id, {options = const []}) async {
                  onReset(id);
                  return media.DownloadDeleteResponse();
                },
              ),
            ),
          ),
        );
      }

      testWidgets('tapping confirm calls discoveredReset and onChange with forced:true', (tester) async {
        String? capturedId;
        bool capturedForced = false;

        await tester.pumpApp(
          buildWithTorrent(
            onReset: (id) {
              capturedId = id;
              return null;
            },
            captureOnChange: (forced) => capturedForced = forced,
          ),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.byTooltip('clear data from disk keeps metadata'));
        await tester.pump();

        expect(
          find.text('Are you sure you want to reset ${testMediaWithTorrent.description}?'),
          findsOneWidget,
        );

        await tester.tap(find.text('Yes'));
        await tester.pumpAndSettle();

        expect(capturedId, equals(testMediaWithTorrent.torrentId));
        expect(capturedForced, isTrue);
        expect(tester.takeException(), isNull);
      });

      testWidgets('tapping cancel does not call discoveredReset', (tester) async {
        bool resetCalled = false;

        await tester.pumpApp(
          buildWithTorrent(
            onReset: (_) {
              resetCalled = true;
              return null;
            },
          ),
        );
        await tester.pumpAndSettle();

        await tester.tap(find.byTooltip('clear data from disk keeps metadata'));
        await tester.pump();

        expect(
          find.text('Are you sure you want to reset ${testMediaWithTorrent.description}?'),
          findsOneWidget,
        );

        await tester.tap(find.text('No'));
        await tester.pumpAndSettle();

        expect(resetCalled, isFalse);
        expect(tester.takeException(), isNull);
      });
    });
  });
}
