import 'package:flutter/material.dart';

class Overlay extends StatelessWidget {
  final Widget child;
  final Widget? overlay;

  const Overlay({super.key, required this.child, this.overlay});

  @override
  Widget build(BuildContext context) {
    return Stack(children: [child, overlay ?? const SizedBox()]);
  }
}
