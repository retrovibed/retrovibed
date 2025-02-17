import 'package:flutter/material.dart';

class Full extends StatelessWidget {
  final Widget? child;
  const Full({this.child, super.key});

  @override
  Widget build(BuildContext context) {
    final media = MediaQuery.of(context);
    return SizedBox(
      width: media.size.width,
      height: media.size.height,
      child: child,
    );
  }
}
