import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/community/community.pb.dart';

class CommunityDetail extends StatelessWidget {
  final Community community;

  const CommunityDetail({super.key, required this.community});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return Column(
      spacing: defaults.spacing,
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          spacing: defaults.spacing,
          children: [
            Flexible(child: Text(community.domain, style: theme.textTheme.titleLarge, overflow: TextOverflow.ellipsis)),
            Visibility(
              visible: community.hidden,
              child: Icon(Icons.lock, size: 16, color: theme.colorScheme.outline),
            ),
          ],
        ),
        SelectableText(
          community.url,
          style: theme.textTheme.bodySmall?.copyWith(
            color: theme.colorScheme.primary,
          ),
        ),
        if (community.description.isNotEmpty) Text(community.description, style: theme.textTheme.bodyMedium),
        ds.Timestamp.iso8601(community.createdAt, leading: Text('Created: ')),
      ],
    );
  }
}
