import 'package:flutter/material.dart';
import '../errors.dart' as errors;
import 'overlay.dart' as s;
import 'glow.dart' as g;

class ErrorScreen extends StatelessWidget {
  final Widget child;
  final Widget cause;
  final BorderRadius borderRadius;
  final List<BoxShadow> tint;

  const ErrorScreen(
    this.child, {
    super.key,
    this.cause = errors.Error.zero,
    this.borderRadius = BorderRadius.zero,
    this.tint = const [],
  });

  bool get hasError => cause != errors.Error.zero;

  @override
  Widget build(BuildContext context) {
    final Widget _overlay = hasError ? cause : const SizedBox();
    final List<BoxShadow> _tint = hasError ? tint : [];

    return LayoutBuilder(
      builder: (context, constraints) {
        final isBounded = constraints.hasBoundedWidth && constraints.hasBoundedHeight;
        final expand = isBounded && hasError;
        return s.Overlay(
          SizedBox(
            width: expand ? double.infinity : null,
            height: expand ? double.infinity : null,
            child: g.Glow(child, tint: _tint, borderRadius: borderRadius),
          ),
          overlay: Positioned.fill(child: _overlay),
          borderRadius: borderRadius,
        );
      },
    );
  }
}
