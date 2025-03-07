import 'package:flutter/material.dart';
import './borders.dart' as borders;

class Debug extends StatelessWidget {
  final Widget? child;
  final Border? border;
  const Debug(this.child, {super.key, this.border});

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(border: border ?? borders.Debug),
      child: child,
    );
  }
}
