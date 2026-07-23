import 'dart:async';
import 'dart:io';
import 'package:fixnum/fixnum.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/media/download.watch.dart';
import 'package:retrovibed/media/api.dart' as api;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

const _reconnectInterval = Duration(milliseconds: 50);

api.Download _download() => api.Download(
  media: api.Media(
    id: 'test-id-1',
    description: 'Test Media Description',
    mimetype: 'video/mp4',
    createdAt: '2025-01-01T00:00:00Z',
    archiveId: uuidx.min(),
    torrentId: uuidx.min(),
    knownMediaId: uuidx.min(),
  ),
  bytes: Int64(1024 * 1024 * 500),
  downloaded: Int64(1024 * 1024 * 250),
);

void main() {
  group('ErrorTests.socketclosed', () {
    test('identifies Reading from a closed socket', () {
      expect(ds.ErrorTests.socketclosed(SocketException('Reading from a closed socket')), isTrue);
    });

    test('does not match other SocketExceptions', () {
      expect(ds.ErrorTests.socketclosed(SocketException('Connection refused')), isFalse);
    });

    test('does not match WebSocketException', () {
      expect(ds.ErrorTests.socketclosed(WebSocketException('abnormal')), isFalse);
    });
  });

  group('RefreshingDownload reconnect', () {
    testWidgets('reconnects after socket closed by server', (WidgetTester tester) async {
      var watchCallCount = 0;
      final controllers = <StreamController<api.Download>>[];

      Future<Stream<api.Download>> mockWatch(
        String id, {
        List<httpx.Option> options = const [],
      }) async {
        watchCallCount++;
        final c = StreamController<api.Download>();
        controllers.add(c);
        return c.stream;
      }

      await tester.pumpApp(
        RefreshingDownload(
          current: _download(),
          interval: _reconnectInterval,
          watch: mockWatch,
        ),
      );
      await tester.pump(); // let initState future resolve

      controllers.first.addError(SocketException('Reading from a closed socket'));
      await tester.pump(_reconnectInterval * 2);

      expect(watchCallCount, equals(2));
      expect(find.text('an unexpected problem has occurred'), findsNothing);
    });

    testWidgets('reconnects after WebSocket closes abnormally', (WidgetTester tester) async {
      var watchCallCount = 0;
      final controllers = <StreamController<api.Download>>[];

      Future<Stream<api.Download>> mockWatch(
        String id, {
        List<httpx.Option> options = const [],
      }) async {
        watchCallCount++;
        final c = StreamController<api.Download>();
        controllers.add(c);
        return c.stream;
      }

      await tester.pumpApp(
        RefreshingDownload(
          current: _download(),
          interval: _reconnectInterval,
          watch: mockWatch,
        ),
      );
      await tester.pump();

      controllers.first.addError(WebSocketException('connection closed abnormally'));
      await tester.pump(_reconnectInterval * 2);

      expect(watchCallCount, equals(2));
      expect(find.text('an unexpected problem has occurred'), findsNothing);
    });

    testWidgets('reconnects after stream closes cleanly', (WidgetTester tester) async {
      var watchCallCount = 0;
      final controllers = <StreamController<api.Download>>[];

      Future<Stream<api.Download>> mockWatch(
        String id, {
        List<httpx.Option> options = const [],
      }) async {
        watchCallCount++;
        final c = StreamController<api.Download>();
        controllers.add(c);
        return c.stream;
      }

      await tester.pumpApp(
        RefreshingDownload(
          current: _download(),
          interval: _reconnectInterval,
          watch: mockWatch,
        ),
      );
      await tester.pump();

      await controllers.first.close();
      await tester.pump(_reconnectInterval * 2);

      expect(watchCallCount, equals(2));
      expect(find.text('an unexpected problem has occurred'), findsNothing);
    });

    testWidgets('shows error and does not reconnect on non-closure errors', (
      WidgetTester tester,
    ) async {
      var watchCallCount = 0;

      Future<Stream<api.Download>> mockWatch(
        String id, {
        List<httpx.Option> options = const [],
      }) async {
        watchCallCount++;
        final c = StreamController<api.Download>();
        c.addError(Exception('unexpected server error'));
        return c.stream;
      }

      await tester.pumpApp(
        RefreshingDownload(
          current: _download(),
          interval: _reconnectInterval,
          watch: mockWatch,
        ),
      );
      await tester.pump(_reconnectInterval * 2);
      await tester.pumpAndSettle();

      expect(watchCallCount, equals(1));
      expect(find.text('an unexpected problem has occurred'), findsOneWidget);
    });

    testWidgets('retains last received state after reconnect', (WidgetTester tester) async {
      final updated = api.Download(
        media: api.Media(
          id: 'test-id-1',
          description: 'Test Media Description',
          mimetype: 'video/mp4',
          createdAt: '2025-01-01T00:00:00Z',
          archiveId: uuidx.min(),
          torrentId: uuidx.min(),
          knownMediaId: uuidx.min(),
        ),
        bytes: Int64(1024 * 1024 * 500),
        downloaded: Int64(1024 * 1024 * 500),
      );

      final controllers = <StreamController<api.Download>>[];

      await tester.pumpApp(
        RefreshingDownload(
          current: _download(),
          interval: _reconnectInterval,
          watch: (id, {options = const []}) async {
            final c = StreamController<api.Download>();
            controllers.add(c);
            return c.stream;
          },
        ),
      );
      await tester.pump();

      controllers.first.add(updated);
      await tester.pumpAndSettle();

      controllers.first.addError(SocketException('Reading from a closed socket'));
      await tester.pump(_reconnectInterval * 2);

      expect(find.text('an unexpected problem has occurred'), findsNothing);
      expect(find.text('100.00%'), findsOneWidget);
    });
  });
}
