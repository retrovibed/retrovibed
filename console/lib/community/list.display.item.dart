import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/community/api.dart';
import 'package:retrovibed/community/list.display.item.management.dart';
import 'package:retrovibed/community/list.display.item.subscriber.dart';

class ListDisplayItem extends StatelessWidget {
  final Community community;
  final void Function(Community)? onChanged;
  final FnSubscribe subscribe;

  const ListDisplayItem({
    super.key,
    required this.community,
    this.onChanged,
    this.subscribe = communities.subscribe,
  });

  @override
  Widget build(BuildContext context) {
    final session = authn.Authenticated.syncSession(context);
    final owned = session.account.id == community.accountId;

    return ds.Guarded(
      key: ValueKey(community.id),
      enabled: !owned,
      overlay: SubscriberListDisplayItem(
        community: community,
        onChanged: onChanged,
        subscribe: subscribe,
      ),
      child: ManagementListDisplayItem(
        community: community,
        onChanged: onChanged,
        subscribe: subscribe,
      ),
    );
  }
}
