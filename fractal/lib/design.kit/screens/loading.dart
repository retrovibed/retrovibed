import 'package:flutter/material.dart';

class Loading extends StatelessWidget {
  final Widget? child;
  final bool loading;
  final Widget overlay;

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
  });

  @override
  Widget build(BuildContext context) {
    if (loading) {
      return Container(child: FractionallySizedBox(child: overlay));
    }

    return child ?? const SizedBox();
  }
}
