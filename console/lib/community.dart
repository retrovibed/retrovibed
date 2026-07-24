import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/community/list.display.dart';

export 'package:retrovibed/community/api.dart';

class AutoHelp extends StatelessWidget {
  final Widget child;
  const AutoHelp(this.child, {super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return ds.HelpAuto(
      child,
      cacheid: 'community',
      title: Text("Community", style: theme.textTheme.titleMedium),
      content: const Text(
        "Connect and share with other users. Browse community-curated "
        "collections, join groups, and discover new content shared by "
        "people with similar tastes.",
      ),
    );
  }
}

class Management extends StatelessWidget {
  const Management({super.key});

  @override
  Widget build(BuildContext context) {
    return ds.build((context) {
      final defaults = ds.Defaults.of(context);
      return ds.Container(
        margin: EdgeInsets.zero,
        padding: defaults.padding,
        ListDisplay(),
      );
    });
  }
}
