import 'package:flutter/material.dart';

class GuardedMask extends StatelessWidget {
  final Widget child;
  final bool protected;

  const GuardedMask({
    super.key,
    required this.child,
    this.protected = false,
  });

  @override
  Widget build(BuildContext context) {
    return IgnorePointer(
      ignoring: protected,
      child: Opacity(
        opacity: protected ? 0.5 : 1.0,
        child: child,
      ),
    );
  }
}
