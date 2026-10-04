import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './media.pb.dart';

class RowDisplay extends StatelessWidget {
  final Media media;
  final List<Widget> leading;
  final List<Widget> trailing;
  final Future<void> Function()? onTap;
  final Future<void> Function()? onDoubleTap;
  final Widget help;
  final Widget expanded;
  final bool highlighted;
  final Border? border;
  final EdgeInsets? padding;
  const RowDisplay({
    super.key,
    required this.media,
    this.leading = const [],
    this.trailing = const [],
    this.onTap,
    this.onDoubleTap,
    this.help = ds.HelpScope.None,
    this.expanded = ds.Empty,
    this.highlighted = false,
    this.border,
    this.padding,
  });

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return ds.Help(
      ds.TableRow(
        padding: padding ?? defaults.padding,
        onTap: onTap,
        tint: highlighted ? defaults.highlightTint : [],
        expanded: expanded,
        autoexpand: true,
        border: border,
        [
          ...leading,
          Expanded(child: Text(media.description, overflow: TextOverflow.ellipsis)),
          ...trailing,
        ],
      ),
      help,
    );
  }
}
