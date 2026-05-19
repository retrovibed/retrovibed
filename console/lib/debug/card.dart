import 'dart:io';

import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/retrovibed.dart' as retro;

class Card extends StatelessWidget {
  final EdgeInsets margin;
  const Card({super.key, this.margin = EdgeInsets.zero});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);
    final mq = MediaQuery.of(context);
    final size = mq.size;

    final rows = [
      ('OS', '${Platform.operatingSystem} ${Platform.operatingSystemVersion}'),
      ('Screen', '${size.width.toInt()} × ${size.height.toInt()}'),
      ('Pixel ratio', mq.devicePixelRatio.toStringAsFixed(2)),
      ('Text scale', mq.textScaler.scale(1.0).toStringAsFixed(2)),
      ('Platform', theme.platform.name),
      ('Desktop', defaults.desktop.toString()),
      ('Mobile', defaults.mobile.toString()),
      ('Compact', defaults.isCompact.toString()),
      ('Metered', retro.metered().toString()),
    ];

    return ds.Card(
      alignment: Alignment.topLeft,
      margin: margin,
      help: ds.Hint(const Text("device and display information")),
      SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          spacing: defaults.spacing / 4,
          children: [
            Text("Device Info", style: theme.textTheme.titleMedium),
            ...rows.map(
              (r) => Row(
                children: [
                  SizedBox(
                    width: 84,
                    child: Text(
                      r.$1,
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: theme.colorScheme.onSurface.withValues(alpha: 0.6),
                      ),
                    ),
                  ),
                  Expanded(
                    child: Text(
                      r.$2,
                      style: theme.textTheme.bodySmall,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
