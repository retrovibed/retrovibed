import 'package:flutter/material.dart';
import './container.dart' as _ds;
import 'theme.defaults.dart';

class Heading extends StatelessWidget {
  final Widget child;
  final MainAxisSize mainAxisSize;
  final MainAxisAlignment mainAxisAlignment;
  final EdgeInsets margin;
  final EdgeInsets? padding;
  final List<Widget> leading;
  final List<Widget> trailing;
  const Heading(
    Widget this.child, {
    super.key,
    this.mainAxisSize = MainAxisSize.max,
    this.mainAxisAlignment = MainAxisAlignment.center,
    this.margin = EdgeInsets.zero,
    this.padding,
    this.leading = const [],
    this.trailing = const [],
  });

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);
    return _ds.Container(
      padding: padding ?? defaults.padding,
      margin: margin,
      Row(
        mainAxisSize: mainAxisSize,
        mainAxisAlignment: mainAxisAlignment,
        children: [
          ...leading,
          Expanded(child: child),
          ...trailing,
        ],
      ),
    );
  }
}
