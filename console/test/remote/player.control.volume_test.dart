import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/remote/api.dart' as remote;
import 'package:retrovibed/remote/player.control.volume.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

class _FakeRemoteControlSocket implements remote.RemoteControlSocket {
  final StreamController<remote.Stream> _incoming = StreamController();
  final List<remote.Stream> sent = [];

  @override
  Stream<remote.Stream> get messages => _incoming.stream;

  @override
  void send(remote.Stream msg) => sent.add(msg);

  @override
  Future<void> close() async {
    await _incoming.close();
  }
}

void main() {
  group('PlayerControlVolume', () {
    testWidgets('shows the muted icon when current.muted is true', (tester) async {
      await tester.pumpApp(
        PlayerControlVolume(
          socket: _FakeRemoteControlSocket(),
          sessionId: 's1',
          current: remote.Sync(volume: 80, muted: true),
        ),
      );

      expect(find.byIcon(Icons.volume_off_rounded), findsOneWidget);
    });

    testWidgets('shows a low-volume icon when unmuted and below 50', (tester) async {
      await tester.pumpApp(
        PlayerControlVolume(
          socket: _FakeRemoteControlSocket(),
          sessionId: 's1',
          current: remote.Sync(volume: 20, muted: false),
        ),
      );

      expect(find.byIcon(Icons.volume_down_rounded), findsOneWidget);
    });

    testWidgets('shows a high-volume icon when unmuted and at or above 50', (tester) async {
      await tester.pumpApp(
        PlayerControlVolume(
          socket: _FakeRemoteControlSocket(),
          sessionId: 's1',
          current: remote.Sync(volume: 80, muted: false),
        ),
      );

      expect(find.byIcon(Icons.volume_up_rounded), findsOneWidget);
    });

    testWidgets('tapping the icon opens a popup slider seeded at the current volume', (tester) async {
      await tester.pumpApp(
        PlayerControlVolume(
          socket: _FakeRemoteControlSocket(),
          sessionId: 's1',
          current: remote.Sync(volume: 42, muted: false),
        ),
      );

      expect(find.byType(Slider), findsNothing);

      await tester.tap(find.byTooltip('volume'));
      await tester.pump();

      final slider = tester.widget<Slider>(find.byType(Slider));
      expect(slider.value, 42);
    });

    testWidgets('tapping the icon again closes the popup', (tester) async {
      await tester.pumpApp(
        PlayerControlVolume(
          socket: _FakeRemoteControlSocket(),
          sessionId: 's1',
          current: remote.Sync(volume: 42, muted: false),
        ),
      );

      await tester.tap(find.byTooltip('volume'));
      await tester.pump();
      expect(find.byType(Slider), findsOneWidget);

      await tester.tap(find.byTooltip('volume'));
      await tester.pump();
      expect(find.byType(Slider), findsNothing);
    });

    testWidgets('tapping outside the popup dismisses it', (tester) async {
      await tester.pumpApp(
        PlayerControlVolume(
          socket: _FakeRemoteControlSocket(),
          sessionId: 's1',
          current: remote.Sync(volume: 42, muted: false),
        ),
      );

      await tester.tap(find.byTooltip('volume'));
      await tester.pump();
      expect(find.byType(Slider), findsOneWidget);

      // corner of the screen, well away from the popup and its anchor.
      await tester.tapAt(const Offset(5, 5));
      await tester.pump();
      expect(find.byType(Slider), findsNothing);
    });

    testWidgets('releasing the popup slider at max sends a volume message to max', (tester) async {
      final socket = _FakeRemoteControlSocket();
      await tester.pumpApp(
        PlayerControlVolume(
          socket: socket,
          sessionId: 's1',
          current: remote.Sync(volume: 30, muted: false),
        ),
      );

      await tester.tap(find.byTooltip('volume'));
      await tester.pump();

      // exercised via the callback directly rather than a simulated drag:
      // the popup slider is rotated (quarterTurns: -1) and re-anchored from
      // the button's RenderBox on open, both of which make pixel-accurate
      // drag gestures brittle in a test harness; onChangeEnd is a plain
      // ValueChanged<double> field, so invoking it exercises the same
      // widget logic without depending on gesture/layout geometry.
      tester.widget<Slider>(find.byType(Slider)).onChangeEnd!(100);
      await tester.pump(const Duration(seconds: 1)); // settle the drag-reset timer

      expect(socket.sent, hasLength(1));
      final sent = socket.sent.single;
      expect(sent.hasVolume(), isTrue);
      expect(sent.sessionId, 's1');
      expect(sent.volume.offset, 70); // 100 (max) - 30 (initial)
    });

    testWidgets('releasing the popup slider at min sends a volume message to min', (tester) async {
      final socket = _FakeRemoteControlSocket();
      await tester.pumpApp(
        PlayerControlVolume(
          socket: socket,
          sessionId: 's1',
          current: remote.Sync(volume: 30, muted: false),
        ),
      );

      await tester.tap(find.byTooltip('volume'));
      await tester.pump();

      tester.widget<Slider>(find.byType(Slider)).onChangeEnd!(0);
      await tester.pump(const Duration(seconds: 1)); // settle the drag-reset timer

      expect(socket.sent, hasLength(1));
      final sent = socket.sent.single;
      expect(sent.hasVolume(), isTrue);
      expect(sent.sessionId, 's1');
      expect(sent.volume.offset, -30); // 0 (min) - 30 (initial)
    });
  });
}
