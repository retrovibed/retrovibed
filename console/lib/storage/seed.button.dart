import 'package:flutter/material.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;
import './seed.dart' as seed;

// compact seed selector, only the current seed's icon is shown until clicked.
class SeedButton extends StatelessWidget {
  final String current;
  final void Function(seed.Seed)? onChange;

  const SeedButton(this.current, {super.key, this.onChange});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final classified = seed.Classifier(community: uuidx.min()).classify(current);
    final options = [
      seed.Seed.global(),
      seed.Seed.unique(uuidx.random()),
    ];

    return PopupMenuButton<String>(
      enabled: onChange != null,
      position: PopupMenuPosition.under,
      color: theme.colorScheme.surface,
      surfaceTintColor: theme.colorScheme.surface,
      tooltip: classified.tooltip,
      icon: Icon(classified.icon),
      onSelected: (id) => onChange?.call(options.firstWhere((s) => s.id == id)),
      itemBuilder: (context) => options.map((s) {
        return PopupMenuItem<String>(
          value: s.id,
          child: ListTile(
            leading: Icon(s.icon),
            title: s.label,
            subtitle: Text(s.id, maxLines: 1, overflow: TextOverflow.ellipsis),
          ),
        );
      }).toList(),
    );
  }
}
