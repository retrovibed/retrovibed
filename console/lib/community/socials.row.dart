import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart';
import 'socials.details.dart';

class SocialCommunityRow extends StatelessWidget {
  final Community community;
  final FnSocialsSearch details;
  final FnSocialsEnable enable;
  final FnSocialsDisable disable;
  final bool focused;
  final VoidCallback onInfo;

  const SocialCommunityRow({
    super.key,
    required this.community,
    required this.details,
    required this.enable,
    required this.disable,
    required this.focused,
    required this.onInfo,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return ds.Card(
      alignment: Alignment.topCenter,
      fit: FlexFit.tight,
      mainAxisSize: MainAxisSize.max,
      padding: EdgeInsets.zero,
      leading: [
        Flexible(
          child: SelectableText(
            communities.domain(community.url),
            style: theme.textTheme.titleLarge,
            maxLines: 1,
          ),
        ),
        Spacer(),
      ],
      Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisSize: MainAxisSize.max,
        children: [
          Text(
            community.description.isNotEmpty ? community.description : community.url,
            style: Theme.of(context).textTheme.titleSmall,
          ),
          Visibility(
            child: SocialCommunityDetails(
              community,
              search: details,
              enable: enable,
              disable: disable,
            ),
            visible: focused,
            maintainState: true,
          ),
        ],
      ),
      trailing: [
        Spacer(),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            const IconButton(
              icon: Icon(Icons.add_a_photo_outlined),
              onPressed: null,
              tooltip: "Photo",
            ),
            const IconButton(
              icon: Icon(Icons.videocam_outlined),
              onPressed: null,
              tooltip: "Video",
            ),
            const IconButton(
              icon: Icon(Icons.video_library_outlined),
              onPressed: null,
              tooltip: "Library",
            ),
            IconButton(
              icon: Icon(focused ? Icons.info : Icons.info_outline),
              onPressed: onInfo,
              tooltip: "Info",
            ),
          ],
        ),
      ],
    );
  }
}
