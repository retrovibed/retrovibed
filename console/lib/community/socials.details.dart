import 'package:flutter/material.dart';
import 'package:retrovibed/community/community.detail.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart';
import 'socials.publishers.dart';

// Fetches its own catalog + enabled-publisher data for a single community —
// the socials search endpoint backs this expanded details view only, not
// the SocialHome grid itself.
class SocialCommunityDetails extends StatefulWidget {
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
  State<SocialCommunityDetails> createState() => _SocialCommunityDetailsState();
}

class _SocialCommunityDetailsState extends State<SocialCommunityDetails> with ds.LoadingState {
  @override
  void initState() {
    super.initState();
  }

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
      ds.Loading(
        loading: loading,
        cause: cause,
        Column(
          children: [
            CommunityDetail(community: widget.community),
            SocialsPublishers(widget.community),
          ],
        ),
      ),
    );
  }
}
