import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/authn/login.dart';

class DeveloperSettings extends StatelessWidget {
  final EdgeInsets margin;
  final BoxConstraints? constraints;
  const DeveloperSettings({super.key, this.margin = EdgeInsets.zero, this.constraints});

  @override
  Widget build(BuildContext context) {
    final flags = Login.cached(context).flags;
    final theme = Theme.of(context);

    return ds.Card(
      alignment: Alignment.topLeft,
      margin: margin,
      constraints: constraints,
      help: ds.Hint(const Text("developer-only feature flags")),
      Column(
        children: [
          Text("Developer Settings", textAlign: TextAlign.center, style: theme.textTheme.titleMedium),
          forms.Checkbox(
            const Text('Networking'),
            description: const Text('Enable networking functionality UX'),
            value: flags.networking,
            onChanged: (v) {
              final s = Login.of(context);
              s?.setState(() => s.flags = flags.copyWith(networking: v ?? false));
            },
          ),
          forms.Checkbox(
            const Text('Subscription'),
            description: const Text('Force enable subscription management UX'),
            value: flags.subscription,
            onChanged: (v) {
              final s = Login.of(context);
              s?.setState(() => s.flags = flags.copyWith(subscription: v ?? false));
            },
          ),
          forms.Checkbox(
            const Text('Recommendations'),
            description: const Text('Toggle recommendations panel'),
            value: flags.recommendations,
            onChanged: (v) {
              final s = Login.of(context);
              s?.setState(() => s.flags = flags.copyWith(recommendations: v ?? false));
            },
          ),
          forms.Checkbox(
            const Text('Releases'),
            description: const Text('Toggle releases panel'),
            value: flags.releases,
            onChanged: (v) {
              final s = Login.of(context);
              s?.setState(() => s.flags = flags.copyWith(releases: v ?? false));
            },
          ),
          forms.Checkbox(
            const Text('Debug'),
            description: const Text('Enable debug-only UX and tuning panels'),
            value: flags.debug,
            onChanged: (v) {
              final s = Login.of(context);
              s?.setState(() => s.flags = flags.copyWith(debug: v ?? false));
            },
          ),
        ],
      ),
    );
  }
}
