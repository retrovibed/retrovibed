import 'package:flutter/material.dart';

class Gauge extends StatelessWidget {
  /// The percentage of the gauge that should be filled (0.0 to 1.0).
  final double fill;

  /// The height/thickness of the gauge bar.
  final double height;

  /// The foreground color (filled part) of the gauge.
  /// Defaults to [ColorScheme.primary].
  final Color? color;

  /// The background color (unfilled part) of the gauge.
  /// Defaults to [ColorScheme.surfaceContainerHighest].
  final Color? background;

  const Gauge(
    this.fill, {
    Key? key,
    this.height = 24.0,
    this.color,
    this.background,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final scheme = Theme.of(context).colorScheme;
    // clamping to prevent visual errors
    final double filled = fill.isNaN ? 0.0 : fill.clamp(0.0, 1.0);

    return Container(
      padding: EdgeInsets.zero,
      margin: EdgeInsets.zero,
      height: height,
      clipBehavior: Clip.antiAlias,
      decoration: BoxDecoration(
        color: background ?? scheme.surfaceContainerHighest,
        borderRadius: BorderRadius.circular(height / 2),
      ),
      child: Align(
        alignment: Alignment.centerLeft,
        child: FractionallySizedBox(
          widthFactor: filled,
          heightFactor: 1.0,
          child: ColoredBox(color: color ?? scheme.primary),
        ),
      ),
    );
  }
}
