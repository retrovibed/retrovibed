import 'dart:async';

import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/timex.dart' as timex;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _infString = timex.formatISO8601(timex.inf);

media.Download _download({
  required String id,
  required String description,
  required String completedAt,
}) {
  return media.Download(
    completedAt: completedAt,
    media: media.Media(
      id: id,
      description: description,
      mimetype: 'video/mp4',
      createdAt: '2025-01-01T00:00:00Z',
      archiveId: uuidx.min(),
      torrentId: uuidx.min(),
      knownMediaId: uuidx.min(),
    ),
  );
}

class _FakeWatch {
  final Map<String, StreamController<media.Download>> _controllers = {};

  StreamController<media.Download> operator [](String id) =>
      _controllers.putIfAbsent(id, () => StreamController<media.Download>());

  Future<Stream<media.Download>> call(String id, {List<httpx.Option> options = const []}) async {
    return this[id].stream;
  }
}

void main() {
  group('DownloadQueue', () {
    testWidgets('shows nothing while the pending list is loading', (tester) async {
      final completer = Completer<List<media.Download>>();

      await tester.pumpApp(media.DownloadQueue(completer.future));
      await tester.pump();

      expect(find.text('1 of 1'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('shows the first download and its position once the queue resolves', (tester) async {
      final watch = _FakeWatch();
      final downloads = [
        _download(id: 'a', description: 'First Download', completedAt: _infString),
        _download(id: 'b', description: 'Second Download', completedAt: _infString),
      ];

      await tester.pumpApp(media.DownloadQueue(Future.value(downloads), watch: watch.call));
      await tester.pumpAndSettle();

      expect(find.text('First Download'), findsOneWidget);
      expect(find.text('1 of 2'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('advances immediately once a watched in-progress download completes', (tester) async {
      final watch = _FakeWatch();
      final downloads = [
        _download(id: 'a', description: 'First Download', completedAt: _infString),
        _download(id: 'b', description: 'Second Download', completedAt: _infString),
      ];

      await tester.pumpApp(media.DownloadQueue(Future.value(downloads), watch: watch.call));
      await tester.pumpAndSettle();
      expect(find.text('First Download'), findsOneWidget);

      watch['a'].add(_download(id: 'a', description: 'First Download', completedAt: '2025-06-01T12:00:00Z'));
      // Three separate frame cycles are needed here: one to deliver the stream
      // event and schedule the zero-delay advance timer, one to fire that
      // timer and advance the index, and one for the newly built item's own
      // connect() chain to resolve.
      await tester.pump(const Duration(seconds: 1));
      await tester.pump(const Duration(seconds: 1));
      await tester.pump(const Duration(seconds: 1));

      expect(find.text('Second Download'), findsOneWidget);
      expect(find.text('2 of 2'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('holds an already-completed download for the configured minimum before advancing', (tester) async {
      final watch = _FakeWatch();
      final downloads = [
        _download(id: 'a', description: 'First Download', completedAt: '2025-06-01T12:00:00Z'),
        _download(id: 'b', description: 'Second Download', completedAt: _infString),
      ];

      await tester.pumpApp(
        media.DownloadQueue(
          Future.value(downloads),
          watch: watch.call,
          minCompletedDisplay: const Duration(milliseconds: 200),
        ),
      );
      // flush the queue future, the item's own connect future, and its
      // postFrameCallback-scheduled completion check without advancing the
      // fake clock, so the 200ms hold below starts from zero.
      await tester.pump();
      await tester.pump();
      await tester.pump();

      expect(find.text('First Download'), findsOneWidget);

      await tester.pump(const Duration(milliseconds: 100));
      expect(find.text('First Download'), findsOneWidget);

      await tester.pump(const Duration(milliseconds: 150));
      await tester.pump();

      expect(find.text('Second Download'), findsOneWidget);
      expect(find.text('2 of 2'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('calls onQueueComplete and renders nothing once the last item finishes', (tester) async {
      final watch = _FakeWatch();
      final downloads = [_download(id: 'a', description: 'Only Download', completedAt: _infString)];
      var completed = 0;

      await tester.pumpApp(
        media.DownloadQueue(
          Future.value(downloads),
          watch: watch.call,
          onQueueComplete: () => completed++,
        ),
      );
      await tester.pumpAndSettle();
      expect(completed, 0);

      watch['a'].add(_download(id: 'a', description: 'Only Download', completedAt: '2025-06-01T12:00:00Z'));
      await tester.pump(const Duration(seconds: 1));
      await tester.pump(const Duration(seconds: 1));
      await tester.pump(const Duration(seconds: 1));

      expect(completed, 1);
      expect(find.text('Only Download'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('calls onQueueComplete immediately for an empty pending list', (tester) async {
      var completed = 0;

      await tester.pumpApp(
        media.DownloadQueue(
          Future.value(const <media.Download>[]),
          onQueueComplete: () => completed++,
        ),
      );
      await tester.pumpAndSettle();

      expect(completed, 1);
      expect(tester.takeException(), isNull);
    });

    testWidgets('surfaces an error without throwing when the queue future fails', (tester) async {
      final completer = Completer<List<media.Download>>();

      await tester.pumpApp(media.DownloadQueue(completer.future));
      await tester.pump(); // let AuthzCache settle and mount DownloadQueue before the future errors
      completer.completeError('boom');
      await tester.pumpAndSettle();

      expect(tester.takeException(), isNull);
    });
  });
}
