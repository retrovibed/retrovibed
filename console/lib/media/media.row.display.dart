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
  final bool highlighted;
  const RowDisplay({
    super.key,
    required this.media,
    this.leading = const [],
    this.trailing = const [],
    this.onTap,
    this.onDoubleTap,
    this.help = ds.HelpScope.None,
    this.highlighted = false,
  });

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return ds.Help(
      ds.TableRow(
        padding: defaults.padding,
        onTap: onTap,
        tint: highlighted ? defaults.highlightTint : [],
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
