import 'dart:io';

import 'package:flutter/foundation.dart' as foundation;
import 'package:flutter/material.dart';
import 'dart:ui' show lerpDouble;

class Defaults extends ThemeExtension<Defaults> {
  static const defaults = Defaults();

  static final _defaultDebug = foundation.kDebugMode;
  static final _defaultDesktop = Platform.isLinux || Platform.isWindows || Platform.isMacOS;
  static final _defaultMobile = Platform.isAndroid || Platform.isIOS;
  static final _defaultBorderRadius = BorderRadius.circular(10.0);
  static const _defaultSpacing = 10.0;
  static const _defaultPadding = EdgeInsets.all(16.0);
  static const _defaultMargin = EdgeInsets.all(5.0);
  static const _defaultDanger = Color.fromRGBO(110, 1, 1, 0.75);
  static const _defaultSuccess = Color.fromARGB(255, 0, 255, 17);
  static const _defaultOpaque = Color.fromRGBO(0, 0, 0, 0.80);
  static const _defaultHighlight = Color.fromRGBO(140, 120, 220, 0.25);
  static const _defaultDangerTint = [
    BoxShadow(
      color: Color.fromRGBO(110, 1, 1, 0.4),
      spreadRadius: 1,
      blurRadius: 10,
    ),
  ];
  static const _defaultHighlightTint = [
    BoxShadow(
      color: Color.fromRGBO(140, 120, 220, 0.4),
      spreadRadius: 1,
      blurRadius: 10,
    ),
  ];
  static const _defaultBorder = Border.fromBorderSide(
    BorderSide(color: Color(0xFF000000)),
  );
  static const _defaultCompact = 400.0;

  final bool? _debug;
  final bool? _desktop;
  final bool? _mobile;
  final double? _spacing;
  final EdgeInsets? _padding;
  final EdgeInsets? _margin;
  final Border? _border;
  final BorderRadius? _borderRadius;
  final Color? _danger;
  final Color? _success;
  final Color? _opaque;
  final Color? _highlight;
  final List<BoxShadow>? _dangerTint;
  final List<BoxShadow>? _highlightTint;
  final double? _compact;
  final bool? _isCompact;

  bool get debug => _debug ?? _defaultDebug;
  bool get desktop => _desktop ?? _defaultDesktop;
  bool get mobile => _mobile ?? _defaultMobile;
  double get spacing => _spacing ?? _defaultSpacing;
  EdgeInsets get padding => _padding ?? _defaultPadding;
  EdgeInsets get margin => _margin ?? _defaultMargin;
  Border get border => _border ?? _defaultBorder;
  BorderRadius get borderRadius => _borderRadius ?? _defaultBorderRadius;
  Color get danger => _danger ?? _defaultDanger;
  Color get success => _success ?? _defaultSuccess;
  Color get opaque => _opaque ?? _defaultOpaque;
  Color get highlight => _highlight ?? _defaultHighlight;
  List<BoxShadow> get dangerTint => _dangerTint ?? _defaultDangerTint;
  List<BoxShadow> get highlightTint => _highlightTint ?? _defaultHighlightTint;
  double get compact => _compact ?? _defaultCompact;
  bool get isCompact => _isCompact ?? false;

  // The constructor accepts nullable values to create partial instances.
  const Defaults({
    bool? debug,
    bool? desktop,
    bool? mobile,
    double? spacing,
    EdgeInsets? padding,
    EdgeInsets? margin,
    Border? border,
    BorderRadius? borderRadius,
    Color? danger,
    Color? success,
    Color? opaque,
    Color? highlight,
    List<BoxShadow>? dangerTint,
    List<BoxShadow>? highlightTint,
    double? compact,
    bool? isCompact,
  }) : _spacing = spacing,
       _padding = padding,
       _margin = margin,
       _border = border,
       _borderRadius = borderRadius,
       _danger = danger,
       _success = success,
       _opaque = opaque,
       _highlight = highlight,
       _dangerTint = dangerTint,
       _highlightTint = highlightTint,
       _debug = debug,
       _desktop = desktop,
       _mobile = mobile,
       _compact = compact,
       _isCompact = isCompact;

  static Defaults of(BuildContext context) {
    return Theme.of(context).extension<Defaults>() ?? Defaults.defaults;
  }

  @override
  Defaults copyWith({
    double? spacing,
    EdgeInsets? padding,
    EdgeInsets? margin,
    Border? border,
    BorderRadius? borderRadius,
    Color? danger,
    Color? success,
    Color? opaque,
    Color? highlight,
    List<BoxShadow>? dangerTint,
    List<BoxShadow>? highlightTint,
    bool? debug,
    bool? desktop,
    bool? mobile,
    double? compact,
    bool? isCompact,
  }) {
    return Defaults(
      spacing: spacing ?? _spacing,
      padding: padding ?? _padding,
      margin: margin ?? _margin,
      border: border ?? _border,
      borderRadius: borderRadius ?? _borderRadius,
      danger: danger ?? _danger,
      success: success ?? _success,
      opaque: opaque ?? _opaque,
      highlight: highlight ?? _highlight,
      dangerTint: dangerTint ?? _dangerTint,
      highlightTint: highlightTint ?? _highlightTint,
      debug: debug ?? _debug,
      desktop: desktop ?? _desktop,
      mobile: mobile ?? _mobile,
      compact: compact ?? _compact,
      isCompact: isCompact ?? _isCompact,
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
      spacing: lerpDouble(_spacing, other._spacing, t),
      padding: EdgeInsets.lerp(_padding, other._padding, t),
      margin: EdgeInsets.lerp(_margin, other._margin, t),
      border: Border.lerp(_border, other._border, t),
      borderRadius: BorderRadius.lerp(_borderRadius, other._borderRadius, t),
      danger: Color.lerp(_danger, other._danger, t),
      success: Color.lerp(_success, other._success, t),
      opaque: Color.lerp(_opaque, other._opaque, t),
      highlight: Color.lerp(_highlight, other._highlight, t),
      dangerTint: BoxShadow.lerpList(
        _dangerTint ?? _defaultDangerTint,
        other._dangerTint ?? _defaultDangerTint,
        t,
      ),
      highlightTint: BoxShadow.lerpList(
        _highlightTint ?? _defaultHighlightTint,
        other._highlightTint ?? _defaultHighlightTint,
        t,
      ),
      debug: other._debug ?? _debug,
      desktop: other._desktop ?? _desktop,
      mobile: other._mobile ?? _mobile,
      compact: lerpDouble(_compact, other._compact, t),
      isCompact: other._isCompact ?? _isCompact,
    );
  }
}
