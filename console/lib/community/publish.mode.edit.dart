import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/billing.dart' as billing;
import 'package:retrovibed/community/community.pb.dart';
import 'package:retrovibed/community/publish.mode.dart';

class PublishModeEdit extends StatelessWidget {
  final PublishMode publishMode;
  final ValueChanged<PublishMode> onChanged;

  const PublishModeEdit({
    super.key,
    required this.publishMode,
    required this.onChanged,
  });

  static const _descriptions = {
    PublishMode.UNLISTED: 'Self hosted distribution.',
    PublishMode.LISTED: 'Discoverable via communities.',
    PublishMode.SYNDICATED: 'Hosted and distributed by retrovibe.',
  };

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);
    final (plan, _) = billing.PlanSummary.fromPlan(
      billing.Registered.of(context).plan,
    );
    final maxMode = maxPublishMode(plan.id);
    final effective = publishMode.value > maxMode.value ? maxMode : publishMode;

    return Column(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.start,
      spacing: defaults.spacing,
      children: [
        Text('Publish Mode', style: theme.textTheme.titleSmall),
        SegmentedButton<PublishMode>(
          showSelectedIcon: false,
          expandedInsets: EdgeInsets.zero,
          segments: [
            ButtonSegment(value: PublishMode.UNLISTED, label: Text('Unlisted')),
            ButtonSegment(
              value: PublishMode.LISTED,
              label: Text('Listed'),
              enabled: maxMode.value >= PublishMode.LISTED.value,
              tooltip: "content is listed by retrovibed. requires a paid plan and to have publishing enabled",
            ),
            ButtonSegment(
              value: PublishMode.SYNDICATED,
              label: Text('Syndicated'),
              enabled: maxMode.value >= PublishMode.SYNDICATED.value,
              tooltip:
                  "content is listed and distributed by retrovibed. requires founder, premium, and family plan and to have publishing enabled",
            ),
          ],
          selected: {effective},
          onSelectionChanged: (selection) => onChanged(selection.first),
        ),
        Text(_descriptions[effective] ?? '', style: theme.textTheme.bodySmall),
      ],
    );
  }
}
