import 'package:flutter/material.dart';
import 'empty.dart';

class LazyVisible extends StatefulWidget {
  final bool visible;
  final Widget child;
  final bool maintainState;

  const LazyVisible(
    this.child, {
    super.key,
    required this.visible,
    this.maintainState = true,
  });

  @override
  State<LazyVisible> createState() => _LazyVisibleState();
}

class _LazyVisibleState extends State<LazyVisible> {
  late bool _everVisible = widget.visible;

  @override
  void didUpdateWidget(covariant LazyVisible oldWidget) {
    super.didUpdateWidget(oldWidget);
    _everVisible = _everVisible || widget.visible;
  }

  @override
  Widget build(BuildContext context) {
    return Visibility(
      visible: widget.visible,
      maintainState: widget.maintainState,
      child: _everVisible ? widget.child : Empty,
    );
  }
}
