import 'package:flutter/material.dart';

class Defaults extends ThemeExtension<Defaults> {
  static const defaults = Defaults();

  final double? spacing;
  final EdgeInsets? padding;
  final EdgeInsetsGeometry? margin;
  final Border? border;
  final Color? danger;
  final Color? opaque;

  const Defaults({
    this.spacing = 10.0,
    this.padding = const EdgeInsets.all(16.0),
    this.margin = const EdgeInsets.all(5.0),
    this.danger = const Color.fromRGBO(110, 1, 1, 0.75),
    this.opaque = const Color.fromRGBO(0, 0, 0, 0.80),
    this.border = const Border.fromBorderSide(
      const BorderSide(color: Color(0xFF000000)),
    ),
  });

  static Defaults of(BuildContext context) {
    return Theme.of(context).extension<Defaults>() ?? Defaults.defaults;
  }

  @override
  Defaults copyWith({
    double? spacing,
    EdgeInsets? padding,
    EdgeInsetsGeometry? margin,
    Border? border,
  }) {
    return Defaults(
      spacing: spacing ?? this.spacing,
      padding: padding ?? this.padding,
      margin: margin ?? this.margin,
      border: border ?? this.border,
      danger: danger ?? this.danger,
      opaque: opaque ?? this.opaque,
    );
  }

  @override
  ThemeExtension<Defaults> lerp(
    covariant ThemeExtension<Defaults>? other,
    double t,
  ) {
    if (other is! Defaults) {
      return this;
    }

    return Defaults(
      padding: padding ?? other.padding,
      margin: margin ?? other.margin,
      border: border ?? other.border,
      danger: danger ?? other.danger,
      opaque: opaque ?? other.opaque,
      spacing: spacing ?? other.spacing,
    );
  }
}
