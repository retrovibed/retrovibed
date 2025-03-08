import 'package:flutter/material.dart';

class Full extends StatelessWidget {
  final Widget? child;
  const Full(this.child, {super.key});

  @override
  Widget build(BuildContext context) {
    final media = MediaQuery.of(context);
    return ConstrainedBox(
      constraints: BoxConstraints(
        minWidth: media.size.width,
        minHeight: media.size.height,
      ),
      child: child,
    );
  }
}
