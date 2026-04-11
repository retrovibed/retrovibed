import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;

class VisibilitySelector extends StatelessWidget {
  final bool hidden;
  final ValueChanged<bool> onChanged;

  const VisibilitySelector({
    super.key,
    required this.hidden,
    required this.onChanged,
  });

  static const _descriptions = {
    false: 'Discoverable via search.',
    true: 'Hidden from search for other accounts.',
  };

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return Column(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.start,
      spacing: defaults.spacing,
      children: [
        Text('Visibility', style: theme.textTheme.titleSmall),
        SegmentedButton<bool>(
          segments: [
            ButtonSegment(value: false, label: Text('Public')),
            ButtonSegment(value: true, label: Text('Private')),
          ],
          selected: {hidden},
          onSelectionChanged: (selection) => onChanged(selection.first),
        ),
        Text(_descriptions[hidden] ?? '', style: theme.textTheme.bodySmall),
      ],
    );
  }
}
