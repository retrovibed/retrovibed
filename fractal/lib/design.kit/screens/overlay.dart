import 'package:flutter/material.dart';

class Overlay extends StatelessWidget {
  final Widget child;
  final Widget? overlay;
  final AlignmentGeometry alignment;
  final Function()? onTap;

  const Overlay({
    super.key,
    required this.child,
    this.overlay,
    this.alignment = Alignment.center,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      child: Stack(
        alignment: alignment,
        children: [child, overlay ?? const SizedBox()],
      ),
    );
  }
}
