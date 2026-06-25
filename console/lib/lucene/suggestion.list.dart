import 'package:flutter/material.dart';

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

  @override
  Widget build(BuildContext context) {
    if (widget.items.isEmpty) return const SizedBox.shrink();
    final theme = Theme.of(context);
    return Material(
      elevation: 4,
      color: theme.colorScheme.surface,
      borderRadius: BorderRadius.circular(4),
      clipBehavior: Clip.antiAlias,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: widget.items.indexed.map(
          (entry) {
            final (i, item) = entry;
            return SizedBox(
              height: 32,
              child: ListTile(
                dense: true,
                selected: i == _selected,
                title: item.$1,
                trailing: const Icon(Icons.chevron_right, size: 16),
                onTap: item.$2,
              ),
            );
          },
        ).toList(),
      ),
    );
  }
}
