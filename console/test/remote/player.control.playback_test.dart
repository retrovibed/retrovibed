import 'dart:async';

import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/remote/api.dart' as remote;
import 'package:retrovibed/remote/player.control.playback.dart';
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
  group('PlayerControlPlayback', () {
    testWidgets('slider is disabled when there is no remaining duration', (tester) async {
      await tester.pumpApp(
        PlayerControlPlayback(
          socket: _FakeRemoteControlSocket(),
          sessionId: 's1',
          current: remote.Sync(
            playback: remote.Playback(
              position: fixnum.Int64(0),
              duration: fixnum.Int64(0),
            ),
          ),
        ),
      );

      final slider = tester.widget<Slider>(find.byType(Slider));
      expect(slider.onChanged, isNull);
      expect(slider.onChangeEnd, isNull);
    });

    testWidgets('slider value/max reflect the remaining playback time', (tester) async {
      await tester.pumpApp(
        PlayerControlPlayback(
          socket: _FakeRemoteControlSocket(),
          sessionId: 's1',
          current: remote.Sync(
            playback: remote.Playback(
              position: fixnum.Int64(15000),
              duration: fixnum.Int64(180000),
            ),
          ),
        ),
      );

      final slider = tester.widget<Slider>(find.byType(Slider));
      // remainingMs = duration - position = 165000
      expect(slider.max, 165000);
      expect(slider.value, 15000);
    });

    testWidgets('releasing the slider at the end sends a seek message', (tester) async {
      final socket = _FakeRemoteControlSocket();
      await tester.pumpApp(
        PlayerControlPlayback(
          socket: socket,
          sessionId: 's1',
          current: remote.Sync(
            playback: remote.Playback(
              position: fixnum.Int64(15000),
              duration: fixnum.Int64(180000),
            ),
          ),
        ),
      );

      // exercised via the callback directly rather than a simulated drag,
      // since onChangeEnd is a plain ValueChanged<double> field - this
      // tests the delta computation without depending on drag/layout
      // geometry.
      tester.widget<Slider>(find.byType(Slider)).onChangeEnd!(165000); // slider's max
      await tester.pump(const Duration(seconds: 1)); // settle the drag-reset timer

      expect(socket.sent, hasLength(1));
      final sent = socket.sent.single;
      expect(sent.hasSeek(), isTrue);
      expect(sent.sessionId, 's1');
      expect(sent.seek.offset, 150000); // remainingMs (165000, the slider's max) - position (15000)
    });

    testWidgets('releasing the slider at the start sends a negative seek message', (tester) async {
      final socket = _FakeRemoteControlSocket();
      await tester.pumpApp(
        PlayerControlPlayback(
          socket: socket,
          sessionId: 's1',
          current: remote.Sync(
            playback: remote.Playback(
              position: fixnum.Int64(15000),
              duration: fixnum.Int64(180000),
            ),
          ),
        ),
      );

      tester.widget<Slider>(find.byType(Slider)).onChangeEnd!(0);
      await tester.pump(const Duration(seconds: 1)); // settle the drag-reset timer

      expect(socket.sent, hasLength(1));
      final sent = socket.sent.single;
      expect(sent.hasSeek(), isTrue);
      expect(sent.sessionId, 's1');
      expect(sent.seek.offset, -15000); // 0 (min) - position (15000)
    });
  });
}
