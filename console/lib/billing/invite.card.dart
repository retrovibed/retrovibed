import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;

class InviteCard extends StatelessWidget {
  final EdgeInsets? margin;
  const InviteCard({super.key, this.margin});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return ds.Card(
      alignment: Alignment.center,
      margin: margin ?? defaults.margin,
      help: ds.Hint(
        label: const Text("Invite"),
        description: const Text("copies an invite link to your clipboard"),
      ),
      Column(
        mainAxisAlignment: MainAxisAlignment.center,
        spacing: defaults.spacing,
        children: [
          Text("Invite", style: theme.textTheme.titleMedium),
          Text(
            "share retrovibed with friends",
            style: theme.textTheme.bodySmall,
            textAlign: TextAlign.center,
          ),
          ds.LoadingIconButton(
            icon: Icon(Icons.person_add),
            tooltip: 'Copy Invite Link',
            onPressed: () {
              return authn.DeeppoolAuthzCache.attributionToken(context).then((token) {
                return Clipboard.setData(ClipboardData(text: 'https://invite.retrovibe.space/?a=$token'));
              });
            },
          ),
        ],
      ),
    );
  }
}
