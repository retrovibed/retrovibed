import 'package:flutter/material.dart';
import 'buttons.dart';
import 'help.scope.dart';

class HelpToggle extends StatelessWidget {
  const HelpToggle({super.key});

  @override
  Widget build(BuildContext context) {
    return buttons.help(onPressed: () => HelpScope.of(context)?.toggle());
  }
}
