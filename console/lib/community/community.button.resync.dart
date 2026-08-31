import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/community/api.dart';
import 'package:retrovibed/authn.dart' as authn;

class ResyncButton extends StatelessWidget {
  final Community community;
  final void Function(Community)? onResynced;
  final FnPublishingSearch apiresync;

  const ResyncButton({
    super.key,
    required this.community,
    this.onResynced,
    this.apiresync = communities.resync,
  });

  @override
  Widget build(BuildContext context) {
    return ds.LoadingIconButton.refresh(
      tooltip: 'Resync Community',
      help: ds.Hint(const Text("refetch metadata and latest published content from the source")),
      onPressed: () {
        final auth = [authn.request(authn.AuthzCache.meta(context))];
        return httpx.withRetry(() => apiresync(community.id, options: auth)).then((resp) {
          onResynced?.call(resp.community);
        });
      },
    );
  }
}
