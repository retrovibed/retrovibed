import 'package:flutter/material.dart' as m;

class Checkbox extends m.StatelessWidget {
  final m.Widget? label;
  final m.Widget? description;
  final m.Widget? trailing;
  final bool dense;
  final bool value;
  final m.Alignment alignment;
  final void Function(bool?)? onChanged;
  Checkbox(
    this.label, {
    super.key,
    this.onChanged,
    this.trailing,
    this.description,
    this.value = false,
    this.dense = false,
    this.alignment = m.Alignment.center,
  });

  @override
  m.Widget build(m.BuildContext context) {
    return m.Material(
      color: m.Colors.transparent,
      child: m.Column(
        mainAxisSize: m.MainAxisSize.min,
        crossAxisAlignment: m.CrossAxisAlignment.start,
        children: [
          m.IntrinsicHeight(
            child: m.CheckboxListTile(
              dense: dense,
              title: label,
              secondary: trailing,
              subtitle: description,
              value: value,
              onChanged: onChanged,
              controlAffinity: m.ListTileControlAffinity.leading,
              visualDensity: m.VisualDensity.compact,
              materialTapTargetSize: m.MaterialTapTargetSize.shrinkWrap,
              contentPadding: m.EdgeInsets.zero,
            ),
          ),
        ],
      ),
    );
  }
}
