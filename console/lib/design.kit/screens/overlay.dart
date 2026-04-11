import 'package:flutter/material.dart';

class Overlay extends StatelessWidget {
  final Widget child;
  final Widget overlay;
  final AlignmentGeometry alignment;
  final BorderRadius borderRadius;

  const Overlay(
    this.child, {
    super.key,
    this.overlay = const SizedBox(),
    this.alignment = Alignment.center,
    this.borderRadius = BorderRadius.zero,
  });

  factory Overlay.tappable(
    Widget child, {
    Key? key,
    Widget overlay = const SizedBox(),
    AlignmentGeometry alignment = Alignment.center,
    BorderRadius borderRadius = BorderRadius.zero,
    required Function()? onTap,
  }) {
    return Overlay(
      InkWell(onTap: onTap, child: child),
      key: key,
      alignment: alignment,
      overlay: overlay,
      borderRadius: borderRadius,
    );
  }

  @override
  Widget build(BuildContext context) {
    return ClipRRect(
      borderRadius: borderRadius,
      child: Stack(
        fit: StackFit.passthrough,
        alignment: alignment,
        children: [child, overlay],
      ),
    );
  }
}
