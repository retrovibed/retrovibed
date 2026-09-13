import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart' as remote;

// An icon button that, on tap, pops up a small floating card with a mute
// toggle and a vertical volume slider.
class PlayerControlVolume extends StatefulWidget {
  final remote.RemoteControlSocket socket;
  final String sessionId;
  final remote.Sync current;

  const PlayerControlVolume({
    Key? key,
    required this.socket,
    required this.sessionId,
    required this.current,
  }) : super(key: key);

  @override
  State<PlayerControlVolume> createState() => _State();
}

class _State extends State<PlayerControlVolume> with ds.LoadingState {
  OverlayEntry? _entry;

  // overrides widget.current while a drag is in progress, since the popup's
  // displayed position is otherwise driven by the daemon's echoed volume,
  // which would make it lag behind the pointer until the round trip returns.
  double? _dragging;

  @override
  void dispose() {
    _entry?.remove();
    _entry = null;
    super.dispose();
  }

  @override
  void didUpdateWidget(covariant PlayerControlVolume oldWidget) {
    super.didUpdateWidget(oldWidget);
    // deferred: the entry lives in the ancestor Overlay, not this subtree,
    // so marking it dirty synchronously here (mid-build of this widget's
    // own ancestors, e.g. on every daemon sync echo) trips Flutter's
    // "setState called during build" assertion.
    if (_entry != null) ds.postframe(() => _entry?.markNeedsBuild());
  }

  void _close() {
    _entry?.remove();
    _entry = null;
  }

  void _toggle() {
    if (_entry != null) {
      _close();
      return;
    }

    final box = context.findRenderObject() as RenderBox;
    final overlay = Overlay.of(context);
    final overlayBox = overlay.context.findRenderObject() as RenderBox;
    final anchor = box.localToGlobal(Offset(box.size.width / 2, 0), ancestor: overlayBox);

    _entry = OverlayEntry(
      builder: (context) {
        final value = (_dragging ?? widget.current.volume).clamp(0.0, 100.0);
        return Stack(
          children: [
            Positioned.fill(
              child: GestureDetector(
                behavior: HitTestBehavior.translucent,
                onTap: _close,
              ),
            ),
            Positioned(
              left: anchor.dx,
              top: anchor.dy,
              child: FractionalTranslation(
                translation: const Offset(-0.5, -1.0),
                child: _VolumePopup(
                  value: value,
                  muted: widget.current.muted,
                  onChanged: (v) {
                    _dragging = v;
                    _entry?.markNeedsBuild();
                  },
                  onChangeEnd: (v) {
                    widget.socket.send(
                      remote.messages.volume((v - widget.current.volume).round(), sessionId: widget.sessionId),
                    );
                    Future.delayed(const Duration(milliseconds: 1000), () {
                      _dragging = null;
                      _entry?.markNeedsBuild();
                    });
                  },
                ),
              ),
            ),
          ],
        );
      },
    );
    if (_entry != null) overlay.insert(_entry!);
  }

  @override
  Widget build(BuildContext context) {
    final value = (_dragging ?? widget.current.volume).clamp(0.0, 100.0);
    final icon = widget.current.muted || value == 0
        ? Icons.volume_off_rounded
        : value < 50
        ? Icons.volume_down_rounded
        : Icons.volume_up_rounded;

    return ds.LoadingIconButton(
      onPressed: ds.LoadingIconButton.convert(_toggle),
      icon: Icon(icon),
      tooltip: "volume",
      help: ds.Hint(const Text("view and adjust the remote device's volume")),
    );
  }
}

class _VolumePopup extends StatelessWidget {
  final double value;
  final bool muted;
  final ValueChanged<double> onChanged;
  final ValueChanged<double> onChangeEnd;

  const _VolumePopup({
    required this.value,
    required this.muted,
    required this.onChanged,
    required this.onChangeEnd,
  });

  @override
  Widget build(BuildContext context) {
    return Material(
      elevation: 8,
      borderRadius: BorderRadius.circular(8),
      child: Container(
        width: 48,
        padding: const EdgeInsets.symmetric(vertical: 8),
        child: SizedBox(
          height: 120,
          child: RotatedBox(
            quarterTurns: -1,
            child: Slider(
              value: value,
              min: 0,
              max: 100,
              label: "${value.round()}%",
              onChanged: onChanged,
              onChangeEnd: onChangeEnd,
            ),
          ),
        ),
      ),
    );
  }
}
