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
  return PopupMenuItem<String>(
    child: ValueListenableBuilder<SearchMode>(
      valueListenable: current,
      builder: (context, currentMode, _) {
        final selected = mode == currentMode;
        return ListTile(
          leading: Icon(selected ? Icons.check : icon),
          title: Text(label),
          selected: selected,
          hoverColor: Colors.transparent,
          onTap: () => onSelect(selected ? SearchMode.library : mode),
        );
      },
    ),
  );
}
