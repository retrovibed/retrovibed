import 'package:flutter/material.dart';
import 'package:flutter/services.dart' as services;
import 'package:retrovibed/designkit.dart' as ds;

class Copyable extends StatelessWidget {
  final Widget content;
  final MainAxisSize mainAxisSize;
  final VoidCallback? onPressed;

  const Copyable(
    this.content, {
    super.key,
    this.onPressed,
    this.mainAxisSize = MainAxisSize.max,
  });

  static VoidCallback copy(String value) {
    return () {
      services.Clipboard.setData(services.ClipboardData(text: value));
    };
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return Row(
      spacing: defaults.spacing,
      mainAxisSize: mainAxisSize,
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Flexible(child: content),
        ds.buttons.copy(onPressed: onPressed),
      ],
    );
  }
}
