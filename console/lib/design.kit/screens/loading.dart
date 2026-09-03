import 'package:flutter/material.dart';
import '../errors.dart' as errors;
import '../empty.dart';
import '../debug.dart';
import 'error.dart';
import 'overlay.dart' as s;

class Loading extends StatelessWidget {
  static const Widget Icon = const Center(
    child: const CircularProgressIndicator(
      padding: EdgeInsets.all(4),
      backgroundColor: Color.fromARGB(0, 0, 0, 0),
      semanticsLabel: 'Linear progress indicator',
    ),
  );

  static Widget Sized({double? width, double? height}) {
    return Container(width: width, height: height, child: Icon);
  }

  final Widget? child;
  final bool loading;
  final bool maintainState;
  final Widget overlay;
  final Widget cause;
  final BorderRadius borderRadius;

  const Loading(
    this.child, {
    super.key,
    this.overlay = Loading.Icon,
    this.loading = false,
    this.maintainState = true,
    this.cause = errors.Error.zero,
    this.borderRadius = BorderRadius.zero,
  });

  @override
  Widget build(BuildContext context) {
    final Widget _overlay = loading ? overlay : Empty;

    return ErrorScreen(
      cause: cause,
      s.Overlay(
        Visibility(
          visible: !loading,
          maintainState: maintainState,
          maintainAnimation: maintainState,
          maintainSize: maintainState,
          child: child ?? Empty,
        ),
        overlay: _overlay,
        borderRadius: borderRadius,
      ),
    );
  }
}
