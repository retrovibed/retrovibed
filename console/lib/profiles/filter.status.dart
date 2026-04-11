import 'package:flutter/material.dart';
import 'package:retrovibed/meta.dart' as meta;

class FilterStatus extends StatelessWidget {
  static const _labels = {
    meta.ProfileStatus.ENABLED: "Enabled",
    meta.ProfileStatus.PENDING: "Pending",
    meta.ProfileStatus.DISABLED: "Disabled",
    meta.ProfileStatus.NONE: "Any",
  };

  static String _label(meta.ProfileStatus s) => _labels[s] ?? "Any";

  final meta.ProfileStatus current;
  final void Function(meta.ProfileStatus? v) onChange;
  const FilterStatus(this.current, this.onChange, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return PopupMenuButton<meta.ProfileStatus>(
      initialValue: current,
      tooltip: "Filter by status",
      onSelected: onChange,
      itemBuilder: (context) => [
        PopupMenuItem(value: meta.ProfileStatus.ENABLED, child: Text("Enabled")),
        PopupMenuItem(value: meta.ProfileStatus.PENDING, child: Text("Pending")),
        PopupMenuItem(value: meta.ProfileStatus.DISABLED, child: Text("Disabled")),
        PopupMenuItem(value: meta.ProfileStatus.NONE, child: Text("Any")),
      ],
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.filter_list),
          Text(_label(current)),
        ],
      ),
    );
  }
}
