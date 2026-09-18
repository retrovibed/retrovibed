import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta/uploads.node.dart';
import 'package:retrovibed/meta/uploads.row.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

const _removalDelay = Duration(milliseconds: 50);

httpx.UploadProgress _event(String id, String name, {required int uploaded, required int total}) =>
    (id, name, "video/mp4", uploaded, total);

void main() {
  group('UploadsRow', () {
    testWidgets('renders nothing when there are no uploads in flight', (tester) async {
      await tester.pumpApp(UploadNode(const UploadsRow()));
      await tester.pump();

      expect(find.byType(UploadsRow), findsOneWidget);
      expect(find.byType(CircularProgressIndicator), findsNothing);
    });

    testWidgets('shows a chip for an upload in progress', (tester) async {
      await tester.pumpApp(UploadNode(const UploadsRow()));
      final state = UploadNode.of(tester.element(find.byType(UploadsRow)));

      state.progress.add(_event('a', 'movie.mp4', uploaded: 40, total: 100));
      await tester.pump();
      await tester.pump();

      expect(find.text('movie.mp4'), findsOneWidget);
    });

    testWidgets('keeps showing a completed upload until the removal delay elapses', (tester) async {
      await tester.pumpApp(UploadNode(const UploadsRow(), delay: _removalDelay));
      final state = UploadNode.of(tester.element(find.byType(UploadsRow)));

      state.progress.add(_event('a', 'movie.mp4', uploaded: 100, total: 100));
      await tester.pump();
      await tester.pump();
      expect(find.text('movie.mp4'), findsOneWidget, reason: 'still shown right after completing');

      await tester.pump(_removalDelay - const Duration(milliseconds: 1));
      expect(find.text('movie.mp4'), findsOneWidget, reason: 'still shown just before the removal delay elapses');

      await tester.pump(const Duration(milliseconds: 2));
      expect(find.text('movie.mp4'), findsNothing, reason: 'removed once the removal delay elapses');
    });

    testWidgets('removing a completed upload leaves one still in flight alone', (tester) async {
      await tester.pumpApp(UploadNode(const UploadsRow(), delay: _removalDelay));
      final state = UploadNode.of(tester.element(find.byType(UploadsRow)));

      state.progress.add(_event('a', 'done.mp4', uploaded: 100, total: 100));
      state.progress.add(_event('b', 'inflight.mp4', uploaded: 40, total: 100));
      await tester.pump();
      await tester.pump();

      await tester.pump(_removalDelay + const Duration(milliseconds: 1));

      expect(find.text('done.mp4'), findsNothing);
      expect(find.text('inflight.mp4'), findsOneWidget);
    });

    testWidgets('an upload completing later gets swept away by an earlier completion timer', (tester) async {
      // documents current behavior rather than a desired one: the removal callback clears every
      // completed entry, not just the one that scheduled it, so two completions coalesce onto
      // whichever timer fires first instead of each being removed `delay` after its own completion.
      await tester.pumpApp(UploadNode(const UploadsRow(), delay: _removalDelay));
      final state = UploadNode.of(tester.element(find.byType(UploadsRow)));

      state.progress.add(_event('a', 'first.mp4', uploaded: 100, total: 100));
      await tester.pump();
      await tester.pump();

      await tester.pump(_removalDelay ~/ 2);
      state.progress.add(_event('b', 'second.mp4', uploaded: 100, total: 100));
      await tester.pump();
      await tester.pump();
      expect(find.text('first.mp4'), findsOneWidget);
      expect(find.text('second.mp4'), findsOneWidget);

      // 'a's timer fires here, well before 'b's own `delay` has elapsed since it completed.
      await tester.pump(_removalDelay ~/ 2 + const Duration(milliseconds: 1));
      expect(find.text('first.mp4'), findsNothing);
      expect(find.text('second.mp4'), findsNothing, reason: "swept away early by a's removal timer");

      // drain 'b's own now-redundant removal timer so it doesn't outlive the test.
      await tester.pump(_removalDelay);
    });

    testWidgets('a later progress event for the same id replaces the earlier one', (tester) async {
      await tester.pumpApp(UploadNode(const UploadsRow()));
      final state = UploadNode.of(tester.element(find.byType(UploadsRow)));

      state.progress.add(_event('a', 'movie.mp4', uploaded: 10, total: 100));
      await tester.pump();
      await tester.pump();
      state.progress.add(_event('a', 'movie.mp4', uploaded: 90, total: 100));
      await tester.pump();
      await tester.pump();

      expect(find.text('movie.mp4'), findsOneWidget);
    });
  });
}
