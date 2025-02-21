import 'package:flutter/material.dart';
import './overlay.dart' as s;

class Loading extends StatelessWidget {
  final Widget? child;
  final bool loading;
  final Widget overlay;
  final Widget? cause;

  const Loading({
    super.key,
    this.child,
    this.overlay = const Center(
      child: CircularProgressIndicator(
        padding: EdgeInsets.all(32.0),
        backgroundColor: Color.fromARGB(0, 0, 0, 0),
        semanticsLabel: 'Linear progress indicator',
      ),
    ),
    this.loading = false,
    this.cause = null,
  });

  @override
  Widget build(BuildContext context) {
    if (loading) {
      return Container(child: overlay);
    }

    return s.Overlay(child: child ?? const SizedBox(), overlay: cause);
    // return child ?? const SizedBox();
  }
}
