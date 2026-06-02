import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart';
import 'community.detail.dart';
import 'community.button.subscribe.dart';
import 'content.display.dart';

class SubscriberListDisplayItem extends StatelessWidget {
  final Community community;
  final void Function(Community)? onChanged;
  final FnSubscribe subscribe;

  const SubscriberListDisplayItem({
    super.key,
    required this.community,
    this.onChanged,
    this.subscribe = API.subscribe,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return ds.TableRow(
      [
        Expanded(child: CommunityDetail(community: community)),
        SubscribeButton(
          community: community,
          onChanged: onChanged,
          subscribe: subscribe,
        ),
      ],
      expanded: Container(
        padding: defaults.padding,
        decoration: BoxDecoration(
          color: theme.colorScheme.surfaceContainerHighest.withValues(
            alpha: 0.3,
          ),
          borderRadius: BorderRadius.only(
            bottomLeft: Radius.circular(12),
            bottomRight: Radius.circular(12),
          ),
        ),
        child: Column(
          spacing: defaults.spacing,
          children: [
            CommunityContentDisplay(community: community),
          ],
        ),
      ),
    );
  }
}
