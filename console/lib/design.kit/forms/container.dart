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
    Widget cause = ds.Error.zero,
    bool loading = false,
  }) : super(
         ds.Loading(
           FocusTraversalGroup(policy: WidgetOrderTraversalPolicy(), child: child),
           cause: cause,
           loading: loading,
         ),
       );
}
