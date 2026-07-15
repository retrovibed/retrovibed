import 'package:flutter/material.dart';
import '../empty.dart';

class Guarded extends StatelessWidget {
  final Widget child;
  final Widget overlay;
  final bool enabled;

  const Guarded({
    super.key,
    required this.child,
    this.overlay = Empty,
    this.enabled = false,
  });

  @override
  Widget build(BuildContext context) {
    return enabled ? overlay : child;
  }
}
