// ignore: deprecated_member_use
import 'package:flutter/material.dart' as sdk show Theme, ThemeData;
import 'package:material_ui/material_ui.dart';

/// Bridges the material_ui [Theme] to the Flutter SDK material theme.
///
/// material_ui defines its own [Theme], separate from `package:flutter/material.dart`.
/// Dependencies still built on the SDK material library (media_kit_video, mobile_scanner, …)
/// find no SDK theme ancestor and fall back to the default light theme, rendering dark
/// icons/text. This installs an SDK theme derived from the ambient material_ui theme.
class MaterialLegacy extends StatelessWidget {
  final Widget child;
  const MaterialLegacy(this.child, {super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return sdk.Theme(
      data: sdk.ThemeData(brightness: theme.brightness, hoverColor: theme.hoverColor),
      child: child,
    );
  }
}
