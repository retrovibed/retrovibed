import 'dart:io';

import 'package:flutter/material.dart';
import 'dart:ui' show lerpDouble;

class Defaults extends ThemeExtension<Defaults> {
  static const defaults = Defaults();

  static final kDesktop = Platform.isLinux || Platform.isWindows || Platform.isMacOS;
  static final kMobile = Platform.isAndroid || Platform.isIOS;
  static final kBorderRadius = BorderRadius.circular(10.0);
  static const kSpacing = 10.0;
  static const kPadding = EdgeInsets.all(16.0);
  static const kMargin = EdgeInsets.all(5.0);
  static const kDanger = Color.fromRGBO(110, 1, 1, 0.75);
  static const kSuccess = Color.fromARGB(255, 0, 255, 17);
  static const kOpaque = Color.fromRGBO(0, 0, 0, 0.80);
  static const kHighlight = Color.fromRGBO(140, 120, 220, 0.25);
  static const kDangerTint = [
    BoxShadow(
      color: Color.fromRGBO(110, 1, 1, 0.4),
      spreadRadius: 1,
      blurRadius: 10,
    ),
  ];
  static const kHighlightTint = [
    BoxShadow(
      color: Color.fromRGBO(140, 120, 220, 0.4),
      spreadRadius: 1,
      blurRadius: 10,
    ),
  ];
  static const kBorder = Border.fromBorderSide(
    BorderSide(color: Color(0xFF000000)),
  );
  static const kCompact = 400.0;
  static const kModalWidth = 0.9;
  static const kModalHeight = 0.85;

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
  final double? _modalWidth;
  final double? _modalHeight;
  final bool? _isCompact;

  bool get desktop => _desktop ?? kDesktop;
  bool get mobile => _mobile ?? kMobile;
  double get spacing => _spacing ?? kSpacing;
  EdgeInsets get padding => _padding ?? kPadding;
  EdgeInsets get margin => _margin ?? kMargin;
  Border get border => _border ?? kBorder;
  BorderRadius get borderRadius => _borderRadius ?? kBorderRadius;
  Color get danger => _danger ?? kDanger;
  Color get success => _success ?? kSuccess;
  Color get opaque => _opaque ?? kOpaque;
  Color get highlight => _highlight ?? kHighlight;
  List<BoxShadow> get dangerTint => _dangerTint ?? kDangerTint;
  List<BoxShadow> get highlightTint => _highlightTint ?? kHighlightTint;
  double get compact => _compact ?? kCompact;
  double get modalWidth => _modalWidth ?? kModalWidth;
  double get modalHeight => _modalHeight ?? kModalHeight;
  bool get isCompact => _isCompact ?? false;

  // the box a modal may occupy. content that has to fit, rather than scroll, needs a
  // bound: the modal node centers and scrolls whatever it is handed, so a child left
  // unbounded renders at its natural size.
  BoxConstraints modal(BuildContext context) {
    final viewport = MediaQuery.sizeOf(context);
    return BoxConstraints(
      maxWidth: viewport.width * modalWidth,
      maxHeight: viewport.height * modalHeight,
    );
  }

  // The constructor accepts nullable values to create partial instances.
  const Defaults({
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
    double? modalWidth,
    double? modalHeight,
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
       _desktop = desktop,
       _mobile = mobile,
       _compact = compact,
       _modalWidth = modalWidth,
       _modalHeight = modalHeight,
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
    bool? desktop,
    bool? mobile,
    double? compact,
    double? modalWidth,
    double? modalHeight,
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
      desktop: desktop ?? _desktop,
      mobile: mobile ?? _mobile,
      compact: compact ?? _compact,
      modalWidth: modalWidth ?? _modalWidth,
      modalHeight: modalHeight ?? _modalHeight,
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
        _dangerTint ?? kDangerTint,
        other._dangerTint ?? kDangerTint,
        t,
      ),
      highlightTint: BoxShadow.lerpList(
        _highlightTint ?? kHighlightTint,
        other._highlightTint ?? kHighlightTint,
        t,
      ),
      desktop: other._desktop ?? _desktop,
      mobile: other._mobile ?? _mobile,
      compact: lerpDouble(_compact, other._compact, t),
      modalWidth: lerpDouble(_modalWidth, other._modalWidth, t),
      modalHeight: lerpDouble(_modalHeight, other._modalHeight, t),
      isCompact: other._isCompact ?? _isCompact,
    );
  }
}
