import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/mimex.dart' as mimex;
import 'community.pb.dart';

class PublishedContentDetail extends StatelessWidget {
  final PublishedContent item;

  const PublishedContentDetail({super.key, required this.item});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return Column(
      spacing: defaults.spacing,
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (item.description.isNotEmpty) Text(item.description, style: theme.textTheme.bodyMedium),
        Wrap(
          spacing: defaults.spacing,
          children: [
            ds.Bytes(item.bytes),
            ds.Timestamp.iso8601(item.publishedAt, leading: Text('published ')),
            SelectableText(
              mimex.maybe(item.mimetype),
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.outline,
              ),
            ),
          ],
        ),
        SelectableText(
          item.magnetUri,
          style: theme.textTheme.bodySmall?.copyWith(
            color: theme.colorScheme.outline,
          ),
        ),
      ],
    );
  }
}
