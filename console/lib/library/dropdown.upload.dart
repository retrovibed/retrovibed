import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;

class DropdownUpload extends StatelessWidget {
  final Widget icon;
  final List<PopupMenuEntry<String>> items;
  final Widget help;

  const DropdownUpload({
    super.key,
    required this.icon,
    required this.items,
    this.help = ds.HelpScope.None,
  });

  @override
  Widget build(BuildContext context) {
    return ds.Help(
      PopupMenuButton<String>(
        position: PopupMenuPosition.under,
        color: Theme.of(context).colorScheme.surface,
        surfaceTintColor: Theme.of(context).colorScheme.surface,
        icon: icon,
        itemBuilder: (context) => items,
      ),
      help,
    );
  }
}
