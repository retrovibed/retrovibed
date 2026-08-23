import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart' as remote;

// A mute toggle paired with a draggable volume slider, meant to occupy its
// own row rather than sit inline with the other (icon-button-only) transport
// controls.
class PlayerControlVolume extends StatefulWidget {
  final remote.RemoteControlSocket socket;
  final remote.Volume current;

  const PlayerControlVolume({Key? key, required this.socket, required this.current}) : super(key: key);

  @override
  State<PlayerControlVolume> createState() => _State();
}

class _State extends State<PlayerControlVolume> {
  // overrides widget.current while a drag is in progress, since the slider's
  // displayed position is otherwise driven by the daemon's echoed volume,
  // which would make it lag behind the pointer until the round trip returns.
  double? _dragging;

  @override
  Widget build(BuildContext context) {
    final value = (_dragging ?? widget.current.level).clamp(0.0, 100.0);
    final help = widget.current.muted ? "unmute the remote device" : "mute the remote device";

    return Row(
      children: [
        ds.LoadingIconButton(
          onPressed: ds.LoadingIconButton.convert(
            () => widget.socket.send(remote.messages.volume(widget.current.level, !widget.current.muted)),
          ),
          icon: Icon(widget.current.muted ? Icons.volume_off_rounded : Icons.volume_up_rounded),
          tooltip: help,
          help: ds.Hint(Text(help)),
        ),
        Expanded(
          child: Slider(
            value: value,
            min: 0,
            max: 100,
            label: "${value.round()}%",
            onChanged: (v) => setState(() => _dragging = v),
            onChangeEnd: (v) {
              widget.socket.send(remote.messages.volume(v, widget.current.muted));
              Future.delayed(const Duration(milliseconds: 1000), () {
                if (!mounted) return;
                setState(() => _dragging = null);
              });
            },
          ),
        ),
      ],
    );
  }
}
