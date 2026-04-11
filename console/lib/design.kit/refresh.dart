import 'package:flutter/material.dart';

class RefreshBoundary extends StatefulWidget {
  final Widget child;
  final void Function()? onReset;
  RefreshBoundary(this.child, {super.key, this.onReset});

  static _RefreshBoundary? of(BuildContext context) {
    return context.findAncestorStateOfType<_RefreshBoundary>();
  }

  @override
  State<StatefulWidget> createState() => _RefreshBoundary();
}

class _RefreshBoundary extends State<RefreshBoundary> {
  Key _refresh = UniqueKey();

  void reset() {
    setState(widget.onReset ?? () {
      _refresh = UniqueKey();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Container(key: _refresh, child: widget.child);
  }
}
