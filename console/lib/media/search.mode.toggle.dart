import 'package:flutter/material.dart';
import 'search.mode.dart';

// Reused by the library, discovery, and remote search trays to switch
// between modes; icon/label are caller-supplied so each mode's presentation
// stays with the screen that knows about it.
PopupMenuItem<String> SearchModeToggle({
  required SearchMode mode,
  required ValueNotifier<SearchMode> current,
  required IconData icon,
  required String label,
  required void Function(SearchMode) onSelect,
}) {
  final selected = mode == current.value;
  return PopupMenuItem<String>(
    enabled: !selected,
    mouseCursor: selected ? SystemMouseCursors.basic : SystemMouseCursors.click,
    onTap: () => onSelect(selected ? SearchMode.library : mode),
    child: Row(
      children: [
        Icon(icon),
        const SizedBox(width: 12),
        Text(label),
      ],
    ),
  );
}
