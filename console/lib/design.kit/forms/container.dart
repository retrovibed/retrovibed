import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;

class Container extends ds.Container {
  Container(
    Widget child, {
    super.key,
    super.decoration,
    super.padding,
    super.margin,
    super.alignment,
  }) : super(FocusTraversalGroup(policy: WidgetOrderTraversalPolicy(), child: child));
}
