import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './api.dart' as api;

class KnownMediaRowDisplay extends StatelessWidget {
  final api.Known current;
  final List<Widget> leading;
  final List<Widget> trailing;
  final Future<void> Function()? onTap;
  final Future<void> Function()? onDoubleTap;
  final Widget help;
  final bool highlighted;
  const KnownMediaRowDisplay(
    this.current, {
    super.key,
    this.leading = const [],
    this.trailing = const [],
    this.onTap,
    this.onDoubleTap,
    this.help = ds.HelpScope.None,
    this.highlighted = false,
  });

  static Widget future(
    Future<api.Known> future, {
    Key? key,
    List<Widget> leading = const [],
    List<Widget> trailing = const [],
    Future<void> Function()? onTap,
    Future<void> Function()? onDoubleTap,
    Widget help = ds.HelpScope.None,
    bool highlighted = false,
  }) {
    return FutureBuilder<api.Known>(
      future: future,
      builder: (context, snapshot) {
        return KnownMediaRowDisplay(
          snapshot.data ?? api.Known(),
          key: key,
          leading: leading,
          trailing: trailing,
          onTap: onTap,
          onDoubleTap: onDoubleTap,
          help: help,
          highlighted: highlighted,
        );
      },
    );
  }

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
          Expanded(child: Text(current.description, overflow: TextOverflow.ellipsis)),
          ...trailing,
        ],
      ),
      help,
    );
  }
}
