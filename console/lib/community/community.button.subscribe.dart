import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/community/api.dart';
import 'package:retrovibed/authn.dart' as authn;

class SubscribeButton extends StatelessWidget {
  final Community community;
  final void Function(Community)? onChanged;
  final FnSubscribe subscribe;

  const SubscribeButton({super.key, required this.community, this.onChanged, this.subscribe = API.subscribe});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final subscribed = community.subscribedAt.isNotEmpty;
    return ds.LoadingIconButton(
      icon: Icon(
        subscribed ? Icons.check_circle : Icons.add_circle_outline,
        color: subscribed ? theme.colorScheme.primary : null,
      ),
      onPressed: () {
        final auth = [authn.AuthzCache.bearer(context)];
        return httpx.withRetry(() => subscribe(community.id, options: auth)).then((v) => onChanged?.call(community));
      },
      tooltip: subscribed ? 'Unsubscribe' : 'Subscribe',
      help: ds.Hint(Text(subscribed ? "tap to unsubscribe from this community" : "tap to subscribe to this community")),
    );
  }
}
