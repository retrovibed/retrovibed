import 'dart:async';
import 'package:flutter/material.dart';

// Poll calls [onTick] every [interval] until it resolves true, then stops.
// interval == Duration.zero means disabled: no timer runs. Switching interval
// from zero to non-zero (re)enables polling; switching to zero cancels it.
class Poll extends StatefulWidget {
  final Widget child;
  final Duration interval;
  final Future<bool> Function() onTick;

  const Poll(
    this.child, {
    super.key,
    required this.interval,
    required this.onTick,
  });

  @override
  State<Poll> createState() => _PollState();
}

class _PollState extends State<Poll> {
  Timer _timer = Timer(Duration.zero, () {});

  @override
  void initState() {
    super.initState();
    _timer.cancel();
    _sync();
  }

  @override
  void didUpdateWidget(Poll old) {
    super.didUpdateWidget(old);
    if (old.interval != widget.interval) {
      _timer.cancel();
      _sync();
    }
  }

  void _sync() {
    if (widget.interval == Duration.zero) return;
    _timer = Timer.periodic(widget.interval, (_) => _tick());
  }

  void _tick() {
    widget.onTick().then((done) {
      if (done) _timer.cancel();
    });
  }

  @override
  void dispose() {
    _timer.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) => widget.child;
}
