import 'package:flutter/material.dart';
import 'package:retrovibed/timex.dart' as timex;
import 'package:retrovibed/design.kit/typography/timerange.dart' as typography;

class TimeRange extends StatelessWidget {
  final List<timex.Range> segments;
  final timex.Range selected;
  final ValueChanged<timex.Range> onChanged;
  final double estimate;

  const TimeRange({
    super.key,
    required this.segments,
    required this.selected,
    required this.onChanged,
    this.estimate = 84.0,
  });

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final estimatedWidth = segments.length * estimate;

        if (constraints.maxWidth < estimatedWidth) {
          return _buildDropdown(context);
        }

        return _buildSegmented(context);
      },
    );
  }

  Widget _buildSegmented(BuildContext context) {
    return SegmentedButton<timex.Range>(
      segments:
          segments
              .map(
                (range) => ButtonSegment(
                  value: range,
                  label: typography.TimeRange(range),
                ),
              )
              .toList(),
      selected: {selected},
      onSelectionChanged: (selection) => onChanged(selection.first),
    );
  }

  Widget _buildDropdown(BuildContext context) {
    final value = segments.contains(selected) ? selected : segments.first;
    return DropdownButton<timex.Range>(
      value: value,
      isExpanded: true,
      onChanged: (value) {
        if (value != null) onChanged(value);
      },
      items:
          segments
              .map(
                (range) => DropdownMenuItem(
                  value: range,
                  child: typography.TimeRange(range),
                ),
              )
              .toList(),
    );
  }
}
