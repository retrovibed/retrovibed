import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart' as api;

class DaemonTypography extends StatelessWidget {
  final api.Daemon current;
  final List<Widget> leading;
  final List<Widget> trailing;

  const DaemonTypography(
    this.current, {
    super.key,
    this.leading = const [],
    this.trailing = const [],
  });

  static String description(
    api.Daemon current,
  ) {
    final isDevice = api.daemons.isLocalDevice(current);
    return isDevice ? "local library" : current.description;
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final defaultDisplay = description(current);

    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      spacing: defaults.spacing,
      children: [
        ...leading,
        Flexible(child: Text(defaultDisplay, overflow: TextOverflow.ellipsis)),
        ...trailing,
      ],
    );
  }
}
