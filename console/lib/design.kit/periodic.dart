import 'dart:async';
import 'package:flutter/material.dart';

class PeriodicBoundary extends StatefulWidget {
  final Widget child;

  PeriodicBoundary(this.child, {super.key});

  static _PeriodicBoundary? of(BuildContext context) {
    return context.findAncestorStateOfType<_PeriodicBoundary>();
  }

  @override
  State<StatefulWidget> createState() => _PeriodicBoundary();
}

class _PeriodicBoundary extends State<PeriodicBoundary> {
  Timer? period;

  Key _refresh = UniqueKey();

  void reset() {
    setState(() {
      _refresh = UniqueKey();
    });
  }

  @override
  void initState() {
    super.initState();
    period = Timer.periodic(const Duration(seconds: 10), (p) => this.reset());
  }

  @override
  void dispose() {
    super.dispose();
    period?.cancel();
  }

  @override
  Widget build(BuildContext context) {
    return Container(key: _refresh, child: widget.child);
  }
}
