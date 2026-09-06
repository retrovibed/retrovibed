import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart' as api;

class PublisherTypography extends StatelessWidget {
  final api.PluginPublisher current;
  final List<Widget> leading;
  final List<Widget> trailing;

  const PublisherTypography(
    this.current, {
    super.key,
    this.leading = const [],
    this.trailing = const [],
  });

  static String description(
    api.PluginPublisher current,
  ) {
    return current.description.isEmpty ? current.id : current.description;
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
