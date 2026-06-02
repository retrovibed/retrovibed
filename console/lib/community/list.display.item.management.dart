import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'api.dart';
import 'qr.attribution.dart';
import 'metrics.dashboard.dart';
import 'community.detail.dart';
import 'community.button.subscribe.dart';
import 'community.button.publish.dart';
import 'community.button.delete.dart';
import 'community.button.share.dart';
import 'community.update.dart';
import 'content.display.dart';

class ManagementListDisplayItem extends StatelessWidget {
  final Community community;
  final void Function(Community)? onChanged;
  final FnSubscribe subscribe;

  const ManagementListDisplayItem({
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
              constraints: BoxConstraints(
                maxWidth: defaults.compact + defaults.padding.horizontal,
                minHeight: defaults.compact + defaults.padding.vertical,
              ),
              community: community,
              update: (updated) {
                var auth = [authn.DeeppoolAuthzCache.bearer(context)];
                return httpx.withRetry(
                  () {
                    return API.update(
                      updated.id,
                      CommunityUpdateRequest(community: updated),
                      options: auth,
                    );
                  },
                );
              },
              onUpdate: (c) => onChanged?.call(c),
              onCancel: () {},
            ),
            QRAttribution(community: community),
            Divider(height: 32),
            MetricsDashboard(community: community),
            CommunityContentDisplay(community: community),
          ],
        ),
      ),
    );
  }
}
