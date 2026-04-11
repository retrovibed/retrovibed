import 'package:flutter/material.dart';

class Glow extends StatelessWidget {
  final Widget child;
  final List<BoxShadow> tint;
  final BorderRadius borderRadius;

  const Glow(
    this.child, {
    super.key,
    this.tint = const [],
    this.borderRadius = BorderRadius.zero,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        boxShadow: tint,
        borderRadius: borderRadius,
      ),
      child: child,
    );
  }
}
