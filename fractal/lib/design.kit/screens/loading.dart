import 'package:flutter/material.dart';

class Loading extends StatelessWidget {
  final Widget? child;
  final bool loading;
  final Widget overlay;

  const Loading({
    super.key,
    this.child,
    this.overlay = const CircularProgressIndicator(
      backgroundColor: Color.fromARGB(0, 0, 0, 0),
      semanticsLabel: 'Linear progress indicator',
    ),
    this.loading = false,
  });

  @override
  Widget build(BuildContext context) {
    if (loading) {
      return FractionallySizedBox(
        child: Center(
          child: overlay,
        ),
      );
    }

    return child ?? const SizedBox();
  }
}
