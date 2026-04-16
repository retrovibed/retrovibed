import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './field.dart';
import './suggestion.list.dart';
import 'parser.results.dart';
import './parser.states.dart';

class SuggestionKeyScope extends InheritedWidget {
  final GlobalKey<SuggestionListState> suggestionKey;

  const SuggestionKeyScope({required this.suggestionKey, required super.child});

  static GlobalKey<SuggestionListState>? of(BuildContext context) =>
      context.dependOnInheritedWidgetOfExactType<SuggestionKeyScope>()?.suggestionKey;

  @override
  bool updateShouldNotify(SuggestionKeyScope old) => old.suggestionKey != suggestionKey;
}

class FilterChip extends StatefulWidget {
  final ParserResult filter;
  final void Function(ParserResult, void Function(ParserResult), VoidCallback) onEdit;
  final VoidCallback onRemove;

  const FilterChip({
    super.key,
    required this.filter,
    required this.onEdit,
    required this.onRemove,
  });

  @override
  State<FilterChip> createState() => FilterChipState();
}

class FilterChipState extends State<FilterChip> {
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
    final bgColor =
        _open
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

class Queryer extends StatefulWidget {
  static const Widget zerobox = SizedBox();

  final void Function(String) onQuery;
  final List<Field> fields;
  final InputDecoration? decoration;
  final bool autofocus;
  final bool disabled;
  final TextEditingController? controller;
  final FocusNode? focusNode;
  final List<Widget> leading;
  final List<Widget> trailing;
  final Widget help;

  const Queryer(
    this.onQuery,
    this.fields, {
    super.key,
    this.decoration,
    this.autofocus = false,
    this.disabled = false,
    this.controller,
    this.focusNode,
    this.leading = const [],
    this.trailing = const [],
    this.help = ds.HelpScope.None,
  });

  @override
  State<Queryer> createState() => _QueryerState();
}

class _QueryerState extends State<Queryer> {
  late TextEditingController _ctrl;
  final GlobalKey<SuggestionListState> _suggestionKey = GlobalKey();
  List<ParserResult> _filters = [];
  Widget? _updating;
  bool _editing = false;
  Parser _parser = Parser([], (ctx, range, content, {completed}) {});

  void _resetParser() => _parser = Parser(widget.fields, _replace);

  @override
  void initState() {
    super.initState();
    _ctrl = widget.controller ?? TextEditingController();
    _resetParser();
    _ctrl.addListener(_onText);
  }

  @override
  void dispose() {
    _ctrl.removeListener(_onText);
    if (widget.controller == null) _ctrl.dispose();
    super.dispose();
  }

  void _replace(
    Context ctx,
    TextRange range,
    String contents, {
    ParserResult? completed,
  }) {
    if (_editing) return; // Prevent recursion
    try {
      _editing = true;
      setState(() {
        if (completed != null) {
          completed.apply(_parser);
          _filters.add(completed);
        }
      });

      WidgetsBinding.instance.addPostFrameCallback((_) {
        setState(() {
          _ctrl.value = _ctrl.value.replaced(range, contents);
        });
        ds.textediting.refocus(_ctrl);
      });
    } finally {
      _editing = false;
    }
  }

  void _onText() {
    if (_ctrl.text.isEmpty && _filters.isEmpty) {
      return setState(_resetParser);
    }

    setState(() {
      _parser.consume(_ctrl);
    });
  }

  void _editFilter(
    ParserResult filter,
    void Function(ParserResult) onChanged,
    VoidCallback closeChip,
  ) {
    var current = filter;
    setState(() {
      final _w = current.edit((upd) {
        setState(() {
          _filters = _filters.map<ParserResult>((e) => e == current ? upd : e).toList();
          current = upd;
        });
        onChanged(upd);
      });
      final focusNode = FocusNode();
      _updating =
          _w == null
              ? null
              : Focus(
                focusNode: focusNode,
                onKeyEvent: (node, event) {
                  if (event.logicalKey != LogicalKeyboardKey.enter) return KeyEventResult.ignored;
                  if (event is! KeyDownEvent) return KeyEventResult.ignored;

                  closeChip();
                  return KeyEventResult.handled;
                },
                child: _w,
              );
      if (_w != null) focusNode.requestFocus();
    });
  }

  void _removeFilter(ParserResult filter) {
    setState(() {
      _filters.removeWhere((v) => v == filter);
    });

    // Reset field to its default value and restore it in the parser's field list.
    filter.reset(_parser);
  }

  bool _partialParse() {
    return !(_parser.current is Query);
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final chips =
        _filters.map((e) {
          return FilterChip(
            filter: e,
            onEdit: (filter, onChanged, closeChip) => _editFilter(filter, onChanged, closeChip),
            onRemove: () => _removeFilter(e),
          );
        }).toList();

    return ds.Shortcuts(
      enabled: defaults.desktop,
      bindings: {
        const SingleActivator(LogicalKeyboardKey.escape): (
          const Text('reset search'),
          () {
            print("queryer: escape");
            setState(_resetParser);
            return KeyEventResult.ignored;
          },
        ),
        const SingleActivator(LogicalKeyboardKey.arrowDown): (
          const Text('next suggestion'),
          () {
            print("queryer: arrowDown");
            _suggestionKey.currentState?.cycle();
            return KeyEventResult.handled;
          },
        ),
        const SingleActivator(LogicalKeyboardKey.arrowUp): (
          const Text('previous suggestion'),
          () {
            print("queryer: arrowUp");
            _suggestionKey.currentState?.cycle(-1);
            return KeyEventResult.handled;
          },
        ),
        const SingleActivator(LogicalKeyboardKey.enter): (
          const Text('select suggestion'),
          () {
            if (_suggestionKey.currentState?.hasItems ?? false) {
              print("queryer: enter -> select suggestion");
              _suggestionKey.currentState?.select();
              return KeyEventResult.handled;
            }
            print("queryer: enter -> no suggestion");
            return KeyEventResult.ignored;
          },
        ),
      },
      Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          ds.CompactingMenu(
            ds.Help(
              TextField(
                controller: _ctrl,
                enabled: !widget.disabled,
                autofocus: widget.autofocus,
                focusNode: widget.focusNode,
                decoration: (widget.decoration ??
                        const InputDecoration(
                          hintText: 'Search… (@ for filters)',
                        ))
                    .copyWith(
                      isDense: true,
                      contentPadding: const EdgeInsets.symmetric(
                        horizontal: 8,
                        vertical: 12,
                      ),
                    ),
                onSubmitted: (v) {
                  if (_partialParse()) return;
                  widget.onQuery(v);
                  ds.textediting.refocus(_ctrl);
                },
              ),
              widget.help,
            ),
            leading: widget.leading,
            trailing: widget.trailing,
          ),
          TextFieldTapRegion(
            child: Wrap(
              spacing: defaults.spacing,
              runSpacing: defaults.spacing / 2,
              children: chips,
            ),
          ),
          TextFieldTapRegion(
            child: SuggestionKeyScope(
              suggestionKey: _suggestionKey,
              child: _updating ?? _parser.current,
            ),
          ),
        ],
      ),
    );
  }
}
