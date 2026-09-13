import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart' as remote;

// A position/duration gauge paired with a draggable seek slider, meant to
// occupy its own row. Modeled on PlayerControlVolume's original inline
// slider layout (see player.control.volume.dart): there is no absolute-seek
// message, only a relative one (remote.messages.seek(offsetMs, ...)), so
// dragging computes the delta between the drop position and the daemon's
// last-known position, same as volume computes a delta against its own
// last-known value.
class PlayerControlPlayback extends StatefulWidget {
  final remote.RemoteControlSocket socket;
  final String sessionId;
  final remote.Sync current;

  const PlayerControlPlayback({
    Key? key,
    required this.socket,
    required this.sessionId,
    required this.current,
  }) : super(key: key);

  @override
  State<PlayerControlPlayback> createState() => _State();
}

class _State extends State<PlayerControlPlayback> with ds.LoadingState {
  // overrides widget.current while a drag is in progress, since the slider's
  // displayed position is otherwise driven by the daemon's echoed position,
  // which would make it lag behind the pointer until the round trip returns.
  double? _dragging;

  @override
  Widget build(BuildContext context) {
    final remainingMs = (widget.current.playback.duration - widget.current.playback.position).toInt();
    final hasDuration = remainingMs > 0;
    final position = (_dragging ?? widget.current.playback.position.toDouble()).clamp(
      0.0,
      hasDuration ? remainingMs.toDouble() : 0.0,
    );

    return Row(
      children: [
        Expanded(
          child: Slider(
            value: position,
            min: 0,
            max: hasDuration ? remainingMs.toDouble() : 1.0,
            onChanged: !hasDuration ? null : (v) => setState(() => _dragging = v),
            onChangeEnd: !hasDuration
                ? null
                : (v) {
                    final delta = (v - widget.current.playback.position.toDouble()).round();
                    widget.socket.send(remote.messages.seek(delta, sessionId: widget.sessionId));
                    Future.delayed(const Duration(milliseconds: 1000), () {
                      setState(() => _dragging = null);
                    });
                  },
          ),
        ),
        ds.Duration.elapsed(Duration(milliseconds: remainingMs)),
      ],
    );
  }
}
