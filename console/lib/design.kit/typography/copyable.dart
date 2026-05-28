import 'package:flutter/material.dart';
import 'package:flutter/services.dart' as services;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/help.dart';

class Copyable extends StatelessWidget {
  final Widget content;
  final MainAxisSize mainAxisSize;
  final ds.AsyncVoidCallback? onPressed;
  final Widget help;

  const Copyable(
    this.content, {
    super.key,
    this.onPressed,
    this.mainAxisSize = MainAxisSize.max,
    this.help = HelpScope.None,
  });

  static ds.AsyncVoidCallback copy(String value) {
    return () async {
      return services.Clipboard.setData(services.ClipboardData(text: value));
    };
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return Help(
      Row(
        spacing: defaults.spacing,
        mainAxisSize: mainAxisSize,
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Flexible(child: content),
          ds.buttons.copy(onPressed: onPressed),
        ],
      ),
      help,
    );
  }
}
