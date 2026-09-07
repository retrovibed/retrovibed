import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './current.dart';
import './authz.meta.display.dart';

class Overview extends StatelessWidget {
  const Overview({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);
    return ds.Container(
      padding: defaults.padding,
      decoration: BoxDecoration(color: theme.colorScheme.surfaceContainerLow),
      Column(
        mainAxisSize: MainAxisSize.max,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        spacing: defaults.spacing,
        children: [
          Current(),
          AuthzMetaDisplay.current(),
        ],
      ),
    );
  }
}
