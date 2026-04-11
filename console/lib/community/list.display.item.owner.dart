import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/community/api.dart';
import 'package:retrovibed/community/qr.attribution.dart';
import 'package:retrovibed/community/metrics.dashboard.dart';
import 'package:retrovibed/community/community.detail.dart';
import 'package:retrovibed/community/community.button.subscribe.dart';
import 'package:retrovibed/community/community.button.publish.dart';
import 'package:retrovibed/community/community.button.delete.dart';
import 'package:retrovibed/community/community.button.share.dart';
import 'package:retrovibed/community/community.update.dart';

class OwnerListDisplayItem extends StatelessWidget {
  final Community community;
  final void Function(Community)? onChanged;
  final FnSubscribe subscribe;

  const OwnerListDisplayItem({
    super.key,
    required this.community,
    this.onChanged,
    this.subscribe = API.subscribe,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return ds.TableRow.single(
      key: ValueKey(community.id),
      ds.CompactingMenu(
        CommunityDetail(community: community),
        trailing: [
          ds.CompactingMenu.pinned(
            SubscribeButton(
              community: community,
              onChanged: onChanged,
              subscribe: subscribe,
            ),
          ),
          ShareButton(community: community),
          PublishButton(community: community),
          DeleteButton(
            community: community,
            onDeleted: onChanged,
          ),
        ],
      ),
      expanded: ds.Container(
        padding: defaults.padding,
        decoration: BoxDecoration(
          color: theme.colorScheme.surfaceContainerHighest.withValues(
            alpha: 0.3,
          ),
          borderRadius: defaults.borderRadius,
        ),
        Wrap(
          alignment: WrapAlignment.center,
          spacing: defaults.spacing,
          children: [
            CommunityUpdate(
              constraints: BoxConstraints(maxWidth: defaults.compact),
              community: community,
              update: (updated) {
                var auth = [authn.AuthzCache.bearer(context)];
                return httpx.withRetry(
                  () => API.update(
                    updated.id,
                    CommunityUpdateRequest(community: updated),
                    options: auth,
                  ),
                );
              },
              onUpdate: (c) => onChanged?.call(c),
              onCancel: () {},
            ),
            Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                QRAttribution(community: community),
                Text(
                  'Scan this QR code to subscribe to this community',
                  style: theme.textTheme.bodyMedium,
                  textAlign: TextAlign.center,
                ),
              ],
            ),
            Divider(height: 32),
            MetricsDashboard(community: community),
          ],
        ),
      ),
    );
  }
}
