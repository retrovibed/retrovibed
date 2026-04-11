import 'package:flutter/material.dart';
import './errors.dart';

mixin LoadingState<T extends StatefulWidget> on State<T> {
  bool loading = true;
  Widget cause = Error.zero;

  @override
  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void resetCause() {
    setState(() {
      cause = Error.zero;
    });
  }
}
