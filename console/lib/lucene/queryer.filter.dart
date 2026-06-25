import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'parser.results.dart';

class QueryerFilterChip extends StatefulWidget {
  final ParserResult filter;
  final void Function(ParserResult, void Function(ParserResult), VoidCallback) onEdit;
  final VoidCallback onRemove;

  const QueryerFilterChip({
    super.key,
    required this.filter,
    required this.onEdit,
    required this.onRemove,
  });

  @override
  State<QueryerFilterChip> createState() => QueryerFilterChipState();
}

class QueryerFilterChipState extends State<QueryerFilterChip> {
  bool _open = false;

  void _toggle() {
    if (_open) {
      setState(() => _open = false);
      widget.onEdit(ParserResult.close, (_) {}, () {});
    } else {
      setState(() => _open = true);
      widget.onEdit(widget.filter, (_) {}, accept);
    }
  }

  void accept() {
    if (!_open) return;
    setState(() => _open = false);
    widget.onEdit(ParserResult.close, (_) {}, () {});
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final chipTheme = theme.chipTheme;
    final bgColor = _open
        ? (chipTheme.selectedColor ?? theme.colorScheme.secondaryContainer)
        : (chipTheme.backgroundColor ?? theme.colorScheme.surfaceContainerLow);

    return CallbackShortcuts(
      bindings: {
        const SingleActivator(LogicalKeyboardKey.backspace): widget.onRemove,
        const SingleActivator(LogicalKeyboardKey.delete): widget.onRemove,
      },
      child: Tooltip(
        message: _open ? 'Press Enter to accept' : '',
        child: InkWell(
          mouseCursor: SystemMouseCursors.click,
          onTap: _toggle,
          borderRadius: BorderRadius.circular(8),
          child: Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
            decoration: BoxDecoration(
              color: bgColor,
              borderRadius: BorderRadius.circular(8),
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                widget.filter,
                const SizedBox(width: 4),
                GestureDetector(
                  onTap: _open ? accept : widget.onRemove,
                  child: Tooltip(
                    message: _open ? 'Accept' : 'Remove',
                    child: Icon(_open ? Icons.check : Icons.close, size: 18),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
