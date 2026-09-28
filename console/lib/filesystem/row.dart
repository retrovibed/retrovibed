import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;

class FilesystemRow extends StatelessWidget {
  final media.Media current;
  final Future<void> Function()? onTap;
  final List<Widget> trailing;
  final Widget expanded;
  // the details (expanded) are open.
  final bool focused;

  const FilesystemRow(
    this.current, {
    super.key,
    this.onTap,
    this.trailing = const [],
    this.expanded = ds.Empty,
    this.focused = false,
  });

  @override
  Widget build(BuildContext context) {
    return media.RowDisplay(
      media: current,
      highlighted: focused,
      leading: [Icon(mimex.icon(current.mimetype))],
      trailing: trailing,
      expanded: expanded,
      onTap: onTap,
    );
  }
}
