import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:retrovibed/designkit.dart' as ds;

class SuggestionList extends StatefulWidget {
  final List<(Widget label, VoidCallback onSelect)> items;

  const SuggestionList(this.items, {super.key});

  @override
  State<SuggestionList> createState() => SuggestionListState();
}

class SuggestionListState extends State<SuggestionList> {
  int _selected = 0;

  void cycle([int delta = 1]) {
    if (widget.items.isEmpty) return;
    final n = widget.items.length;
    setState(() => _selected = (_selected + delta % n + n) % n);
  }

  bool get hasItems => widget.items.isNotEmpty;

  void select() {
    if (widget.items.isEmpty) return;
    widget.items[_selected.clamp(0, widget.items.length - 1)].$2();
  }

  @override
  void didUpdateWidget(SuggestionList old) {
    super.didUpdateWidget(old);
    if (old.items != widget.items) {
      _selected = 0;
    }
  }

  // Arrow keys move the highlight whenever focus is anywhere inside this
  // list (e.g. a ListTile reached via Tab). Flutter walks up the focus
  // chain calling each ancestor FocusNode's onKeyEvent until one handles
  // the event, so this fires for any descendant without the list's caller
  // needing to know anything about its internal structure.
  KeyEventResult _onKeyEvent(FocusNode node, KeyEvent event) {
    if (event is! KeyDownEvent || !hasItems) return KeyEventResult.ignored;
    if (event.logicalKey == LogicalKeyboardKey.arrowDown) {
      cycle();
      return KeyEventResult.handled;
    }
    if (event.logicalKey == LogicalKeyboardKey.arrowUp) {
      cycle(-1);
      return KeyEventResult.handled;
    }
    return KeyEventResult.ignored;
  }

  @override
  Widget build(BuildContext context) {
    if (widget.items.isEmpty) return ds.Empty;
    final defaults = ds.Defaults.of(context);
    return Focus(
      onKeyEvent: _onKeyEvent,
      child: ds.Container(
        padding: defaults.padding / 8,
        Material(
          elevation: 4,
          clipBehavior: Clip.antiAlias,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: widget.items.indexed.map(
              (entry) {
                final (i, item) = entry;
                return ListTile(
                  dense: true,
                  selected: i == _selected,
                  selectedTileColor: defaults.highlight,
                  title: item.$1,
                  trailing: const Icon(Icons.chevron_right, size: 16),
                  onTap: item.$2,
                );
              },
            ).toList(),
          ),
        ),
      ),
    );
  }
}
