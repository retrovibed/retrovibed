import 'package:flutter/material.dart';
import 'package:retrovibed/community/community.detail.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart';
import 'socials.publishers.dart';

// The expanded Info panel for a single community: what it is, and the
// publishers it publishes through.
class SocialCommunityDetails extends StatelessWidget {
  final Community community;
  final FnSocialsSearch search;
  final FnSocialsEnable enable;
  final FnSocialsDisable disable;

  const SocialCommunityDetails(
    this.community, {
    super.key,
    required this.search,
    required this.enable,
    required this.disable,
  });

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final theme = Theme.of(context);

    return ds.Container(
      padding: defaults.padding,
      decoration: BoxDecoration(
        color: theme.colorScheme.surfaceContainerLow,
        border: defaults.border,
        borderRadius: defaults.borderRadius,
      ),
      Column(
        children: [
          CommunityDetail(community: community),
          SocialsPublishers(
            community,
            search: search,
            enable: enable,
            disable: disable,
          ),
        ],
      ),
    );
  }
}
