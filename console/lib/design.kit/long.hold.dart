import 'dart:async';
import 'package:flutter/material.dart';

// Fires [onHold] after [duration] of continuous pointer contact.
// Cancels silently if the pointer is lifted or cancelled before [duration].
class LongHold extends StatefulWidget {
  final Widget child;
  final VoidCallback onHold;
  final Duration duration;

  const LongHold({
    super.key,
    required this.child,
    required this.onHold,
    this.duration = const Duration(seconds: 5),
  });

  @override
  State<LongHold> createState() => _LongHoldState();
}

class _LongHoldState extends State<LongHold> {
  Timer? _timer;

  void _start(PointerDownEvent _) {
    _timer = Timer(widget.duration, widget.onHold);
  }

  void _cancel(PointerEvent _) {
    _timer?.cancel();
    _timer = null;
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Listener(
      onPointerDown: _start,
      onPointerUp: _cancel,
      onPointerCancel: _cancel,
      child: widget.child,
    );
  }
}
