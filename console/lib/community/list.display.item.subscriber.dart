import 'package:flutter/material.dart';
import 'package:qr_flutter/qr_flutter.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart';
import 'qr.dart';
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
    final qrData = encodeQRPayload(community);

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
            ConstrainedBox(
              constraints: BoxConstraints(maxHeight: 256, maxWidth: 256),
              child: QrImageView(
                data: qrData,
                version: QrVersions.auto,
                backgroundColor: Colors.white,
                dataModuleStyle: QrDataModuleStyle(color: Colors.black),
              ),
            ),
            Text(
              'Scan this QR code to subscribe to this community',
              style: theme.textTheme.bodyMedium,
              textAlign: TextAlign.center,
            ),
            CommunityContentDisplay(community: community),
          ],
        ),
      ),
    );
  }
}
