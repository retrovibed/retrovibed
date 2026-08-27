import 'package:flutter/material.dart';
import 'buttons.dart';
import 'help.scope.dart';

class HelpToggle extends StatelessWidget {
  const HelpToggle({super.key});

  @override
  Widget build(BuildContext context) {
    if (!HelpScope.visible(context)) return const SizedBox.shrink();
    return buttons.help(onPressed: () => HelpScope.of(context)?.toggle());
  }
}
